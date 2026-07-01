package usersync

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Session lifecycle timings. The client ends idle sessions cleanly before the
// server's 60-min backstop fires (see SESSION_HANDSHAKE_SPEC.md §6).
const (
	sessTick       = 30 * time.Second // how often we evaluate active sessions
	sessHeartbeat  = 2 * time.Minute  // ping quiet-but-active sessions this often
	sessClientIdle = 30 * time.Minute // end a session after this long with no uploads
)

// sessionState is one locally-tracked session. id is 0 until /start returns (or
// stays 0 when logged out — uploads then fall back to server-side auto-grouping).
type sessionState struct {
	kind      string
	serverID  int
	player    string
	id        int64
	startedAt time.Time
	lastActAt time.Time // last real upload (drives idle-end + whether to heartbeat)
	lastBeat  time.Time // last heartbeat we sent
}

// SessionManager owns the client side of the session handshake (start/end/
// heartbeat/end-all) and is the source of truth for "is a session active".
type SessionManager struct {
	baseURL    string
	token      func() (jwt, userID string, ok bool)
	serverID   func() int
	playerName func() string
	client     *http.Client
	logf       func(string)

	mu     sync.Mutex
	active map[string]*sessionState

	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewSessionManager builds a manager. baseURL is the ingest origin; token returns
// the bearer JWT (ok=false when logged out).
func NewSessionManager(baseURL string, token func() (string, string, bool), serverID func() int, playerName func() string, logf func(string)) *SessionManager {
	if logf == nil {
		logf = func(string) {}
	}
	return &SessionManager{
		baseURL:    baseURL,
		token:      token,
		serverID:   serverID,
		playerName: playerName,
		client:     &http.Client{Timeout: 10 * time.Second},
		logf:       logf,
		active:     map[string]*sessionState{},
	}
}

// Start marks a kind active immediately and, if logged in, opens a server session
// in the background and records its id.
func (m *SessionManager) Start(kind string) {
	now := time.Now()
	s := &sessionState{kind: kind, startedAt: now, lastActAt: now, lastBeat: now}
	if m.serverID != nil {
		s.serverID = m.serverID()
	}
	if m.playerName != nil {
		s.player = m.playerName()
	}
	m.mu.Lock()
	m.active[kind] = s
	m.mu.Unlock()
	go m.startOnServer(s)
}

func (m *SessionManager) startOnServer(s *sessionState) {
	jwt, _, ok := m.token()
	if !ok {
		return // logged out: no server session; uploads fall back to auto-grouping
	}
	var resp struct {
		Ok        bool  `json:"ok"`
		SessionID int64 `json:"sessionId"`
	}
	body := map[string]any{"kind": s.kind, "serverId": s.serverID}
	if s.player != "" {
		body["playerName"] = s.player
	}
	if !m.postJSON(jwt, "/user/session/start", body, &resp) || resp.SessionID == 0 {
		return
	}
	m.mu.Lock()
	cur := m.active[s.kind]
	if cur == s {
		cur.id = resp.SessionID
		m.mu.Unlock()
		m.logf(fmt.Sprintf("Opened %s session #%d on your account.", s.kind, resp.SessionID))
		return
	}
	m.mu.Unlock()
	// the session was stopped before /start returned — close the orphan immediately
	go m.endOnServer(resp.SessionID, s.kind)
}

// End clears a kind's session locally and closes it server-side. With a known id
// we end it precisely; otherwise we fall back to end-all-by-kind so no server
// session leaks until the backstop.
func (m *SessionManager) End(kind string) {
	m.mu.Lock()
	s := m.active[kind]
	delete(m.active, kind)
	m.mu.Unlock()
	if s == nil {
		return
	}
	if s.id != 0 {
		go m.endOnServer(s.id, kind)
		return
	}
	// id unknown — close any active session of this kind server-side.
	go m.postEndAll(kind)
}

func (m *SessionManager) endOnServer(id int64, kind string) {
	jwt, _, ok := m.token()
	if !ok {
		return
	}
	if m.postJSON(jwt, "/user/session/end", map[string]any{"sessionId": id}, nil) {
		m.logf(fmt.Sprintf("Closed %s session #%d.", kind, id))
	}
}

// EndAll closes all active sessions (kind == "") or all of one kind, locally and
// server-side. Synchronous (best-effort) so it can run on shutdown / region change.
func (m *SessionManager) EndAll(kind string) {
	m.mu.Lock()
	had := false
	for k := range m.active {
		if kind == "" || k == kind {
			delete(m.active, k)
			had = true
		}
	}
	m.mu.Unlock()
	if had {
		m.postEndAll(kind)
	}
}

// postEndAll sends the end-all request (kind == "" closes everything). It does
// not touch local state, so it's safe to call after the active map is cleared.
func (m *SessionManager) postEndAll(kind string) {
	jwt, _, ok := m.token()
	if !ok {
		return
	}
	body := map[string]any{}
	if kind != "" {
		body["kind"] = kind
	}
	if m.postJSON(jwt, "/user/session/end-all", body, nil) {
		if kind != "" {
			m.logf("Closed " + kind + " session(s) on your account.")
		} else {
			m.logf("Closed all sessions on your account.")
		}
	}
}

// OnRegionChange closes the old region's sessions and restarts the same kinds
// against the new region so data never lands under the wrong server.
func (m *SessionManager) OnRegionChange() {
	m.mu.Lock()
	kinds := make([]string, 0, len(m.active))
	for k := range m.active {
		kinds = append(kinds, k)
	}
	m.mu.Unlock()
	if len(kinds) == 0 {
		return
	}
	m.EndAll("")
	for _, k := range kinds {
		m.Start(k) // Start() reads the now-current serverID
	}
	m.logf("Region changed — restarted active sessions on the new server.")
}

// CurrentID returns the server session id for a kind (0 if none / not yet open).
func (m *SessionManager) CurrentID(kind string) int64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s := m.active[kind]; s != nil {
		return s.id
	}
	return 0
}

