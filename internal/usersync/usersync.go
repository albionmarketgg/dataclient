// Package usersync uploads the user's own captured data to their Albion Market
// account, gated on login AND per-type consent. Trades/mails sync idempotently by
// id; gathering/loot/party as throttled snapshots; dungeon per-event. See
// USER_DATA_UPLOAD_SPEC.md.
package usersync

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/niick1231/albionmarket_dataclient/internal/store"
)

// Snapshots provides current in-memory tracker state as JSON-marshalable values.
// ServerID (0 if unknown) is merged into each snapshot body to label sessions.
type Snapshots struct {
	Gathering  func() any
	Loot       func() any
	Party      func() any
	Awakened   func() any
	ServerID   func() int
	PlayerName func() string // in-game character that produced the data
}

// Stats reports how many items of each kind have been uploaded this session.
type Stats struct {
	Trades    int64 `json:"trades"`
	Mails     int64 `json:"mails"`
	Gathering int64 `json:"gathering"`
	Dungeon   int64 `json:"dungeon"`
	Loot      int64 `json:"loot"`
	Party     int64 `json:"party"`
	Specs     int64 `json:"specs"`
	Awakened  int64 `json:"awakened"`
}

// SpecEntry is one Destiny Board node + mastery level, for the specs upload.
type SpecEntry struct {
	ID    string `json:"id"`
	Level int    `json:"level"`
}

// Syncer periodically uploads user-owned data when logged in and opted in.
type Syncer struct {
	baseURL string
	store   *store.Store
	token   func() (jwt, userID string, ok bool)
	enabled func(kind string) bool // per-type consent (kind: trades/mails/gathering/combat/loot/party)
	snaps   Snapshots
	client  *http.Client
	logf    func(string)

	cancel context.CancelFunc
	wg     sync.WaitGroup

	sessions *SessionManager // tags uploads with the active sessionId (may be nil)

	awMu    sync.Mutex
	awTimer *time.Timer // debounce for event-driven awakened sync

	mu       sync.Mutex
	lastSig  map[string]string
	counters map[string]*atomic.Int64
}

// SetSessions attaches the session manager so uploads can carry their sessionId.
func (s *Syncer) SetSessions(sm *SessionManager) { s.sessions = sm }

// New builds a Syncer. baseURL is the account API origin; token returns the
// bearer JWT + user id (ok=false when logged out); enabled reports per-type consent.
func New(baseURL string, st *store.Store, token func() (string, string, bool), enabled func(kind string) bool, snaps Snapshots, logf func(string)) *Syncer {
	if logf == nil {
		logf = func(string) {}
	}
	counters := map[string]*atomic.Int64{}
	for _, k := range []string{"trades", "mails", "gathering", "dungeon", "loot", "party", "specs", "awakened"} {
		counters[k] = &atomic.Int64{}
	}
	return &Syncer{
		baseURL:  strings.TrimRight(baseURL, "/"),
		store:    st,
		token:    token,
		enabled:  enabled,
		snaps:    snaps,
		client:   &http.Client{Timeout: 20 * time.Second},
		logf:     logf,
		lastSig:  map[string]string{},
		counters: counters,
	}
}

// Stats returns per-kind upload counts for this session.
func (s *Syncer) Stats() Stats {
	return Stats{
		Trades:    s.counters["trades"].Load(),
		Mails:     s.counters["mails"].Load(),
		Gathering: s.counters["gathering"].Load(),
		Dungeon:   s.counters["dungeon"].Load(),
		Loot:      s.counters["loot"].Load(),
		Party:     s.counters["party"].Load(),
		Specs:     s.counters["specs"].Load(),
		Awakened:  s.counters["awakened"].Load(),
	}
}

func (s *Syncer) on(kind string) bool { return s.enabled != nil && s.enabled(kind) }

// Start runs the periodic sync loop (first pass shortly after start).
func (s *Syncer) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		first := time.NewTimer(12 * time.Second)
		defer first.Stop()
		select {
		case <-ctx.Done():
			return
		case <-first.C:
			s.syncOnce(ctx)
		}
		t := time.NewTicker(60 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				s.syncOnce(ctx)
			}
		}
	}()
}

// Stop halts the sync loop.
func (s *Syncer) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	s.awMu.Lock()
	if s.awTimer != nil {
		s.awTimer.Stop()
	}
	s.awMu.Unlock()
	s.wg.Wait()
}

func (s *Syncer) syncOnce(ctx context.Context) {
	jwt, _, ok := s.token() // logged out -> nothing
	if !ok {
		return
	}
	if s.store != nil {
		if s.on("trades") {
			s.syncTrades(ctx, jwt)
		}
		if s.on("mails") {
			s.syncMails(ctx, jwt)
		}
	}
	if s.on("loot") {
		s.syncSnapshot(ctx, jwt, "loot", s.snaps.Loot)
	}
	if s.on("party") {
		s.syncSnapshot(ctx, jwt, "party", s.snaps.Party)
	}
	if s.on("gathering") {
		s.syncSnapshot(ctx, jwt, "gathering", s.snaps.Gathering)
	}
	// awakened is event-driven (see TriggerAwakened), not on this loop.
}

