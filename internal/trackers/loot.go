package trackers

import (
	"sync"
	"time"

	"github.com/niick1231/albionmarket_dataclient/internal/dispatch"
	"github.com/niick1231/albionmarket_dataclient/internal/photon"
)

// LootRecord is a single loot pickup.
type LootRecord struct {
	Time   time.Time `json:"time"`
	Player string    `json:"player"`
	Item   string    `json:"item"`
	Amount int64     `json:"amount"`
	Value  int64     `json:"value"`
	Source string    `json:"source"`
}

// LootSnapshot lists recent loot records (most recent first).
type LootSnapshot struct {
	Records []LootRecord `json:"records"`
}

// Loot tracks who looted what.
type Loot struct {
	info ItemInfo

	mu       sync.Mutex
	records  []LootRecord
	onChange func()
	onFeed   func(kind, detail string, count int)
}

// OnFeed registers a live-feed callback for loot pickups.
func (l *Loot) OnFeed(fn func(kind, detail string, count int)) { l.onFeed = fn }

// NewLoot creates a loot tracker.
func NewLoot(info ItemInfo) *Loot {
	if info == nil {
		info = nopItems{}
	}
	return &Loot{info: info}
}

// OnChange registers a change callback.
func (l *Loot) OnChange(fn func()) { l.onChange = fn }

// Register attaches loot handlers.
func (l *Loot) Register(d *dispatch.Dispatcher) {
	d.OnEvent(photon.EvOtherGrabbedLoot, l.onGrabbed)
	d.OnEvent(photon.EvPartyLootItems, l.onPartyLoot)
}

func (l *Loot) add(rec LootRecord) {
	l.mu.Lock()
	l.records = append(l.records, rec)
	if len(l.records) > 500 {
		l.records = l.records[len(l.records)-500:]
	}
	l.mu.Unlock()
	if l.onFeed != nil {
		l.onFeed("loot", rec.Player+": "+rec.Item, int(rec.Amount))
	}
	if l.onChange != nil {
		l.onChange()
	}
}

func (l *Loot) onGrabbed(m map[byte]any) {
	if isSilver, _ := m[3].(bool); isSilver {
		return // silver pickups are not items
	}
	itemID := iN(m[4])
	if itemID <= 0 {
		return
	}
	amount := i64(m[5])
	player := str(m[2])
	if player == "" {
		player = "Unknown"
	}
	l.add(LootRecord{
		Time:   time.Now(),
		Player: player,
		Item:   itemName(l.info, itemID),
		Amount: amount,
		Value:  itemValue(l.info, itemID, 1) * amount,
		Source: str(m[1]),
	})
}

func (l *Loot) onPartyLoot(m map[byte]any) {
	itemIDs := dispatch.Int64Slice(m[2])
	qualities := dispatch.Int64Slice(m[4])
	amounts := dispatch.Int64Slice(m[9])
	players := dispatch.Strings(m[10])
	for i := range itemIDs {
		id := int(itemIDs[i])
		if id <= 0 {
			continue
		}
		quality := 1
		if i < len(qualities) && qualities[i] > 0 {
			quality = int(qualities[i])
		}
		var amount int64 = 1
		if i < len(amounts) {
			amount = amounts[i]
		}
		player := "Party"
		if i < len(players) && players[i] != "" {
			player = players[i]
		}
		l.add(LootRecord{
			Time:   time.Now(),
			Player: player,
			Item:   itemName(l.info, id),
			Amount: amount,
			Value:  itemValue(l.info, id, quality) * amount,
			Source: "party chest",
		})
	}
}

// Clear empties the loot log.
func (l *Loot) Clear() {
	l.mu.Lock()
	l.records = nil
	l.mu.Unlock()
	if l.onChange != nil {
		l.onChange()
	}
}

// Snapshot returns recent loot records, most recent first.
func (l *Loot) Snapshot() LootSnapshot {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]LootRecord, 0, len(l.records))
	for i := len(l.records) - 1; i >= 0; i-- {
		out = append(out, l.records[i])
	}
	return LootSnapshot{Records: out}
}
