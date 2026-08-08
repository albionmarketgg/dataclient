package upload

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/albionmarketgg/dataclient/internal/config"
	"github.com/albionmarketgg/dataclient/internal/market"
)

type capturedReq struct {
	path    string
	auth    string
	userID  string
	private string
}

// recorder serves the PoW handshake and records every upload request.
func recorder(t *testing.T) (*httptest.Server, func() []capturedReq) {
	t.Helper()
	var mu sync.Mutex
	var got []capturedReq
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/pow" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"key":"abc","wanted":"0"}`)) // 0 difficulty = instant solve
			return
		}
		mu.Lock()
		got = append(got, capturedReq{
			path:    r.URL.Path,
			auth:    r.Header.Get("Authorization"),
			userID:  r.Header.Get("X-User-Id"),
			private: r.Header.Get("X-Private"),
		})
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	return srv, func() []capturedReq {
		mu.Lock()
		defer mu.Unlock()
		out := make([]capturedReq, len(got))
		copy(out, got)
		return out
	}
}

func waitFor(t *testing.T, get func() []capturedReq, n int) []capturedReq {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if r := get(); len(r) >= n {
			return r
		}
		time.Sleep(10 * time.Millisecond)
	}
	return get()
}

func testCfg(base string) config.Config {
	cfg := config.Default()
	cfg.IngestBaseURL = base
	cfg.RequirePoW = true
	return cfg
}

// X-Private is sent only while the private switch is on.
func TestPrivateHeader(t *testing.T) {
	srv, got := recorder(t)
	defer srv.Close()

	private := false
	u := New(testCfg(srv.URL), nil)
	u.SetPrivateProvider(func() bool { return private })
	u.Start()
	defer u.Stop()

	u.EnqueueGold(market.GoldPriceUpload{Prices: []uint32{1}, Timestamps: []int64{1}})
	if r := waitFor(t, got, 1); len(r) != 1 || r[0].private != "" {
		t.Fatalf("public upload should not be flagged: %+v", r)
	}
	private = true
	u.EnqueueGold(market.GoldPriceUpload{Prices: []uint32{2}, Timestamps: []int64{2}})
	r := waitFor(t, got, 2)
	if len(r) != 2 || r[1].private != "1" {
		t.Fatalf("private upload not flagged: %+v", r)
	}
}

// The AODP mirror must receive market data WITHOUT any account identity, and
// must never be fed while private uploads are on.
func TestAodpMirrorIsAnonymousAndInterlocked(t *testing.T) {
	ours, ourReqs := recorder(t)
	defer ours.Close()
	theirs, theirReqs := recorder(t)
	defer theirs.Close()

	private, contribute := false, true
	u := New(testCfg(ours.URL), nil)
	u.SetTokenProvider(func() (string, string, bool) { return "jwt-secret", "discord-123", true })
	u.SetPrivateProvider(func() bool { return private })

	mirror := NewAnonymous(testCfg(theirs.URL), func() string {
		if private || !contribute {
			return ""
		}
		return theirs.URL
	}, nil)
	mirror.SetTokenProvider(func() (string, string, bool) { return "jwt-secret", "discord-123", true }) // must be ignored
	u.SetMirror(mirror, func() bool { return contribute && !private })

	u.Start()
	mirror.Start()
	defer u.Stop()
	defer mirror.Stop()

	u.EnqueueGold(market.GoldPriceUpload{Prices: []uint32{1}, Timestamps: []int64{1}})

	ourGot := waitFor(t, ourReqs, 1)
	if len(ourGot) != 1 || ourGot[0].auth == "" || ourGot[0].userID != "discord-123" {
		t.Fatalf("our upload lost its identity: %+v", ourGot)
	}
	mirrorGot := waitFor(t, theirReqs, 1)
	if len(mirrorGot) != 1 {
		t.Fatalf("mirror got %d uploads, want 1", len(mirrorGot))
	}
	if mirrorGot[0].auth != "" || mirrorGot[0].userID != "" || mirrorGot[0].private != "" {
		t.Fatalf("account identity leaked to third party: %+v", mirrorGot[0])
	}

	// private on -> mirror must receive nothing further
	private = true
	u.EnqueueGold(market.GoldPriceUpload{Prices: []uint32{2}, Timestamps: []int64{2}})
	waitFor(t, ourReqs, 2) // ours still uploads (flagged private)
	time.Sleep(200 * time.Millisecond)
	if r := theirReqs(); len(r) != 1 {
		t.Fatalf("mirror uploaded while private: %d requests", len(r))
	}
}

// With contribution off, only our endpoint is used.
func TestMirrorDisabled(t *testing.T) {
	ours, ourReqs := recorder(t)
	defer ours.Close()
	theirs, theirReqs := recorder(t)
	defer theirs.Close()

	u := New(testCfg(ours.URL), nil)
	mirror := NewAnonymous(testCfg(theirs.URL), func() string { return theirs.URL }, nil)
	u.SetMirror(mirror, func() bool { return false })
	u.Start()
	mirror.Start()
	defer u.Stop()
	defer mirror.Stop()

	u.EnqueueMarket(market.Upload{Orders: []market.Order{{ID: 1}}})
	waitFor(t, ourReqs, 1)
	time.Sleep(200 * time.Millisecond)
	if r := theirReqs(); len(r) != 0 {
		t.Fatalf("disabled mirror still uploaded: %+v", r)
	}
}