// TriggerAwakened schedules an awakened sync ~2s after the latest inventory
// change (debounced), so a burst of item packets results in a single push. The
// push itself is still change-deduped, so identical data won't re-POST.
func (s *Syncer) TriggerAwakened() {
	s.awMu.Lock()
	defer s.awMu.Unlock()
	if s.awTimer != nil {
		s.awTimer.Stop()
	}
	s.awTimer = time.AfterFunc(2*time.Second, func() {
		jwt, _, ok := s.token()
		if !ok || !s.on("awakened") {
			return
		}
		s.syncAwakened(context.Background(), jwt)
	})
}

// syncAwakened pushes the awakened inventory to /user/awakened/sync (login-gated,
// only when the set changed since the last push). serverId is merged top-level.
func (s *Syncer) syncAwakened(ctx context.Context, jwt string) {
	if s.snaps.Awakened == nil {
		return
	}
	body, err := json.Marshal(s.snaps.Awakened())
	if err != nil {
		return
	}
	sid := 0
	if s.snaps.ServerID != nil {
		sid = s.snaps.ServerID()
	}
	if sid < 1 || sid > 3 {
		return // server not detected yet — a serverId of 0 would 400
	}
	var m map[string]any
	if json.Unmarshal(body, &m) == nil && m != nil {
		m["serverId"] = sid
		if merged, e := json.Marshal(m); e == nil {
			body = merged
		}
	}
	sig := hashBytes(body)
	s.mu.Lock()
	same := s.lastSig["awakened"] == sig
	s.mu.Unlock()
	if same {
		return
	}
	if s.postRaw(ctx, jwt, "/user/awakened/sync", body) {
		s.mu.Lock()
		s.lastSig["awakened"] = sig
		s.mu.Unlock()
		s.counters["awakened"].Add(1)
	}
}

// ---- wire shapes (camelCase) ----

type tradeDTO struct {
	ID          string  `json:"id"`
	ServerID    int     `json:"serverId"`
	ItemID      string  `json:"itemId"`
	Quality     int     `json:"quality"`
	LocationID  string  `json:"locationId"`
	Operation   string  `json:"operation"` // buy|sell
	TradeType   string  `json:"tradeType"` // instant|order
	Amount      int     `json:"amount"`
	UnitSilver  float64 `json:"unitSilver"` // already divided (display silver)
	SalesTaxPct float64 `json:"salesTaxPct"`
	TradedAt    string  `json:"tradedAt"`
	PlayerName  string  `json:"playerName"`
	Deleted     bool    `json:"deleted"`
}

type mailDTO struct {
	ID            int64   `json:"id"`
	ServerID      int     `json:"serverId"`
	AuctionType   int     `json:"auctionType"`
	ItemID        string  `json:"itemId"`
	LocationID    string  `json:"locationId"`
	PartialAmount int     `json:"partialAmount"`
	TotalAmount   int     `json:"totalAmount"`
	ReceivedAt    string  `json:"receivedAt"`
	TaxesPct      float64 `json:"taxesPct"`
	TotalSilver   int64   `json:"totalSilver"` // already divided
	TotalTaxes    int64   `json:"totalTaxes"`
	UnitSilver    float64 `json:"unitSilver"`
	Type          int     `json:"type"`
	PlayerName    string  `json:"playerName"`
	Deleted       bool    `json:"deleted"`
}

func (s *Syncer) syncTrades(ctx context.Context, jwt string) {
	rows, err := s.store.UnsyncedTrades(500)
	if err != nil || len(rows) == 0 {
		return
	}
	out := make([]tradeDTO, 0, len(rows))
	ids := make([]string, 0, len(rows))
	for _, t := range rows {
		op := "buy"
		if t.Operation == 1 {
			op = "sell"
		}
		tt := "instant"
		if t.Type == 1 {
			tt = "order"
		}
		out = append(out, tradeDTO{
			ID: t.ID, ServerID: t.AlbionServerID, ItemID: t.ItemID, Quality: t.QualityLevel,
			LocationID: t.RawLocationID, Operation: op, TradeType: tt, Amount: t.Amount,
			UnitSilver: t.UnitSilver, SalesTaxPct: t.SalesTaxesPercent,
			TradedAt: t.DateTime.UTC().Format(time.RFC3339), PlayerName: t.PlayerName, Deleted: t.Deleted,
		})
		ids = append(ids, t.ID)
	}
	if s.post(ctx, jwt, "/user/trades", map[string]any{"trades": out}) {
		if err := s.store.MarkTradesSynced(ids); err == nil {
			s.counters["trades"].Add(int64(len(ids)))
			s.logf(fmt.Sprintf("Synced %d trade(s) to your account.", len(ids)))
		}
	}
}