// MarkActivity records that a real upload happened for a kind (resets its idle
// clock and suppresses the next heartbeat).
func (m *SessionManager) MarkActivity(kind string) {
	now := time.Now()
	m.mu.Lock()
	if s := m.active[kind]; s != nil {
		s.lastActAt = now
		s.lastBeat = now
	}
	m.mu.Unlock()
}

// Active reports whether a kind has an active session.
func (m *SessionManager) Active(kind string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.active[kind]
	return ok
}

// Any reports whether any session is active.
func (m *SessionManager) Any() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.active) > 0
}

// StartedAt returns a kind's session start time (ok=false if inactive).
func (m *SessionManager) StartedAt(kind string) (time.Time, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s := m.active[kind]; s != nil {
		return s.startedAt, true
	}
	return time.Time{}, false
}

// Run starts the background heartbeat / idle-end loop.
func (m *SessionManager) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		t := time.NewTicker(sessTick)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				m.tick()
			}
		}
	}()
}

// Stop halts the background loop (does not end sessions — call EndAll for that).
func (m *SessionManager) Stop() {
	if m.cancel != nil {
		m.cancel()
	}
	m.wg.Wait()
}

func (m *SessionManager) tick() {
	now := time.Now()
	var toEnd []string
	type beat struct {
		id   int64
		kind string
	}
	var beats []beat
	m.mu.Lock()
	for k, s := range m.active {
		idle := now.Sub(s.lastActAt)
		if idle >= sessClientIdle {
			toEnd = append(toEnd, k)
			continue
		}
		if s.id != 0 && idle >= sessHeartbeat && now.Sub(s.lastBeat) >= sessHeartbeat {
			s.lastBeat = now
			beats = append(beats, beat{id: s.id, kind: k})
		}
	}
	m.mu.Unlock()
	for _, k := range toEnd {
		m.logf(fmt.Sprintf("Ending idle %s session (no activity for %d min).", k, int(sessClientIdle.Minutes())))
		m.End(k)
	}
	for _, b := range beats {
		go m.heartbeat(b.id)
	}
}

func (m *SessionManager) heartbeat(id int64) {
	jwt, _, ok := m.token()
	if !ok {
		return
	}
	m.postJSON(jwt, "/user/session/heartbeat", map[string]any{"sessionId": id}, nil)
}

func (m *SessionManager) postJSON(jwt, path string, body any, out any) bool {
	b, err := json.Marshal(body)
	if err != nil {
		return false
	}
	req, err := http.NewRequest(http.MethodPost, m.baseURL+path, bytes.NewReader(b))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("User-Agent", "AlbionMarketDataClient")
	if _, userID, ok := m.token(); ok && userID != "" {
		req.Header.Set("X-User-Id", userID)
	}
	resp, err := m.client.Do(req)
	if err != nil {
		m.logf("Session " + path + " failed: " + err.Error())
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		m.logf(fmt.Sprintf("Session %s -> %d", path, resp.StatusCode))
		return false
	}
	if out != nil {
		_ = json.NewDecoder(resp.Body).Decode(out)
	}
	return true
}
