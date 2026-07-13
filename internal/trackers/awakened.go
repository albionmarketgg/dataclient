package trackers

import (
	"encoding/hex"
	"math"
	"sort"
	"strconv"
	"sync"

	"github.com/albionmarketgg/dataclient/internal/dispatch"
	"github.com/albionmarketgg/dataclient/internal/photon"
)

func round2(f float64) float64 { return math.Round(f*100) / 100 }

// Trait values are NOT computed here. Strict separation of concerns: the client
// captures the raw roll and sends it; the website backend computes the real value
// (base × rollFactor × progression^(itemPower/100)) and returns it in the sync
// response, which we apply via ApplyValues. See AWAKENED_TRAIT_VALUE_MIGRATION.md.

// traitNames maps raw packet trait ids to display labels. Only these 19 exist, so
// a static map is enough (no runtime lookup pipeline). Used for local display;
// uploads send the raw id.
var traitNames = map[string]string{
	"TRAIT_ITEM_POWER":             "Item Power",
	"TRAIT_HITPOINTS_MAX":          "Max Health",
	"TRAIT_ENERGY_MAX":             "Max Energy",
	"TRAIT_ENERGY_COST_REDUCTION":  "Energy Cost Reduction",
	"TRAIT_DEFENSE_BONUS":          "Defense Bonus",
	"TRAIT_ABILITY_DAMAGE":         "Ability Damage",
	"TRAIT_AUTO_ATTACK_DAMAGE":     "Auto Attack Damage",
	"TRAIT_HEALING_DEALT":          "Healing Dealt",
	"TRAIT_HEALING_RECEIVED":       "Healing Received",
	"TRAIT_CC_RESIST":              "CC Resistance",
	"TRAIT_CC_DURATION":            "CC Duration",
	"TRAIT_ATTACK_SPEED":           "Attack Speed",
	"TRAIT_ATTACK_RANGE":           "Attack Range",
	"TRAIT_CAST_SPEED_INCREASE":    "Cast Speed Increase",
	"TRAIT_COOLDOWN_REDUCTION":     "Cooldown Reduction",
	"TRAIT_THREAT_BONUS":           "Threat Bonus",
	"TRAIT_MOB_FAME":               "Mob Fame Modifier",
	"TRAIT_RESILIENCE_PENETRATION": "Resilience Penetration",
	"TRAIT_LIFESTEAL":              "Lifesteal",
}

func traitName(id string) string {
	if n, ok := traitNames[id]; ok {
		return n
	}
	return id
}

// AwakenedTrait is one rolled trait on an awakened item. ID is the raw packet id
// (e.g. "TRAIT_ITEM_POWER"); Name is labeled locally for display. Roll is the raw
// 0-1 packet roll we upload; Value/Percent are filled from the backend sync
// response (see ApplyValue) — the client computes no values itself.
type AwakenedTrait struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Roll     float64 `json:"roll"`     // raw 0-1 roll (uploaded; backend computes the value)
	Value    float64 `json:"value"`    // backend-computed display value (0 until resolved)
	Percent  bool    `json:"percent"`  // backend-provided: render as a percentage
	Resolved bool    `json:"resolved"` // true once the backend returned a value
}

// AwakenedItem is a single awakened/attuned item in the player's inventory,
// assembled from NewEquipmentItem (30) + NewEquipmentItemLegendarySoul (37).
type AwakenedItem struct {
	Key        string          `json:"key"`        // stable id (soulId when known, else objectId)
	SoulID     string          `json:"soulId"`     // stable item-instance id (dedup key)
	ItemID     string          `json:"itemId"`     // unique item name → icon lookup
	ServerID   int             `json:"serverId"`   // region
	Name       string          `json:"name"`       // display name
	Quality    int             `json:"quality"`    // 1-5
	AttunedTo  string          `json:"attunedTo"`  // character it's attuned to
	Attunement float64         `json:"attunement"` // attunement progress (/10000)
	Strain     float64         `json:"strain"`     // accumulated strain (/10000)
	Traits     []AwakenedTrait `json:"traits"`     // up to 3
	Location   string          `json:"location"`   // where the item lives (from container events)

	updated int64 // per-tracker sequence of the last packet update (dedup tie-break)
}

// AwakenedSnapshot is the current awakened-inventory list.
type AwakenedSnapshot struct {
	Items []AwakenedItem `json:"items"`
}

