package handlers

import (
	"testing"

	"github.com/albionmarketgg/dataclient/internal/config"
	"github.com/albionmarketgg/dataclient/internal/state"
	"github.com/albionmarketgg/dataclient/internal/upload"
)

// Low-volume items get their history's sale-count array (and sometimes silver)
// wire-encoded as a Photon ByteArray -> []byte. That previously coerced to nil,
// zeroed the row count and silently dropped the whole upload — "some items never
// show history". The upload must go through for both array shapes.
func TestHistoryUploadByteArrayAmounts(t *testing.T) {
	cfg := config.Default()
	up := upload.New(cfg, nil)
	st := state.New()
	st.SetServer(&config.Server{ID: 3, Name: "Europe"})
	st.SetLocation(state.Location{ID: "3005", Name: "Caerleon"})
	st.MarkPacket()
	m := NewMarket(st, up, fakeNamer{}, nil)

	// request arms the ring (item/quality/timescale + location), keyed by msgID
	m.handleHistoryRequest(map[byte]any{
		1: int64(944), 2: int64(1), 3: int64(0), 255: int64(7001),
	})
	// low-volume response: amounts arrive as []byte (incl. zero days),
	// silver/timestamps as regular []any arrays
	m.handleHistoryResponse(0, "", map[byte]any{
		0:   []byte{0, 3, 0, 12},
		1:   []any{int64(0), int64(36000), int64(0), int64(150000)},
		2:   []any{int64(133700000), int64(133700001), int64(133700002), int64(133700003)},
		255: int64(7001),
	})
	if q := up.Stats().Queued; q != 1 {
		t.Fatalf("byte-array history dropped: queued=%d (want 1)", q)
	}

	// regular []any amounts still work
	m.handleHistoryRequest(map[byte]any{
		1: int64(945), 2: int64(1), 3: int64(0), 255: int64(7002),
	})
	m.handleHistoryResponse(0, "", map[byte]any{
		0:   []any{int64(500), int64(410)},
		1:   []any{int64(1), int64(2)},
		2:   []any{int64(133700000), int64(133700001)},
		255: int64(7002),
	})
	if q := up.Stats().Queued; q != 2 {
		t.Fatalf("[]any history broken: queued=%d (want 2)", q)
	}
}
