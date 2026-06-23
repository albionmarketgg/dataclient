package specs

import (
	"strings"
	"testing"
)

const sampleXML = `<?xml version="1.0" encoding="utf-8"?>
<achievements Version="abc-123">
  <template name="REFINE_T4"><baselevels>x</baselevels></template>
  <achievement id="ADVENTURER_ADEPT"><requirement/></achievement>
  <templateachievement id="CRAFT_QUARTERSTAFFS">
    <achievementlevels>
      <level value="1"/>
    </achievementlevels>
  </templateachievement>
  <achievement id="CRAFT_QUARTERSTAFFS_DOUBLEBLADED"/>
  <somethingelse id="IGNORED"/>
  <achievement><noid/></achievement>
  <achievement id="COMBAT_AXES"/>
</achievements>`

func TestParseIndexMapping(t *testing.T) {
	s := New("", "")
	if !s.parse(strings.NewReader(sampleXML)) {
		t.Fatal("parse failed")
	}
	if !s.Ready() {
		t.Fatal("not ready")
	}
	if s.Version() != "abc-123" {
		t.Fatalf("version: %q", s.Version())
	}
	// nested <achievement> inside templateachievement and non-achievement elements
	// must NOT shift the index; elements without an id are skipped.
	want := map[int]string{
		0: "ADVENTURER_ADEPT",
		1: "CRAFT_QUARTERSTAFFS",
		2: "CRAFT_QUARTERSTAFFS_DOUBLEBLADED",
		3: "COMBAT_AXES",
	}
	for i, id := range want {
		got, ok := s.IDByIndex(i)
		if !ok || got != id {
			t.Fatalf("index %d: got %q (ok=%v), want %q", i, got, ok, id)
		}
	}
	if _, ok := s.IDByIndex(4); ok {
		t.Fatal("index 4 should be out of range")
	}
}
