package handlers

import (
	"strconv"
	"testing"
)

type fakeRes struct{}

func (fakeRes) IDByIndex(i int) (string, bool) { return "NODE_" + strconv.Itoa(i), true }
func (fakeRes) Version() string                { return "v1" }

func TestSpecsDecodeFullAchievement(t *testing.T) {
	h := NewSpecs(fakeRes{})
	var gotVer string
	var got []SpecEntry
	h.OnFull(func(version string, e []SpecEntry) { gotVer = version; got = e })

	// param2 = indices, param3 = levels (parallel); param1 = maxed indices -> 100.
	h.handleFull(map[byte]any{
		1: []any{int16(50)},                    // index 50 is maxed
		2: []any{int16(6), int16(7), int16(8)}, // indices
		3: []byte{30, 0, 100},                  // levels: 6->30, 7->0(drop), 8->100
	})

	if gotVer != "v1" {
		t.Fatalf("version: %q", gotVer)
	}
	// expect NODE_6=30, NODE_8=100, NODE_50=100 ; NODE_7 dropped (level 0)
	m := map[string]int{}
	for _, e := range got {
		m[e.ID] = e.Level
	}
	if len(m) != 3 || m["NODE_6"] != 30 || m["NODE_8"] != 100 || m["NODE_50"] != 100 {
		t.Fatalf("bad decode: %+v", got)
	}
	if _, ok := m["NODE_7"]; ok {
		t.Fatal("level-0 node should be dropped")
	}
}

func TestSpecsNoOutputWhenEmpty(t *testing.T) {
	h := NewSpecs(fakeRes{})
	fired := false
	h.OnFull(func(string, []SpecEntry) { fired = true })
	h.handleFull(map[byte]any{2: []any{int16(6)}, 3: []byte{0}}) // only a level-0 node
	if fired {
		t.Fatal("should not fire with no positive-level specs")
	}
}
