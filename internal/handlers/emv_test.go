package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/albionmarketgg/dataclient/internal/config"
	"github.com/albionmarketgg/dataclient/internal/state"
	"github.com/albionmarketgg/dataclient/internal/upload"
)

type fakeNamer struct{}

func (fakeNamer) UniqueName(id int) (string, bool) { return "ITEM_" + strconv.Itoa(id), true }

func TestItemEMVUpload(t *testing.T) {
	var mu sync.Mutex
	var bodies []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/itemEstimatedMarketValues" {
			b, _ := io.ReadAll(r.Body)
			mu.Lock()
			bodies = append(bodies, string(b))
			mu.Unlock()
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cfg := config.Default()
	cfg.IngestBaseURL = srv.URL
	up := upload.New(cfg, nil)
	up.Start()
	defer up.Stop()

	st := state.New()
	st.SetServer(&config.Server{ID: 1, Name: "Americas"})
	m := NewMarket(st, up, fakeNamer{}, nil)

	// NewEquipmentItem: id=param[1], emv=param[4]/10000, quality=param[6]
	m.handleItemEMV(map[byte]any{1: int16(3409), 4: int64(1353134052), 6: byte(3)})
	// NewSimpleItem: no quality -> defaults to 1
	m.handleItemEMV(map[byte]any{1: int16(914), 4: int64(310674)})
	m.flushEMV()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		joined := strings.Join(bodies, "")
		mu.Unlock()
		if strings.Contains(joined, "ITEM_3409") && strings.Contains(joined, "ITEM_914") {
			if !strings.Contains(joined, `"emv":135313`) || !strings.Contains(joined, `"emv":31`) {
				t.Fatalf("wrong emv values: %s", joined)
			}
			if !strings.Contains(joined, `"quality":3`) {
				t.Fatalf("equipment quality not carried: %s", joined)
			}
			// second flush must NOT re-upload (deduped)
			before := len(bodies)
			m.flushEMV()
			time.Sleep(200 * time.Millisecond)
			mu.Lock()
			after := len(bodies)
			mu.Unlock()
			if after != before {
				t.Fatalf("re-uploaded deduped EMV: %d -> %d", before, after)
			}
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("EMV not uploaded: %v", bodies)
}
