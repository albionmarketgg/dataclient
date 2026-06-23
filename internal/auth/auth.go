// Package auth implements optional "Login with Discord" via a backend-brokered
// device-authorization flow (browser + poll, RFC 8628 shape); the client never
// holds a Discord secret. Optional: with no AuthBaseURL it's a no-op and uploads
// stay anonymous. See DISCORD_AUTH_DESIGN.md.
package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// User is the authenticated user surfaced to the UI.
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// Verification is a character the account has verified.
type Verification struct {
	Server string `json:"server"`
	Name   string `json:"name"`
}

// VerifyResult is the outcome of a /auth/verify call.
type VerifyResult struct {
	Status        string         `json:"status"`
	Server        string         `json:"server"`
	Name          string         `json:"name"`
	Message       string         `json:"message"`
	Verifications []Verification `json:"verifications"`
}

// StoredAuth is the persisted login (refresh token is stored encrypted).
type StoredAuth struct {
	UserID       string
	RefreshToken string
	Username     string
	Avatar       string
}

// RefreshStore persists the rotating refresh token across restarts.
type RefreshStore interface {
	Save(StoredAuth) error
	Load() (StoredAuth, bool, error)
	Clear() error
}

// Provider is the auth surface used by the engine/uploader/UI.
type Provider interface {
	Configured() bool
	Token(ctx context.Context) (string, bool)
	User() (User, bool)
	Login(ctx context.Context) error
	Logout(ctx context.Context) error
}

const accessTTL = 14 * time.Minute // refresh a bit before the ~15m server TTL

// Discord is the backend-brokered Discord auth provider.
type Discord struct {
	authBase  string
	store     RefreshStore
	client    *http.Client
	openURL   func(string) error
	logf      func(string)
	onPending func(userCode, verificationURL string)

	mu          sync.Mutex
	accessToken string
	accessAt    time.Time
	refresh     string
	user        User
}

// New builds a Discord auth provider. Empty authBase disables auth. openURL
// opens a URL in the browser (inject runtime.BrowserOpenURL from Wails).
func New(authBase string, store RefreshStore, openURL func(string) error, logf func(string)) *Discord {
	if logf == nil {
		logf = func(string) {}
	}
	if openURL == nil {
		openURL = openInBrowser
	}
	return &Discord{
		authBase: strings.TrimRight(authBase, "/"),
		store:    store,
		client:   &http.Client{Timeout: 20 * time.Second},
		openURL:  openURL,
		logf:     logf,
	}
}

// SetOnPending registers a callback for when login awaits browser approval, so
// the UI can show a manual fallback (user code + URL).
func (d *Discord) SetOnPending(fn func(userCode, verificationURL string)) { d.onPending = fn }

// Configured reports whether auth is enabled.
func (d *Discord) Configured() bool { return d.authBase != "" }

// Restore attempts a silent login from a stored refresh token at startup.
func (d *Discord) Restore(ctx context.Context) {
	if !d.Configured() || d.store == nil {
		return
	}
	sa, ok, err := d.store.Load()
	if err != nil || !ok || sa.RefreshToken == "" {
		return
	}
	plain, err := unprotect(sa.RefreshToken)
	if err != nil {
		return
	}
	d.mu.Lock()
	d.refresh = plain
	d.user = User{ID: sa.UserID, Username: sa.Username, Avatar: sa.Avatar}
	d.mu.Unlock()
	if _, err := d.refreshAccess(ctx); err != nil {
		d.logf("Discord session expired; logged out.")
		d.clear()
	} else {
		d.logf("Signed in to Discord as " + sa.Username + ".")
	}
}

// User returns the current user, if logged in.
func (d *Discord) User() (User, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.user, d.user.ID != ""
}

// Token returns a valid access token, refreshing if needed.
func (d *Discord) Token(ctx context.Context) (string, bool) {
	if !d.Configured() {
		return "", false
	}
	d.mu.Lock()
	tok, at, refresh := d.accessToken, d.accessAt, d.refresh
	d.mu.Unlock()
	if tok != "" && time.Since(at) < accessTTL {
		return tok, true
	}
	if refresh == "" {
		return "", false
	}
	newTok, err := d.refreshAccess(ctx)
	if err != nil {
		return "", false
	}
	return newTok, true
}

func (d *Discord) refreshAccess(ctx context.Context) (string, error) {
	d.mu.Lock()
	refresh := d.refresh
	d.mu.Unlock()
	if refresh == "" {
		return "", errors.New("auth: no refresh token")
	}
	body, _ := json.Marshal(map[string]string{"refreshToken": refresh})
	var out struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
	if err := d.postJSON(ctx, "/auth/refresh", "", body, &out); err != nil {
		return "", err
	}
	d.mu.Lock()
	d.accessToken = out.AccessToken
	d.accessAt = time.Now()
	if out.RefreshToken != "" {
		d.refresh = out.RefreshToken
	}
	user := d.user
	refreshTok := d.refresh
	d.mu.Unlock()
	d.persist(user, refreshTok)
	return out.AccessToken, nil
}