// Awakened tracks awakened items from equipment + legendary-soul packets.
type Awakened struct {
	info       ItemInfo
	serverID   func() int
	playerName func() string
	location   func() string // current mapped location name ("" if unmapped)

	mu       sync.Mutex
	items    map[string]*AwakenedItem // keyed by objectId (string)
	seq      int64                    // monotonic update counter (stamps AwakenedItem.updated)
	onChange func()
}

// NewAwakened creates the tracker. info resolves item names; serverID/playerName/
// location supply the current region, local character, and mapped location (all
// may be nil).
func NewAwakened(info ItemInfo, serverID func() int, playerName func() string, location func() string) *Awakened {
	if info == nil {
		info = nopItems{}
	}
	return &Awakened{info: info, serverID: serverID, playerName: playerName, location: location, items: map[string]*AwakenedItem{}}
}

// curLocation returns the current mapped location name (empty when unmapped).
func (a *Awakened) curLocation() string {
	if a.location == nil {
		return ""
	}
	return a.location()
}

// OnChange registers a change callback.
func (a *Awakened) OnChange(fn func()) { a.onChange = fn }

// Register attaches the equipment + legendary-soul handlers.
func (a *Awakened) Register(d *dispatch.Dispatcher) {
	d.OnEvent(photon.EvNewEquipmentItem, a.onEquipment)
	d.OnEvent(photon.EvNewEquipmentItemLegendarySoul, a.onSoul)
}

func (a *Awakened) at(obj int64) *AwakenedItem {
	key := strconv.FormatInt(obj, 10)
	it := a.items[key]
	if it == nil {
		it = &AwakenedItem{Key: key}
		if a.serverID != nil {
			it.ServerID = a.serverID()
		}
		a.items[key] = it
	}
	return it
}

// onEquipment handles NewEquipmentItem (30): only awakened items (param 10) are
// tracked; supplies itemId, name and quality.
func (a *Awakened) onEquipment(m map[byte]any) {
	if i64(m[10]) == 0 { // not awakened
		return
	}
	obj := i64(m[0])
	if obj == 0 {
		return
	}
	idx := iN(m[1])
	uniq, _ := a.info.UniqueName(idx)
	name := uniq
	if dn, ok := a.info.DisplayName(idx); ok && dn != "" {
		name = dn
	}
	loc := a.curLocation()
	a.mu.Lock()
	it := a.at(obj)
	it.ItemID = uniq
	it.Name = name
	it.Quality = iN(m[6])
	if loc != "" {
		it.Location = loc
	}
	a.seq++
	it.updated = a.seq
	a.mu.Unlock()
	a.fire()
}

// onSoul handles NewEquipmentItemLegendarySoul (37): attuned-to, strain,
// attunement and the parallel trait id/value arrays.
func (a *Awakened) onSoul(m map[byte]any) {
	obj := i64(m[0])
	if obj == 0 {
		return
	}
	attunedTo := str(m[5])
	if b, _ := m[4].(bool); b && a.playerName != nil { // attunedToMe → the local player
		if p := a.playerName(); p != "" {
			attunedTo = p
		}
	}
	ids := dispatch.Strings(m[8])
	vals := toFloatSlice(m[9])
	traits := make([]AwakenedTrait, 0, len(ids))
	for i, id := range ids {
		roll := 0.0
		if i < len(vals) {
			roll = vals[i]
		}
		// Store the raw roll only; the backend computes and returns the value.
		traits = append(traits, AwakenedTrait{ID: id, Name: traitName(id), Roll: roll})
	}
	soulID := ""
	if b, ok := m[1].([]byte); ok && len(b) > 0 {
		soulID = hex.EncodeToString(b)
	}
	loc := a.curLocation()
	a.mu.Lock()
	it := a.at(obj)
	if soulID != "" {
		it.SoulID = soulID
		it.Key = soulID // stable identity → dedup + selection
	}
	if loc != "" {
		it.Location = loc
	}
	it.AttunedTo = attunedTo
	it.Strain = float64(i64(m[6])) / 10000
	it.Attunement = float64(i64(m[7])) / 10000
	it.Traits = traits
	a.seq++
	it.updated = a.seq
	a.mu.Unlock()
	a.fire()
}

func (a *Awakened) fire() {
	if a.onChange != nil {
		a.onChange()
	}
}

