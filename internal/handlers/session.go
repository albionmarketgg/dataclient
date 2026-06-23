package handlers

import (
	"strconv"
	"time"

	"github.com/niick1231/albionmarket_dataclient/internal/dispatch"
	"github.com/niick1231/albionmarket_dataclient/internal/locations"
	"github.com/niick1231/albionmarket_dataclient/internal/photon"
	"github.com/niick1231/albionmarket_dataclient/internal/state"
)

// Session handles the Join response, establishing the local player's object id,
// character name, and location. Market uploads are gated on a known location, so
// without this location and player name stay unknown during live capture.
type Session struct {
	st   *state.State
	feed FeedFunc
}

// NewSession builds the session handler.
func NewSession(st *state.State, feed FeedFunc) *Session {
	if feed == nil {
		feed = func(CaptureEvent) {}
	}
	return &Session{st: st, feed: feed}
}

// Register attaches the Join handler to the dispatcher.
func (s *Session) Register(d *dispatch.Dispatcher) {
	d.OnResponse(photon.OpJoin, s.handleJoin)
}

func (s *Session) handleJoin(_ int16, _ string, p map[byte]any) {
	if id, ok := dispatch.Int64(p[0]); ok {
		s.st.SetUserObjectID(id)
	}
	if name, ok := dispatch.String(p[2]); ok && name != "" {
		s.st.SetPlayerName(name)
		// Join carries the player's own fame at param[35] (raw, /10000) — record
		// this detected character for in-client verification.
		fame, _ := dispatch.Int64(p[35])
		s.st.RecordCharacter(name, fame/10000)
	}
	raw := locationParam(p[8])
	if raw != "" {
		name, _ := locations.Resolve(raw)
		s.st.SetLocation(state.Location{ID: raw, Name: name})
		s.feed(CaptureEvent{Time: time.Now(), Kind: "session", Detail: "entered " + name, Count: 1})
	}
}

// locationParam coerces the Join location parameter (string or numeric) to a string.
func locationParam(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	if n, ok := dispatch.Int64(v); ok {
		return strconv.FormatInt(n, 10)
	}
	return ""
}
