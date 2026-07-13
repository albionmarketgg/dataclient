package trackers

import (
	"sync"
	"time"

	"github.com/albionmarketgg/data-client/internal/dispatch"
	"github.com/albionmarketgg/data-client/internal/photon"
)

// LootRecord is a single loot pickup.
type LootRecord struct {
	Time    time.Time `json:"time"`
	Player  string    `json:"player"`
	Item    string    `json:"item"`              // resolved display name (fallback)
	ItemID  string    `json:"itemId,omitempty"`  // canonical unique name incl. @enchant
	Quality int       `json:"quality,omitempty"` // 1-5 (omitted when unknown)
	Amount  int64     `json:"amount"`
	Value   int64     `json:"value"` // client EMV estimate × amount, display silver
	Source  string    `json:"source"`
}

// LootSnapshot lists recent loot records (most recent first).
type LootSnapshot struct {
	Records []LootRecord `json:"records"`
}

// lootItemInfo is a discovered item's type/quality/amount (from New*Item events),
// keyed by its object id for self-loot correlation.
type lootItemInfo struct {
	itemIdx int
	quality int
	amount  int64
}

// selfLootWindow is how long after your InventoryMoveItem a container item leaving
// still counts as your pickup (the move + delete/put land within ~200ms).
const selfLootWindow = 3 * time.Second

// Loot tracks who looted what. "Other" pickups come from OtherGrabbedLoot; your own
// pickups have no single packet, so they're derived by correlating:
//   - New*Item        → item objectId ⇒ type/quality/amount
//   - AttachItemContainer (public, empty privateContainerId) → the item objectIds
//     that sit in a loot container (chest / mob corpse)
//   - InventoryMoveItem (a request YOU send) → arms "I'm taking loot" so another
//     player looting a shared container never counts as yours
//   - InventoryDeleteItem / InventoryPutItem → the item leaving that container ⇒
//     which item (by its original objectId), recorded while armed.
type Loot struct {
	info       ItemInfo
	playerName func() string // local character name (may be nil)
	partyOther func() string // the single OTHER party member's name, "" if 0 or 2+ (may be nil)

	mu      sync.Mutex
	records []LootRecord
	// self-loot correlation state:
	discovered map[int64]lootItemInfo // item objectId -> type/quality/amount
	bagSource  map[int64]string       // loot-container objectId -> source name
	bagItems   map[int64]int64        // item objectId -> its loot-container objectId
	bagGuid    map[string]int64       // container guid -> loot-container objectId
	armedAt    time.Time              // time of your last move OUT of a loot container
	lastPlayer string                 // cached local character name (survives zone gaps)
	onChange   func()
	onFeed     func(kind, detail string, count int)
}

// resolveSelf returns the local character's name, caching the last non-empty value
// so a momentary gap (zone transition) doesn't flip a pickup to the "You" fallback.
func (l *Loot) resolveSelf() string {
	name := ""
	if l.playerName != nil {
		name = l.playerName()
	}
	l.mu.Lock()
	if name != "" {
		l.lastPlayer = name
	} else {
		name = l.lastPlayer
	}
	l.mu.Unlock()
	if name == "" {
		name = "You"
	}
	return name
}

// OnFeed registers a live-feed callback for loot pickups.
func (l *Loot) OnFeed(fn func(kind, detail string, count int)) { l.onFeed = fn }

// NewLoot creates a loot tracker. playerName supplies the local character (for
// attributing your own pickups); partyOther supplies the single other party
// member's name (for attributing their chest pickups). Both may be nil.
func NewLoot(info ItemInfo, playerName, partyOther func() string) *Loot {
	if info == nil {
		info = nopItems{}
	}
	return &Loot{
		info:       info,
		playerName: playerName,
		partyOther: partyOther,
		discovered: map[int64]lootItemInfo{},
		bagSource:  map[int64]string{},
		bagItems:   map[int64]int64{},
		bagGuid:    map[string]int64{},
	}
}

// OnChange registers a change callback.
func (l *Loot) OnChange(fn func()) { l.onChange = fn }

