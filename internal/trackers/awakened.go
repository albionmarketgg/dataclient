package trackers

import (
	"encoding/hex"
	"math"
	"sort"
	"strconv"
	"sync"

	"github.com/niick1231/albionmarket_dataclient/internal/dispatch"
	"github.com/niick1231/albionmarket_dataclient/internal/photon"
)

// traitBase is each trait's base value (from the legendary spell data). The
// displayed value is base × rollFactor, where rollFactor = minfactor +
// raw×(maxfactor−minfactor) with the game's single "default" config
// (minfactor 0.01, maxfactor 1). base >= 1 → a flat stat (rounded int); base < 1
// → a percentage (×100).
var traitBase = map[string]float64{
	"TRAIT_ITEM_POWER":             120,
	"TRAIT_HITPOINTS_MAX":          260,
	"TRAIT_ENERGY_MAX":             52,
	"TRAIT_CC_RESIST":              16,
	"TRAIT_THREAT_BONUS":           2,
	"TRAIT_ENERGY_COST_REDUCTION":  0.24,
	"TRAIT_HEALING_DEALT":          0.15,
	"TRAIT_HEALING_RECEIVED":       0.15,
	"TRAIT_ATTACK_SPEED":           0.24,
	"TRAIT_ATTACK_RANGE":           0.20,
	"TRAIT_CAST_SPEED_INCREASE":    0.20,
	"TRAIT_COOLDOWN_REDUCTION":     0.12,
	"TRAIT_MOB_FAME":               0.30,
	"TRAIT_RESILIENCE_PENETRATION": 0.14,
	"TRAIT_DEFENSE_BONUS":          0.14,
	"TRAIT_ABILITY_DAMAGE":         0.165,
	"TRAIT_AUTO_ATTACK_DAMAGE":     0.24,
	"TRAIT_CC_DURATION":            0.32,
	"TRAIT_LIFESTEAL":              0.10,
}

func round2(f float64) float64 { return math.Round(f*100) / 100 }

// computeTraitValue turns the raw packet roll (0-1) into the real trait value.
// Returns the value and whether it's a percentage.
func computeTraitValue(id string, raw float64) (value float64, percent bool) {
	base, ok := traitBase[id]
	if !ok {
		return raw, false
	}
	scaled := base * (0.01 + raw*0.99)
	if base >= 1 { // flat stat, 2 dp
		return round2(scaled), false
	}
	return round2(scaled * 100), true // percentage, 2 dp
}

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
// (e.g. "TRAIT_ITEM_POWER"); Name is resolved locally for display only — uploads
// send the id and the backend resolves the label.
type AwakenedTrait struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Value   float64 `json:"value"`   // resolved display value (flat int or % number)
	Percent bool    `json:"percent"` // display/interpret as a percentage
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
	hidden   map[string]bool          // keys the user removed from the table
	onChange func()
}

// NewAwakened creates the tracker. info resolves item names; serverID/playerName/
// location supply the current region, local character, and mapped location (all
// may be nil).
func NewAwakened(info ItemInfo, serverID func() int, playerName func() string, location func() string) *Awakened {
	if info == nil {
		info = nopItems{}
	}
	return &Awakened{info: info, serverID: serverID, playerName: playerName, location: location, items: map[string]*AwakenedItem{}, hidden: map[string]bool{}}
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
		v := 0.0
		if i < len(vals) {
			v = vals[i]
		}
		val, pct := computeTraitValue(id, v)
		traits = append(traits, AwakenedTrait{ID: id, Name: traitName(id), Value: val, Percent: pct})
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

// Remove hides an item from the table by its Key (so re-capture won't re-add it).
func (a *Awakened) Remove(key string) {
	if key == "" {
		return
	}
	a.mu.Lock()
	a.hidden[key] = true
	for oid, it := range a.items {
		if it.Key == key {
			delete(a.items, oid)
		}
	}
	a.mu.Unlock()
	a.fire()
}

// Reset clears all tracked items and un-hides everything.
func (a *Awakened) Reset() {
	a.mu.Lock()
	a.items = map[string]*AwakenedItem{}
	a.hidden = map[string]bool{}
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
			traits = append(traits, map[string]any{"id": tr.ID, "value": tr.Value})
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
	return map[string]any{"items": items}
}

// Snapshot returns tracked items deduped by stable Key (soulId), skipping
// user-removed ones, sorted by name then quality.
func (a *Awakened) Snapshot() AwakenedSnapshot {
	a.mu.Lock()
	defer a.mu.Unlock()
	byKey := make(map[string]AwakenedItem, len(a.items))
	for _, it := range a.items {
		if a.hidden[it.Key] {
			continue
		}
		// prefer the most complete record for a given key (item details present)
		if prev, ok := byKey[it.Key]; ok && prev.ItemID != "" && it.ItemID == "" {
			continue
		}
		byKey[it.Key] = *it
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