// Upsert adds or updates an item by its Key (for tests / manual injection).
func (a *Awakened) Upsert(it AwakenedItem) {
	if it.Key == "" {
		return
	}
	a.mu.Lock()
	cp := it
	a.items[it.Key] = &cp
	a.mu.Unlock()
	a.fire()
}

// Reset clears all tracked items.
func (a *Awakened) Reset() {
	a.mu.Lock()
	a.items = map[string]*AwakenedItem{}
	a.mu.Unlock()
	a.fire()
}

// SyncBody returns the awakened inventory as the /user/awakened/sync body
// (soulId-keyed items). Items without a soulId (soul packet not yet seen) are
// skipped since they have no stable key.
func (a *Awakened) SyncBody() any {
	snap := a.Snapshot()
	items := make([]map[string]any, 0, len(snap.Items))
	for _, it := range snap.Items {
		if it.SoulID == "" {
			continue
		}
		traits := make([]map[string]any, 0, len(it.Traits))
		for _, tr := range it.Traits {
			// v2 contract: send the raw roll only; backend computes the value.
			traits = append(traits, map[string]any{"id": tr.ID, "roll": tr.Roll})
		}
		m := map[string]any{
			"soulId":     it.SoulID,
			"itemId":     it.ItemID,
			"quality":    it.Quality,
			"strain":     round2(it.Strain),
			"attunement": round2(it.Attunement),
			"traits":     traits,
		}
		if it.AttunedTo != "" {
			m["attunedTo"] = it.AttunedTo
		}
		if it.Location != "" {
			m["location"] = it.Location
		}
		items = append(items, m)
	}
	return map[string]any{"v": 2, "items": items}
}

// ResolvedTrait is a backend-computed trait value applied back onto the table.
type ResolvedTrait struct {
	ID      string
	Value   float64
	Percent bool
}

// ApplyValues writes backend-computed trait values (from the sync response) onto
// every tracked item with the given soulId, matching by trait id, and marks them
// resolved. Fires a change so the UI repaints skeletons into values. The uploaded
// body (id + roll) is unaffected, so this never triggers a re-sync loop.
func (a *Awakened) ApplyValues(soulID string, traits []ResolvedTrait) {
	if soulID == "" || len(traits) == 0 {
		return
	}
	byID := make(map[string]ResolvedTrait, len(traits))
	for _, t := range traits {
		byID[t.ID] = t
	}
	a.mu.Lock()
	changed := false
	for _, it := range a.items {
		if it.SoulID != soulID {
			continue
		}
		for i := range it.Traits {
			if rv, ok := byID[it.Traits[i].ID]; ok {
				it.Traits[i].Value = rv.Value
				it.Traits[i].Percent = rv.Percent
				it.Traits[i].Resolved = true
				changed = true
			}
		}
	}
	a.mu.Unlock()
	if changed {
		a.fire()
	}
}

// Snapshot returns tracked items deduped by stable Key (soulId), skipping
// user-removed ones, sorted by name then quality.
func (a *Awakened) Snapshot() AwakenedSnapshot {
	a.mu.Lock()
	defer a.mu.Unlock()
	// Several objectId records can share one soulId (the same item re-observed under
	// a new objectId). Merge deterministically: the most recently updated packet wins
	// the volatile fields (attunement/strain/traits/...), and static item details are
	// backfilled from whichever record has them. Without this the display flickers
	// between two records' attunement values on random map iteration order.
	byKey := make(map[string]AwakenedItem, len(a.items))
	for _, it := range a.items {
		cur := *it
		if prev, ok := byKey[cur.Key]; ok {
			newer, older := cur, prev
			if prev.updated > cur.updated {
				newer, older = prev, cur
			}
			if newer.ItemID == "" {
				newer.ItemID = older.ItemID
			}
			if newer.Name == "" {
				newer.Name = older.Name
			}
			if newer.Quality == 0 {
				newer.Quality = older.Quality
			}
			if newer.AttunedTo == "" {
				newer.AttunedTo = older.AttunedTo
			}
			if newer.Location == "" {
				newer.Location = older.Location
			}
			cur = newer
		}
		byKey[cur.Key] = cur
	}
	out := make([]AwakenedItem, 0, len(byKey))
	for _, it := range byKey {
		out = append(out, it)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Name != out[j].Name {
			return out[i].Name < out[j].Name
		}
		return out[i].Quality > out[j].Quality
	})
	return AwakenedSnapshot{Items: out}
}