type deviceStart struct {
	DeviceCode              string `json:"deviceCode"`
	UserCode                string `json:"userCode"`
	VerificationURL         string `json:"verificationUrl"`
	VerificationURLComplete string `json:"verificationUrlComplete"`
	Interval                int    `json:"interval"`
	ExpiresIn               int    `json:"expiresIn"`
}

// Login runs the device-authorization flow: open a browser to our backend and
// poll until the user approves the sign-in there.
func (d *Discord) Login(ctx context.Context) error {
	if !d.Configured() {
		return errors.New("auth: not configured")
	}
	startBody, _ := json.Marshal(map[string]string{"client": "AlbionMarketDataClient", "version": "0.1.0"})
	var start deviceStart
	if err := d.postJSON(ctx, "/auth/device/start", "", startBody, &start); err != nil {
		return err
	}
	if start.DeviceCode == "" {
		return errors.New("auth: backend returned no device code")
	}

	openTarget := start.VerificationURLComplete
	if openTarget == "" {
		openTarget = start.VerificationURL
	}
	if err := d.openURL(openTarget); err != nil {
		d.logf("Couldn't open browser. Visit " + start.VerificationURL + " and enter code " + start.UserCode)
	}
	if d.onPending != nil {
		d.onPending(start.UserCode, start.VerificationURL)
	}
	d.logf("Waiting for Discord sign-in in your browser…")

	interval := start.Interval
	if interval <= 0 {
		interval = 5
	}
	expires := start.ExpiresIn
	if expires <= 0 {
		expires = 600
	}
	deadline := time.Now().Add(time.Duration(expires) * time.Second)
	pollBody, _ := json.Marshal(map[string]string{"deviceCode": start.DeviceCode})

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(interval) * time.Second):
		}
		if time.Now().After(deadline) {
			return errors.New("auth: login timed out")
		}

		var poll struct {
			Status       string `json:"status"`
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
			User         User   `json:"user"`
		}
		if err := d.postJSON(ctx, "/auth/device/poll", "", pollBody, &poll); err != nil {
			continue // transient; keep polling until the deadline
		}
		switch poll.Status {
		case "pending", "":
			continue
		case "slow_down":
			interval += 2
		case "complete":
			d.adopt(poll.AccessToken, poll.RefreshToken, poll.User)
			d.logf("Signed in to Discord as " + poll.User.Username + ".")
			return nil
		case "denied":
			return errors.New("auth: login was denied")
		case "expired":
			return errors.New("auth: login request expired")
		}
	}
}

// adopt installs freshly minted tokens + user and persists the refresh token.
func (d *Discord) adopt(access, refresh string, user User) {
	d.mu.Lock()
	d.accessToken = access
	d.accessAt = time.Now()
	if refresh != "" {
		d.refresh = refresh
	}
	r := d.refresh
	d.user = user
	d.mu.Unlock()
	d.persist(user, r)
}

// Logout revokes the session server-side and clears local state.
func (d *Discord) Logout(ctx context.Context) error {
	d.mu.Lock()
	refresh := d.refresh
	d.mu.Unlock()
	if d.Configured() && refresh != "" {
		body, _ := json.Marshal(map[string]string{"refreshToken": refresh})
		_ = d.postJSON(ctx, "/auth/logout", "", body, nil)
	}
	d.clear()
	d.logf("Signed out of Discord.")
	return nil
}

func (d *Discord) clear() {
	d.mu.Lock()
	d.accessToken, d.refresh, d.user = "", "", User{}
	d.mu.Unlock()
	if d.store != nil {
		_ = d.store.Clear()
	}
}

// Verify submits a detected character to /auth/verify. Returns the parsed result
// for ALL status codes (the body carries status/message even on 4xx).
func (d *Discord) Verify(ctx context.Context, serverID int, name string, fame int64) (VerifyResult, error) {
	tok, ok := d.Token(ctx)
	if !ok {
		return VerifyResult{Status: "unauthorized", Message: "Sign in with Discord first."}, nil
	}
	body, _ := json.Marshal(map[string]any{"serverId": serverID, "ingameName": name, "fame": fame})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.authBase+"/auth/verify", strings.NewReader(string(body)))
	if err != nil {
		return VerifyResult{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tok)
	resp, err := d.client.Do(req)
	if err != nil {
		return VerifyResult{}, err
	}
	defer resp.Body.Close()
	var out VerifyResult
	json.NewDecoder(resp.Body).Decode(&out)
	if out.Status == "" {
		out.Status = "error"
		if out.Message == "" {
			out.Message = fmt.Sprintf("Verification failed (%d).", resp.StatusCode)
		}
	}
	return out, nil
}

func (d *Discord) persist(user User, refresh string) {
	if d.store == nil || refresh == "" {
		return
	}
	enc, err := protect(refresh)
	if err != nil {
		return
	}
	_ = d.store.Save(StoredAuth{UserID: user.ID, RefreshToken: enc, Username: user.Username, Avatar: user.Avatar})
}

func (d *Discord) postJSON(ctx context.Context, path, bearer string, body []byte, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.authBase+path, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("auth: %s returned %d", path, resp.StatusCode)
	}
	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}
