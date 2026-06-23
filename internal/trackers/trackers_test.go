package trackers

import (
	"sync"
	"testing"
	"time"

	"github.com/niick1231/albionmarket_dataclient/internal/dispatch"
	"github.com/niick1231/albionmarket_dataclient/internal/photon"
	"github.com/niick1231/albionmarket_dataclient/internal/phototest"
)

type fakeItems struct{}

func (fakeItems) UniqueName(id int) (string, bool)  { return "ITEM_" + itoa(id), true }
func (fakeItems) DisplayName(id int) (string, bool) { return "Item " + itoa(id), true }
func (fakeItems) EMV(id, q int) (int64, bool)       { return int64(id * 100), true }

type fakeLocal struct{ id int64 }

func (f fakeLocal) UserObjectID() (int64, bool) { return f.id, true }

func newPipe(reg func(*dispatch.Dispatcher)) *photon.Parser {
	d := dispatch.New()
	reg(d)
	return photon.NewParser(d)
}

func ev(code photon.EventCode, params map[byte]any) []byte {
	params[252] = int16(code)
	return phototest.EventPacket(1, params)
}

func TestPartyTracker(t *testing.T) {
	p := NewParty()
	parser := newPipe(p.Register)
	parser.ReceivePacket(phototest.ResponsePacket(byte(photon.OpJoin), 0, map[byte]any{
		253: int16(photon.OpJoin), 0: int64(50), 1: "me-guid", 2: "Me", 8: "3005",
	}))
	parser.ReceivePacket(ev(photon.EvPartyJoined, map[byte]any{
		5: []string{"a-guid", "b-guid"}, 6: []string{"Alice", "Bob"},
	}))
	snap := p.Snapshot()
	if len(snap.Members) != 3 {
		t.Fatalf("got %d members: %+v", len(snap.Members), snap.Members)
	}
	if !snap.Members[0].IsLocal || snap.Members[0].Name != "Me" {
		t.Fatalf("local player not first: %+v", snap.Members[0])
	}
}

func TestGatheringTracker(t *testing.T) {
	g := NewGathering(fakeItems{}, fakeLocal{id: 50})
	parser := newPipe(g.Register)
	parser.ReceivePacket(ev(photon.EvHarvestFinished, map[byte]any{
		0: int64(50), 4: int32(10), 5: int32(3), 6: int32(2), 7: int32(0),
	}))
	snap := g.Snapshot()
	if !snap.Active || snap.TotalAmount != 5 {
		t.Fatalf("gathering: active=%v total=%d", snap.Active, snap.TotalAmount)
	}
	if snap.TotalValue != 5*int64(10*100) {
		t.Fatalf("gathering value: %d", snap.TotalValue)
	}
	// harvest from a different player should be ignored
	parser.ReceivePacket(ev(photon.EvHarvestFinished, map[byte]any{
		0: int64(999), 4: int32(11), 5: int32(50), 6: int32(0), 7: int32(0),
	}))
	if g.Snapshot().TotalAmount != 5 {
		t.Fatalf("non-local harvest was recorded")
	}
}

func TestCombatTracker(t *testing.T) {
	c := NewCombat(nil)
	parser := newPipe(c.Register)
	parser.ReceivePacket(ev(photon.EvNewCharacter, map[byte]any{0: int64(50), 1: "Me", 7: "me-guid"}))
	parser.ReceivePacket(ev(photon.EvNewMob, map[byte]any{0: int64(900), 1: int32(7)}))
	for i := 0; i < 3; i++ {
		parser.ReceivePacket(ev(photon.EvHealthUpdate, map[byte]any{
			0: int64(900), 2: int64(-100), 6: int64(50),
		}))
	}
	snap := c.Snapshot()
	if len(snap.Encounters) != 1 {
		t.Fatalf("encounters: %d", len(snap.Encounters))
	}
	ps := snap.Encounters[0].Participants
	if len(ps) == 0 || ps[0].Name != "Me" || ps[0].DamageDealt != 300 {
		t.Fatalf("combat attribution wrong: %+v", ps)
	}
}

func TestDungeonEventPairing(t *testing.T) {
	old := dungeonPairWindow
	dungeonPairWindow = 60 * time.Millisecond
	defer func() { dungeonPairWindow = old }()

	c := NewCombat(nil)
	parser := newPipe(c.Register)

	var mu sync.Mutex
	var events []DungeonEvent
	c.OnDungeonEvent(func(e DungeonEvent) { mu.Lock(); events = append(events, e); mu.Unlock() })

	// fame, then silver shortly after -> they pair into one event
	parser.ReceivePacket(ev(photon.EvUpdateFame, map[byte]any{2: int64(3650000), 10: int64(0)}))
	time.Sleep(15 * time.Millisecond)
	parser.ReceivePacket(ev(photon.EvTakeSilver, map[byte]any{3: int64(12000000)}))
	// silver with no fame nearby -> dropped
	time.Sleep(200 * time.Millisecond)
	parser.ReceivePacket(ev(photon.EvTakeSilver, map[byte]any{3: int64(99990000)}))
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(events) != 1 {
		t.Fatalf("want 1 dungeon event, got %d: %+v", len(events), events)
	}
	if events[0].Fame != 365 || events[0].Silver != 1200 {
		t.Fatalf("bad pairing: %+v", events[0])
	}
}

func TestLootTracker(t *testing.T) {
	l := NewLoot(fakeItems{})
	parser := newPipe(l.Register)
	parser.ReceivePacket(ev(photon.EvOtherGrabbedLoot, map[byte]any{
		1: "Mob", 2: "Alice", 3: false, 4: int32(20), 5: int64(4),
	}))
	// silver pickups ignored
	parser.ReceivePacket(ev(photon.EvOtherGrabbedLoot, map[byte]any{
		1: "Mob", 2: "Alice", 3: true, 4: int32(0), 5: int64(1000),
	}))
	snap := l.Snapshot()
	if len(snap.Records) != 1 {
		t.Fatalf("loot records: %d", len(snap.Records))
	}
	if snap.Records[0].Player != "Alice" || snap.Records[0].Amount != 4 || snap.Records[0].Value != 4*2000 {
		t.Fatalf("loot record wrong: %+v", snap.Records[0])
	}
}
