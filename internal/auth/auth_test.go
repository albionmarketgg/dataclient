package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func TestProtectRoundTrip(t *testing.T) {
	enc, err := protect("secret-refresh-token")
	if err != nil {
		t.Fatal(err)
	}
	dec, err := unprotect(enc)
	if err != nil {
		t.Fatal(err)
	}
	if dec != "secret-refresh-token" {
		t.Fatalf("roundtrip: %q", dec)
	}
}

type memStore struct {
	mu sync.Mutex
	sa *StoredAuth
}

func (m *memStore) Save(s StoredAuth) error { m.mu.Lock(); defer m.mu.Unlock(); cp := s; m.sa = &cp; return nil }
func (m *memStore) Load() (StoredAuth, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.sa == nil {
		return StoredAuth{}, false, nil
	}
	return *m.sa, true, nil
}
func (m *memStore) Clear() error { m.mu.Lock(); defer m.mu.Unlock(); m.sa = nil; return nil }

// stubBackend implements the device-authorization flow: start -> a couple of
// pending polls -> complete; plus refresh + logout.
func stubBackend() *httptest.Server {
	const deviceCode = "device-secret-abcdef0123456789"
	var polls int32
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/device/start", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"deviceCode":              deviceCode,
			"userCode":                "ABCD-1234",
			"verificationUrl":         "https://albionmarket.gg/link",
			"verificationUrlComplete": "https://albionmarket.gg/link?code=ABCD-1234",
			"interval":                1, // 1s polls for a fast test
			"expiresIn":               60,
		})
	})
	mux.HandleFunc("/auth/device/poll", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			DeviceCode string `json:"deviceCode"`
		}
		json.NewDecoder(r.Body).Decode(&in)
		if in.DeviceCode != deviceCode {
			http.Error(w, "bad device code", http.StatusForbidden)
			return
		}
		n := atomic.AddInt32(&polls, 1)
		if n < 2 {
			json.NewEncoder(w).Encode(map[string]string{"status": "pending"})
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"status":       "complete",
			"accessToken":  "access-1",
			"refreshToken": "refresh-1",
			"user":         User{ID: "u123", Username: "Tester", Avatar: "av"},
		})
	})
	mux.HandleFunc("/auth/refresh", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			RefreshToken string `json:"refreshToken"`
		}
		json.NewDecoder(r.Body).Decode(&in)
		if !strings.HasPrefix(in.RefreshToken, "refresh-") {
			http.Error(w, "bad refresh", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"accessToken": "access-2", "refreshToken": "refresh-2"})
	})
	mux.HandleFunc("/auth/logout", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	return httptest.NewServer(mux)
}

func TestDeviceLoginAndRefresh(t *testing.T) {
	srv := stubBackend()
	defer srv.Close()

	store := &memStore{}
	var opened atomic.Bool
	open := func(u string) error { opened.Store(true); return nil }

	d := New(srv.URL, store, open, nil)
	if !d.Configured() {
		t.Fatal("should be configured")
	}
	var pendingSeen atomic.Bool
	d.SetOnPending(func(code, url string) {
		if code == "ABCD-1234" {
			pendingSeen.Store(true)
		}
	})

	if err := d.Login(context.Background()); err != nil {
		t.Fatalf("login: %v", err)
	}
	if !opened.Load() {
		t.Fatal("browser was not opened")
	}
	if !pendingSeen.Load() {
		t.Fatal("onPending not fired")
	}
	u, ok := d.User()
	if !ok || u.ID != "u123" || u.Username != "Tester" {
		t.Fatalf("user: %+v ok=%v", u, ok)
	}
	tok, ok := d.Token(context.Background())
	if !ok || tok != "access-1" {
		t.Fatalf("token: %q ok=%v", tok, ok)
	}
	if store.sa == nil || store.sa.RefreshToken == "refresh-1" {
		t.Fatalf("refresh token not encrypted at rest: %+v", store.sa)
	}

	// new provider restores from store and refreshes silently
	d2 := New(srv.URL, store, open, nil)
	d2.Restore(context.Background())
	if u2, ok := d2.User(); !ok || u2.ID != "u123" {
		t.Fatalf("restore user: %+v", u2)
	}
	if tok2, ok := d2.Token(context.Background()); !ok || tok2 != "access-2" {
		t.Fatalf("refreshed token: %q ok=%v", tok2, ok)
	}

	d2.Logout(context.Background())
	if _, ok := d2.User(); ok {
		t.Fatal("still logged in after logout")
	}
}
