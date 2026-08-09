package usersync

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// PrivateSettings mirrors GET /user/private-uploads. PrivateUploads is a
// pointer because the backend distinguishes "never set" (null) from an explicit
// false — the client only pushes its own value over a null.
type PrivateSettings struct {
	Enabled        bool  `json:"enabled"`
	PrivateUploads *bool `json:"private_uploads"`
	HoldMinutes    int   `json:"hold_minutes"`
}

// GetPrivateSettings reads the account-level private-uploads flag and the live
// hold period. Login-gated; ok=false when logged out or unreachable.
func (s *Syncer) GetPrivateSettings(ctx context.Context) (PrivateSettings, bool) {
	var out PrivateSettings
	jwt, userID, ok := s.token()
	if !ok {
		return out, false
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/user/private-uploads", nil)
	if err != nil {
		return out, false
	}
	req.Header.Set("Authorization", "Bearer "+jwt)
	if userID != "" {
		req.Header.Set("X-User-Id", userID)
	}
	if s.clientVersion != "" {
		req.Header.Set("X-Client-Version", s.clientVersion)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return out, false
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		io.Copy(io.Discard, resp.Body)
		return out, false
	}
	if json.NewDecoder(resp.Body).Decode(&out) != nil {
		return out, false
	}
	return out, true
}

// PutPrivateSettings pushes the account-level flag so the backend can fail
// closed if the per-upload header is ever stripped in transit.
func (s *Syncer) PutPrivateSettings(ctx context.Context, private bool) bool {
	jwt, _, ok := s.token()
	if !ok {
		return false
	}
	body, _ := json.Marshal(map[string]bool{"private_uploads": private})
	return s.postRaw(ctx, jwt, "/user/private-uploads", body)
}

// SyncPrivateSettings converges the client and the account on login and
// returns the value the client should hold, plus the live hold period in
// minutes (0 = unknown).
//
// Rule: **private wins** — effective = local OR account. It matches the
// backend's own semantics (it holds prices when either the header or the
// account flag says so), and it's the only rule that can't silently drop
// privacy the user asked for. A toggle made while logged out never reaches the
// account, so "adopt whatever the account says" would quietly undo it; taking
// the union can only ever err toward holding data back, never toward
// publishing it. Whichever side is behind gets updated.
func (s *Syncer) SyncPrivateSettings(local bool) (effective bool, holdMinutes int) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	got, ok := s.GetPrivateSettings(ctx)
	if !ok {
		return local, 0
	}
	account := got.PrivateUploads != nil && *got.PrivateUploads
	effective = local || account
	// push when the account is unset or behind (never seen it, or it says off
	// while this client says on)
	if got.PrivateUploads == nil || account != effective {
		s.PutPrivateSettings(ctx, effective)
	}
	return effective, got.HoldMinutes
}
