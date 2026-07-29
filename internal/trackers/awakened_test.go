package trackers

import "testing"

// Real packets captured 2026-07 after the game patch that inserted a value field
// at [5] (shifting crafter→[6], quality→[7]) and moved the awakened flag from an
// int at [10] (now an opaque bytes blob on every item) to a byte at [11].
var (
	// NewEquipmentItem (30) — the awakened mace, quality 4, item index 9086.
	maceEquip = map[byte]any{
		0: int64(224), 1: int16(9086), 2: int16(1),
		4: int64(7620015055), 5: int64(8162751111),
		6: "lissychann", 7: byte(4), 8: int64(611006000),
		10: []byte{0x4a}, 11: byte(1),
	}
	// NewEquipmentItem (30) — a non-awakened item ([11]=0), quality 3.
	plainEquip = map[byte]any{
		0: int64(215), 1: int16(7188), 2: int16(1),
		4: int64(17122232932), 5: int64(9237000000),
		6: "Shrinks", 7: byte(3), 8: int64(3251390000),
		10: []byte{0x3d}, 11: byte(0),
	}
	maceSoulID = []byte{0xe2, 0x70, 0x16, 0x29, 0x84, 0xe4, 0xe1, 0x4c, 0xb8, 0x49, 0x69, 0x20, 0x1f, 0xdb, 0xf6, 0xda}
)

// maceSoul returns the mace's NewEquipmentItemLegendarySoul (37) packet.
func maceSoul() map[byte]any {
	return map[byte]any{
		0: int64(224), 1: append([]byte(nil), maceSoulID...), 3: int32(1), 4: true,
		6: int64(23316), 7: int64(48497280),
		8: []string{"TRAIT_ITEM_POWER"}, 9: []any{0.014731574811793848},
		12: int64(2080000),
	}
}

func oneAwakened(t *testing.T, a *Awakened) AwakenedItem {
	t.Helper()
	snap := a.Snapshot()
	if len(snap.Items) != 1 {
		t.Fatalf("want 1 tracked item, got %d: %+v", len(snap.Items), snap.Items)
	}
	return snap.Items[0]
}

// New wire layout: equipment + soul merge into a complete item, and the sync
// body carries itemId + quality (without them the backend cannot compute trait
// values — the "value unavailable" bug).
func TestAwakenedNewLayout(t *testing.T) {
	a := NewAwakened(fakeItems{}, func() int { return 3 }, func() string { return "Hero" }, nil)
	a.onEquipment(maceEquip)
	a.onSoul(maceSoul())

	it := oneAwakened(t, a)
	if it.SoulID != "e270162984e4e14cb84969201fdbf6da" {
		t.Fatalf("soulId: %q", it.SoulID)
	}
	if it.ItemID != "ITEM_9086" || it.Quality != 4 {
		t.Fatalf("static info missing: itemId=%q quality=%d", it.ItemID, it.Quality)
	}
	if it.AttunedTo != "Hero" || len(it.Traits) != 1 || it.Traits[0].ID != "TRAIT_ITEM_POWER" {
		t.Fatalf("soul fields: %+v", it)
	}

	body := a.SyncBody().(map[string]any)
	items := body["items"].([]map[string]any)
	if len(items) != 1 || items[0]["itemId"] != "ITEM_9086" || items[0]["quality"] != 4 {
		t.Fatalf("sync body incomplete: %+v", items)
	}
}

// Non-awakened equipment ([11]=0) must not be tracked (would flood the UI with
// every piece of gear on screen).
func TestAwakenedIgnoresPlainEquipment(t *testing.T) {
	a := NewAwakened(fakeItems{}, nil, nil, nil)
	a.onEquipment(plainEquip)
	if n := len(a.Snapshot().Items); n != 0 {
		t.Fatalf("plain equipment tracked: %d items", n)
	}
}

// Resilience: even if the awakened flag breaks again (equipment looks
// non-awakened), the soul packet must still materialize a complete item from the
// equip cache — the sync body keeps itemId/quality.
func TestAwakenedSoulBackfillsWhenFlagMisses(t *testing.T) {
	a := NewAwakened(fakeItems{}, nil, nil, nil)
	eq := map[byte]any{}
	for k, v := range maceEquip {
		eq[k] = v
	}
	eq[11] = byte(0) // simulate the flag breaking/moving again
	a.onEquipment(eq)
	if len(a.Snapshot().Items) != 0 {
		t.Fatal("unflagged equipment should not be tracked yet")
	}
	a.onSoul(maceSoul())
	it := oneAwakened(t, a)
	if it.ItemID != "ITEM_9086" || it.Quality != 4 {
		t.Fatalf("cache backfill failed: itemId=%q quality=%d", it.ItemID, it.Quality)
	}
}

// Soul arriving before its equipment packet: the later equipment packet must
// backfill the already-materialized item.
func TestAwakenedEquipmentAfterSoul(t *testing.T) {
	a := NewAwakened(fakeItems{}, nil, nil, nil)
	a.onSoul(maceSoul())
	it := oneAwakened(t, a)
	if it.ItemID != "" {
		t.Fatalf("unexpected itemId before equipment packet: %q", it.ItemID)
	}
	a.onEquipment(maceEquip)
	it = oneAwakened(t, a)
	if it.ItemID != "ITEM_9086" || it.Quality != 4 {
		t.Fatalf("late equipment did not backfill: itemId=%q quality=%d", it.ItemID, it.Quality)
	}
}

// The pre-patch layout (int flag at [10], quality at [6]) must keep working.
func TestAwakenedOldLayout(t *testing.T) {
	a := NewAwakened(fakeItems{}, nil, nil, nil)
	a.onEquipment(map[byte]any{0: int64(300), 1: int16(500), 6: int32(2), 10: int64(1)})
	it := oneAwakened(t, a)
	if it.ItemID != "ITEM_500" || it.Quality != 2 {
		t.Fatalf("old layout broken: itemId=%q quality=%d", it.ItemID, it.Quality)
	}
}