func (s *Syncer) syncMails(ctx context.Context, jwt string) {
	rows, err := s.store.UnsyncedMails(500)
	if err != nil || len(rows) == 0 {
		return
	}
	out := make([]mailDTO, 0, len(rows))
	ids := make([]int64, 0, len(rows))
	for _, m := range rows {
		out = append(out, mailDTO{
			ID: m.ID, ServerID: m.AlbionServerID, AuctionType: m.AuctionType, ItemID: m.ItemID,
			LocationID: m.RawLocationID, PartialAmount: m.PartialAmount, TotalAmount: m.TotalAmount,
			ReceivedAt: m.Received.UTC().Format(time.RFC3339), TaxesPct: m.TaxesPercent,
			TotalSilver: m.TotalSilver, TotalTaxes: m.TotalTaxes, UnitSilver: m.UnitSilver,
			Type: m.Type, PlayerName: m.PlayerName, Deleted: m.Deleted,
		})
		ids = append(ids, m.ID)
	}
	if s.post(ctx, jwt, "/user/mails", map[string]any{"mails": out}) {
		if err := s.store.MarkMailsSynced(ids); err == nil {
			s.counters["mails"].Add(int64(len(ids)))
			s.logf(fmt.Sprintf("Synced %d mail(s) to your account.", len(ids)))
		}
	}
}

// UploadDungeonEvent posts a single dungeon gain (one fame event + nearest silver)
// immediately, when logged in and the dungeon kind is enabled. fame<=0 is never
// sent, so silver-only activity produces nothing.
func (s *Syncer) UploadDungeonEvent(fame, silver int64, at time.Time) {
	if fame <= 0 || !s.on("dungeon") {
		return
	}
	jwt, _, ok := s.token()
	if !ok {
		return
	}
	body := map[string]any{
		"fame":   fame,
		"silver": silver,
		"at":     at.UTC().Format(time.RFC3339),
	}
	if s.snaps.ServerID != nil {
		body["serverId"] = s.snaps.ServerID()
	}
	if s.snaps.PlayerName != nil {
		body["playerName"] = s.snaps.PlayerName()
	}
	if s.sessions != nil {
		if id := s.sessions.CurrentID("dungeon"); id != 0 {
			body["sessionId"] = id
		}
	}
	if s.post(context.Background(), jwt, "/user/dungeon", body) {
		s.counters["dungeon"].Add(1)
		if s.sessions != nil {
			s.sessions.MarkActivity("dungeon")
		}
	}
}

// UploadSpecs posts a character's full Destiny Board immediately, when logged in
// and specs upload is enabled. Login-gated only (no session). playerName is the
// in-game character the board belongs to.
func (s *Syncer) UploadSpecs(serverID int, playerName, version string, entries []SpecEntry) {
	if len(entries) == 0 || playerName == "" || !s.on("specs") {
		return
	}
	jwt, _, ok := s.token()
	if !ok {
		return
	}
	body := map[string]any{
		"serverId":            serverID,
		"playerName":          playerName,
		"achievementsVersion": version,
		"specs":               entries,
	}
	if s.post(context.Background(), jwt, "/user/specs", body) {
		s.counters["specs"].Add(1)
		s.logf(fmt.Sprintf("Synced %d character spec(s) for %s to your account.", len(entries), playerName))
	}
}

// syncSnapshot uploads an in-memory tracker snapshot, but only when it changed
// since the last upload (avoids a per-tick firehose).
func (s *Syncer) syncSnapshot(ctx context.Context, jwt, name string, fn func() any) {
	if fn == nil {
		return
	}
	snap := fn()
	body, err := json.Marshal(snap)
	if err != nil {
		return
	}
	// merge serverId, in-game character, and sessionId into the body for attribution.
	var sessionID int64
	if s.sessions != nil {
		sessionID = s.sessions.CurrentID(name)
	}
	if s.snaps.ServerID != nil || s.snaps.PlayerName != nil || sessionID != 0 {
		var m map[string]any
		if json.Unmarshal(body, &m) == nil && m != nil {
			if s.snaps.ServerID != nil {
				m["serverId"] = s.snaps.ServerID()
			}
			if s.snaps.PlayerName != nil {
				m["playerName"] = s.snaps.PlayerName()
			}
			if sessionID != 0 {
				m["sessionId"] = sessionID
			}
			if merged, e := json.Marshal(m); e == nil {
				body = merged
			}
		}
	}
	sig := hashBytes(body)
	s.mu.Lock()
	same := s.lastSig[name] == sig
	s.mu.Unlock()
	if same {
		return
	}
	if s.postRaw(ctx, jwt, "/user/"+name, body) {
		s.mu.Lock()
		s.lastSig[name] = sig
		s.mu.Unlock()
		if c := s.counters[name]; c != nil {
			c.Add(1)
		}
		if s.sessions != nil {
			s.sessions.MarkActivity(name)
		}
	}
}

func (s *Syncer) post(ctx context.Context, jwt, path string, payload any) bool {
	body, err := json.Marshal(payload)
	if err != nil {
		return false
	}
	return s.postRaw(ctx, jwt, path, body)
}

func (s *Syncer) postRaw(ctx context.Context, jwt, path string, body []byte) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("User-Agent", "AlbionMarketDataClient")
	if _, userID, ok := s.token(); ok && userID != "" {
		req.Header.Set("X-User-Id", userID)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		s.logf(fmt.Sprintf("User sync %s -> %d", path, resp.StatusCode))
		return false
	}
	return true
}

func hashBytes(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:8])
}
