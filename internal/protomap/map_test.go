package protomap

import (
	"testing"

	"github.com/albionmarketgg/dataclient/internal/photon"
	"github.com/albionmarketgg/dataclient/internal/phototest"
)

// A 1:1 mirror of the compiled defaults must be a pure no-op.
func TestIdentityPayload(t *testing.T) {
	r := parseRemap([]byte(`{
		"version": 1,
		"events": {
			"UpdateFame": {"code": 82, "params": {"gainedFame": 2, "isPremium": 5, "satchelFame": 10}},
			"NewEquipmentItem": {"code": 30, "params": {"quality": 7, "isAwakened": 11}}
		},
		"operations": {
			"Join": {"code": 2, "params": {"characterName": 2, "locationId": 8}}
		}
	}`))
	if r == nil {
		t.Fatal("mirror payload rejected")
	}
	if !r.identity() {
		t.Fatalf("mirror payload produced translations: %+v", r)
	}
}

// A moved param must translate the wire position back to the compiled one.
func TestParamShift(t *testing.T) {
	r := parseRemap([]byte(`{
		"version": 2,
		"events": {"NewEquipmentItem": {"params": {"quality": 9}}}
	}`))
	if r == nil {
		t.Fatal("payload rejected")
	}
	code, params := r.Event(photon.EvNewEquipmentItem, map[byte]any{0: int64(1), 9: byte(4)})
	if code != photon.EvNewEquipmentItem {
		t.Fatalf("code changed unexpectedly: %v", code)
	}
	if q, _ := params[7].(byte); q != 4 {
		t.Fatalf("quality not moved 9->7: %+v", params)
	}
	if params[0] != int64(1) {
		t.Fatalf("untouched param lost: %+v", params)
	}
}

// A renumbered event code must route the new wire code to the compiled one.
func TestCodeShift(t *testing.T) {
	r := parseRemap([]byte(`{
		"version": 2,
		"events": {"OtherGrabbedLoot": {"code": 281}}
	}`))
	if r == nil {
		t.Fatal("payload rejected")
	}
	code, _ := r.Event(photon.EventCode(281), map[byte]any{})
	if code != photon.EvOtherGrabbedLoot {
		t.Fatalf("wire 281 not remapped: got %v", code)
	}
	// the old (compiled) code must not be double-translated
	code, _ = r.Event(photon.EvUpdateFame, map[byte]any{})
	if code != photon.EvUpdateFame {
		t.Fatalf("unrelated code touched: %v", code)
	}
}

// Operations translate per direction; shared "params" + per-direction override.
func TestOperationDirections(t *testing.T) {
	r := parseRemap([]byte(`{
		"version": 2,
		"operations": {
			"AuctionGetItemAverageStats": {"responseParams": {"silver": 5}}
		}
	}`))
	if r == nil {
		t.Fatal("payload rejected")
	}
	// response: silver moved 5 -> compiled 1
	_, params := r.Response(photon.OpAuctionGetItemAverageStats, map[byte]any{5: int64(777)})
	if params[1] != int64(777) {
		t.Fatalf("response silver not moved: %+v", params)
	}
	// request untouched (its [1] means itemId; no request move was declared)
	_, params = r.Request(photon.OpAuctionGetItemAverageStats, map[byte]any{1: int64(42)})
	if params[1] != int64(42) {
		t.Fatalf("request wrongly translated: %+v", params)
	}
}

// Rejection & tolerance rules.
func TestValidation(t *testing.T) {
	if parseRemap([]byte(`{"version": 0}`)) != nil {
		t.Fatal("stale version accepted")
	}
	if parseRemap([]byte(`not json`)) != nil {
		t.Fatal("junk accepted")
	}
	if parseRemap([]byte(`{"version":1,"events":{"UpdateFame":{"params":{"gainedFame":999}}}}`)) != nil {
		t.Fatal("out-of-range position accepted")
	}
	// unknown packet names / params and v2 extractor objects are ignored, not fatal
	r := parseRemap([]byte(`{
		"version": 1,
		"events": {
			"SomeFuturePacket": {"code": 900, "params": {"x": 1}},
			"UpdateFame": {"params": {"someNewField": 4, "isPremium": {"pos":10,"type":"bitflag"}}}
		}
	}`))
	if r == nil || !r.identity() {
		t.Fatalf("tolerant cases mishandled: %+v", r)
	}
}

// End to end: a payload-shifted packet flows through the real photon parser and
// arrives at the handler translated back to the compiled layout.
type captureHandler struct {
	code   photon.EventCode
	params map[byte]any
}

func (c *captureHandler) HandleEvent(code photon.EventCode, p map[byte]any) { c.code, c.params = code, p }
func (c *captureHandler) HandleRequest(photon.OperationCode, map[byte]any)  {}
func (c *captureHandler) HandleResponse(photon.OperationCode, int16, string, map[byte]any) {
}

func TestParserIntegration(t *testing.T) {
	// simulate a future patch: UpdateFame renumbered to 84, gained fame moved to [20]
	r := parseRemap([]byte(`{
		"version": 2,
		"events": {"UpdateFame": {"code": 84, "params": {"gainedFame": 20}}}
	}`))
	if r == nil {
		t.Fatal("payload rejected")
	}
	h := &captureHandler{}
	parser := photon.NewParser(h)
	parser.SetRemapper(r)
	parser.ReceivePacket(phototest.EventPacket(1, map[byte]any{
		252: int16(84), 20: int64(3650000), 10: int64(0),
	}))
	if h.code != photon.EvUpdateFame {
		t.Fatalf("wire 84 not routed to UpdateFame: %v", h.code)
	}
	if h.params[2] != int64(3650000) {
		t.Fatalf("gainedFame not translated to [2]: %+v", h.params)
	}
}
