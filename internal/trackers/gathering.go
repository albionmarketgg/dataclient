package trackers

import (
	"sync"
	"time"

	"github.com/albionmarketgg/dataclient/internal/dispatch"
	"github.com/albionmarketgg/dataclient/internal/photon"
)

// LocalPlayer exposes the local player's object id (satisfied by state.State).
type LocalPlayer interface {
	UserObjectID() (int64, bool)
}

// GatherItem is one aggregated item in a gathering session.
type GatherItem struct {
	ItemID int    `json:"itemId"`
	Name   string `json:"name"`
	Amount int64  `json:"amount"`
	Value  int64  `json:"value"`
	Source string `json:"source"` // "gathering" | "fishing"
}

// GatherSnapshot is the gathering session state.
type GatherSnapshot struct {
	Active         bool         `json:"active"`
	ElapsedSeconds int64        `json:"elapsedSeconds"`
	TotalAmount    int64        `json:"totalAmount"`
	TotalValue     int64        `json:"totalValue"`
	SilverPerHour  int64        `json:"silverPerHour"`
	Items          []GatherItem `json:"items"`
}

// Gathering tracks harvested and fished items in a session.
type Gathering struct {
	info  ItemInfo
	local LocalPlayer

	mu        sync.Mutex
	started   time.Time
	lastAt    time.Time
	active    bool
	items     map[int]*GatherItem
	onChange  func()
	onFeed    func(kind, detail string, count int)
}

// OnFeed registers a live-feed callback for gathered items.
func (g *Gathering) OnFeed(fn func(kind, detail string, count int)) { g.onFeed = fn }

const gatherIdle = 30 * time.Minute

// NewGathering creates a gathering tracker.
func NewGathering(info ItemInfo, local LocalPlayer) *Gathering {
	if info == nil {
		info = nopItems{}
	}
	return &Gathering{info: info, local: local, items: map[int]*GatherItem{}}
}

// OnChange registers a change callback.
func (g *Gathering) OnChange(fn func()) { g.onChange = fn }

// Register attaches gathering handlers.
func (g *Gathering) Register(d *dispatch.Dispatcher) {
	d.OnEvent(photon.EvHarvestFinished, g.onHarvest)
	d.OnEvent(photon.EvRewardGranted, g.onReward)
}

func (g *Gathering) isLocal(objectID int64) bool {
	if g.local == nil {
		return true
	}
	id, ok := g.local.UserObjectID()
	return !ok || id == objectID
}

func (g *Gathering) onHarvest(m map[byte]any) {
	if !g.isLocal(i64(m[0])) {
		return
	}
	itemID := iN(m[4])
	amount := int64(iN(m[5]) + iN(m[6]) + iN(m[7]))
	if itemID <= 0 || amount <= 0 {
		return
	}
	g.record(itemID, 1, amount, "gathering")
}

func (g *Gathering) onReward(m map[byte]any) {
	itemID := iN(m[1])
	amount := int64(iN(m[3]))
	if itemID <= 0 || amount <= 0 {
		return
	}
	g.record(itemID, 1, amount, "fishing")
}

func (g *Gathering) record(itemID, quality int, amount int64, source string) {
	g.mu.Lock()
	now := time.Now()
	if !g.active || now.Sub(g.lastAt) > gatherIdle {
		g.started = now
		g.items = map[int]*GatherItem{}
		g.active = true
	}
	g.lastAt = now
	value := itemValue(g.info, itemID, quality) * amount
	it := g.items[itemID]
	if it == nil {
		it = &GatherItem{ItemID: itemID, Name: itemName(g.info, itemID), Source: source}
		g.items[itemID] = it
	}
	it.Amount += amount
	it.Value += value
	name := it.Name
	g.mu.Unlock()
	if g.onFeed != nil {
		g.onFeed("gather", name, int(amount))
	}
	if g.onChange != nil {
		g.onChange()
	}
}

// Reset clears the current session.
func (g *Gathering) Reset() {
	g.mu.Lock()
	g.active = false
	g.items = map[int]*GatherItem{}
	g.mu.Unlock()
	if g.onChange != nil {
		g.onChange()
	}
}

// Snapshot returns the current gathering session.
func (g *Gathering) Snapshot() GatherSnapshot {
	g.mu.Lock()
	defer g.mu.Unlock()
	snap := GatherSnapshot{Active: g.active}
	if !g.active {
		return snap
	}
	elapsed := time.Since(g.started)
	snap.ElapsedSeconds = int64(elapsed.Seconds())
	for _, it := range g.items {
		snap.Items = append(snap.Items, *it)
		snap.TotalAmount += it.Amount
		snap.TotalValue += it.Value
	}
	if hours := elapsed.Hours(); hours > 0 {
		snap.SilverPerHour = int64(float64(snap.TotalValue) / hours)
	}
	sortItemsDesc(snap.Items)
	return snap
}

func sortItemsDesc(items []GatherItem) {
	for i := 1; i < len(items); i++ {
		for j := i; j > 0 && items[j].Value > items[j-1].Value; j-- {
			items[j], items[j-1] = items[j-1], items[j]
		}
	}
}

var _ = dispatch.Int
