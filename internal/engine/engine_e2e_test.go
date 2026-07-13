package engine

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/albionmarketgg/dataclient/internal/config"
	"github.com/albionmarketgg/dataclient/internal/photon"
	"github.com/albionmarketgg/dataclient/internal/phototest"
	"github.com/albionmarketgg/dataclient/internal/pow"
	"github.com/albionmarketgg/dataclient/internal/state"
)

// mockIngest is an in-test ingest server with the PoW handshake.
type mockIngest struct {
	mu       sync.Mutex
	issued   map[string]string
	received []string // topics
	bodies   []string
}

func newMockIngest() *mockIngest { return &mockIngest{issued: map[string]string{}} }

func (m *mockIngest) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/pow", func(w http.ResponseWriter, r *http.Request) {
		key := "key-1234567890abcdef"
		wanted := "00000000" // 8 bits
		m.mu.Lock()
		m.issued[key] = wanted
		m.mu.Unlock()
		json.NewEncoder(w).Encode(pow.Request{Key: key, Wanted: wanted})
	})
	mux.HandleFunc("/pow/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		key := r.Form.Get("key")
		m.mu.Lock()
		wanted, ok := m.issued[key]
		m.mu.Unlock()
		if !ok || !pow.Verify(pow.Request{Key: key, Wanted: wanted}, r.Form.Get("solution")) {
			http.Error(w, "bad pow", http.StatusForbidden)
			return
		}
		m.mu.Lock()
		m.received = append(m.received, strings.TrimPrefix(r.URL.Path, "/pow/"))
		m.bodies = append(m.bodies, r.Form.Get("natsmsg"))
		m.mu.Unlock()
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/itemEstimatedMarketValues", func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		m.mu.Lock()
		m.received = append(m.received, "emv")
		m.bodies = append(m.bodies, string(b))
		m.mu.Unlock()
		w.WriteHeader(http.StatusOK)
	})
	return mux
}

func (m *mockIngest) topics() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]string, len(m.received))
	copy(out, m.received)
	return out
}

func waitFor(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatal("condition not met within timeout")
}

func TestEndToEndMarketUpload(t *testing.T) {
	mock := newMockIngest()
	srv := httptest.NewServer(mock.handler())
	defer srv.Close()

	cfg := config.Default()
	cfg.IngestBaseURL = srv.URL
	e := New(cfg, nil, "")
	e.Up.Start()
	defer e.Up.Stop()

	// make OkToUpload pass
	e.State.SetServer(&config.Server{ID: 1, Name: "Americas"})
	e.State.SetLocation(state.Location{ID: "3005", Name: "Caerleon"})

	order := `{"Id":1,"ItemTypeId":"T4_BAG","LocationId":"3005","QualityLevel":1,"UnitPriceSilver":12340000,"Amount":3,"AuctionType":"offer","Expires":"2026-07-01T00:00:00"}`
	pkt := phototest.ResponsePacket(byte(photon.OpAuctionGetOffers), 0, map[byte]any{
		253: int16(photon.OpAuctionGetOffers),
		0:   []string{order},
	})
	if st := e.Inject(pkt); st != photon.StatusSuccess {
		t.Fatalf("inject status %v", st)
	}

	waitFor(t, func() bool {
		for _, top := range mock.topics() {
			if top == cfg.MarketOrdersTopic {
				return true
			}
		}
		return false
	})

	// verify body round-trips an order
	mock.mu.Lock()
	body := mock.bodies[len(mock.bodies)-1]
	mock.mu.Unlock()
	if !strings.Contains(body, "T4_BAG") {
		t.Fatalf("uploaded body missing order: %s", body)
	}
}

func TestEndToEndGoldUpload(t *testing.T) {
	mock := newMockIngest()
	srv := httptest.NewServer(mock.handler())
	defer srv.Close()

	cfg := config.Default()
	cfg.IngestBaseURL = srv.URL
	e := New(cfg, nil, "")
	e.Up.Start()
	defer e.Up.Stop()
	e.State.SetServer(&config.Server{ID: 3, Name: "Europe"})
	e.State.SetLocation(state.Location{ID: "3005"})

	pkt := phototest.ResponsePacket(byte(photon.OpGoldMarketGetAverageInfo), 0, map[byte]any{
		253: int16(photon.OpGoldMarketGetAverageInfo),
		0:   []int32{4500, 4510, 4490},
		1:   []int64{1719800000, 1719800060, 1719800120},
	})
	e.Inject(pkt)

	waitFor(t, func() bool {
		for _, top := range mock.topics() {
			if top == cfg.GoldPricesTopic {
				return true
			}
		}
		return false
	})
}