// Register attaches loot handlers.
func (l *Loot) Register(d *dispatch.Dispatcher) {
	d.OnEvent(photon.EvOtherGrabbedLoot, l.onGrabbed)
	// party members' chest pickups arrive as "types removed from the loot pool"
	// (item type + count, no name) — attributed to the other member.
	d.OnEvent(photon.EvPartyLootItemTypesRemoved, l.onPartyRemoved)
	// self-loot correlation:
	for _, code := range []photon.EventCode{
		photon.EvNewSimpleItem, photon.EvNewEquipmentItem, photon.EvNewFurnitureItem,
		photon.EvNewJournalItem, photon.EvNewLaborerItem, photon.EvNewKillTrophyItem,
		photon.EvNewSiegeBannerItem,
	} {
		d.OnEvent(code, l.onNewItem)
	}
	d.OnEvent(photon.EvNewLoot, l.onLootSource)          // mob-corpse source
	d.OnEvent(photon.EvNewLootChest, l.onNewLootChest)   // chest source (with name)
	d.OnEvent(photon.EvLootChestOpened, l.onLootChest)   // chest source (fallback)
	d.OnEvent(photon.EvAttachItemContainer, l.onAttachContainer)
	d.OnEvent(photon.EvDetachItemContainer, l.onDetachContainer)
	// your move out of a loot container arms a pickup; the item then leaves via
	// delete (it stacked → new objectId) or put (it kept its objectId).
	d.OnRequest(photon.OpInventoryMoveItem, l.onMoveItem)
	d.OnEvent(photon.EvInventoryDeleteItem, l.onItemTaken)
	d.OnEvent(photon.EvInventoryPutItem, l.onItemTaken)
}

// onNewItem records an item's type/quality/amount by object id (New*Item events:
// [0]=objectId [1]=itemId [2]=amount [6]=quality).
func (l *Loot) onNewItem(m map[byte]any) {
	obj := i64(m[0])
	idx := iN(m[1])
	if obj == 0 || idx <= 0 {
		return
	}
	amount := i64(m[2])
	if amount <= 0 {
		amount = 1
	}
	quality := 1
	if q := iN(m[6]); q > 0 {
		quality = q
	}
	l.mu.Lock()
	if len(l.discovered) > 8192 {
		l.discovered = map[int64]lootItemInfo{}
	}
	l.discovered[obj] = lootItemInfo{itemIdx: idx, quality: quality, amount: amount}
	l.mu.Unlock()
}

// onLootSource labels a mob-corpse loot container (NewLoot: [0]=objectId [3]=source).
func (l *Loot) onLootSource(m map[byte]any) {
	obj := i64(m[0])
	if obj == 0 {
		return
	}
	src := str(m[3])
	if src == "" {
		return
	}
	l.mu.Lock()
	if len(l.bagSource) > 4096 {
		l.bagSource = map[int64]string{}
	}
	l.bagSource[obj] = src
	l.mu.Unlock()
}

// onNewLootChest labels a loot-chest container when it appears (NewLootChest:
// [0]=objectId [3]=chest unique name). This usually precedes the attach.
func (l *Loot) onNewLootChest(m map[byte]any) {
	obj := i64(m[0])
	if obj == 0 {
		return
	}
	name := str(m[3])
	if name == "" {
		name = "Loot chest"
	}
	l.mu.Lock()
	if len(l.bagSource) > 4096 {
		l.bagSource = map[int64]string{}
	}
	l.bagSource[obj] = name
	l.mu.Unlock()
}

// onLootChest labels a loot-chest container (LootChestOpened: [0]=objectId).
func (l *Loot) onLootChest(m map[byte]any) {
	obj := i64(m[0])
	if obj == 0 {
		return
	}
	l.mu.Lock()
	if _, ok := l.bagSource[obj]; !ok {
		l.bagSource[obj] = "Loot chest"
	}
	l.mu.Unlock()
}

// onAttachContainer records the items inside a loot container you open (Attach:
// [0]=objectId [1]=containerId [3]=slot item objectIds). Only containers whose
// objectId was announced as a loot source (NewLoot / NewLootChest / LootChestOpened)
// are tracked — your own inventory / bank are never announced, so rearranging your
// items never counts as loot. (Do NOT use privateContainerId: real loot chests set
// it too.)
func (l *Loot) onAttachContainer(m map[byte]any) {
	obj := i64(m[0])
	if obj == 0 {
		return
	}
	guid := guidKey(m[1])
	slots := dispatch.Int64Slice(m[3])
	l.mu.Lock()
	if _, isLoot := l.bagSource[obj]; !isLoot {
		l.mu.Unlock()
		return // not an identified loot container
	}
	if guid != "" {
		l.bagGuid[guid] = obj
	}
	for _, it := range slots {
		if it != 0 {
			l.bagItems[it] = obj
		}
	}
	l.mu.Unlock()
}

