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

// SyncPrivateSettings converges the client and the account on login.
//
// The account flag is shared across devices and the website, and effective
// privacy is (header OR account flag) — so if the account already says private
// but this client doesn't, the client UI would be lying about what happens.
// Hence: adopt an account value that's already set, and only push ours when the
// backend has never been told (null). Returns the value the client should now
// hold, plus the live hold period in minutes (0 = unknown).
func (s *Syncer) SyncPrivateSettings(local bool) (effective bool, holdMinutes int) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	got, ok := s.GetPrivateSettings(ctx)
	if !ok {
		return local, 0
	}
	if got.PrivateUploads == nil {
		s.PutPrivateSettings(ctx, local) // never set -> seed it from this client
		return local, got.HoldMinutes
	}
	return *got.PrivateUploads, got.HoldMinutes
}