// onDetachContainer clears a container's item links when you close it (Detach:
// [0]=containerId), so items that despawn afterwards aren't miscounted.
func (l *Loot) onDetachContainer(m map[byte]any) {
	guid := guidKey(m[0])
	if guid == "" {
		return
	}
	l.mu.Lock()
	if obj, ok := l.bagGuid[guid]; ok {
		for it, b := range l.bagItems {
			if b == obj {
				delete(l.bagItems, it)
			}
		}
		delete(l.bagGuid, guid)
	}
	l.mu.Unlock()
}

// onMoveItem arms a self-pickup when YOU move an item out of a loot container into
// a non-loot container (InventoryMoveItem request: [1]=source container guid,
// [4]=destination container guid). It's a request, so it only fires for your own
// actions — another player looting a shared container never arms it.
func (l *Loot) onMoveItem(m map[byte]any) {
	src, dst := guidKey(m[1]), guidKey(m[4])
	if src == "" {
		return
	}
	l.mu.Lock()
	_, fromLoot := l.bagGuid[src]
	_, toLoot := l.bagGuid[dst]
	if fromLoot && !toLoot {
		l.armedAt = time.Now()
	}
	l.mu.Unlock()
}

// onItemTaken fires when a container item leaves — via InventoryDeleteItem (it
// stacked, so a new objectId was made) or InventoryPutItem (it kept its objectId).
// If the item was in a loot container AND you armed a pickup recently, it's yours.
func (l *Loot) onItemTaken(m map[byte]any) {
	itemObj := i64(m[0])
	if itemObj == 0 {
		return
	}
	l.mu.Lock()
	bag, ok := l.bagItems[itemObj]
	if !ok || time.Now().Sub(l.armedAt) > selfLootWindow {
		l.mu.Unlock()
		return
	}
	info, disc := l.discovered[itemObj]
	src := l.bagSource[bag]
	delete(l.bagItems, itemObj) // dedup: only record once
	delete(l.discovered, itemObj)
	l.mu.Unlock()
	if !disc {
		return // unknown item type — don't guess
	}
	l.add(LootRecord{
		Time:    time.Now().UTC(),
		Player:  l.resolveSelf(),
		Item:    itemName(l.info, info.itemIdx),
		ItemID:  itemUnique(l.info, info.itemIdx),
		Quality: info.quality,
		Amount:  info.amount,
		Value:   itemValue(l.info, info.itemIdx, info.quality) * info.amount,
		Source:  src,
	})
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
		Time:    time.Now().UTC(),
		Player:  player,
		Item:    itemName(l.info, itemID),
		ItemID:  itemUnique(l.info, itemID),
		Quality: 1, // OtherGrabbedLoot carries no quality
		Amount:  amount,
		Value:   itemValue(l.info, itemID, 1) * amount,
		Source:  str(m[1]),
	})
}

// onPartyRemoved records another party member taking items from a shared loot
// chest (PartyLootItemTypesRemoved: [1]=item types [2]/[4]=counts [3]=isSilver
// flags). It fires only for OTHER members (your own pickups go through the move
// path), so it never double-counts you. Attribution: the single other party
// member's name (2-person party), else "Party".
func (l *Loot) onPartyRemoved(m map[byte]any) {
	types := dispatch.Int64Slice(m[1])
	if len(types) == 0 {
		return
	}
	silver := boolSlice(m[3])
	player := "Party"
	if l.partyOther != nil {
		if n := l.partyOther(); n != "" {
			player = n
		}
	}
	for i, t := range types {
		if t <= 0 || (i < len(silver) && silver[i]) {
			continue
		}
		id := int(t)
		l.add(LootRecord{
			Time:    time.Now().UTC(),
			Player:  player,
			Item:    itemName(l.info, id),
			ItemID:  itemUnique(l.info, id),
			Quality: 1,
			Amount:  1,
			Value:   itemValue(l.info, id, 1),
			Source:  "Loot chest (party)",
		})
	}
}

// boolSlice extracts a []bool from a bool-array param.
func boolSlice(v any) []bool {
	switch a := v.(type) {
	case []bool:
		return a
	case []any:
		out := make([]bool, len(a))
		for i, e := range a {
			if b, ok := e.(bool); ok {
				out[i] = b
			}
		}
		return out
	}
	return nil
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
