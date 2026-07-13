package trackers

import (
	"math"
	"sort"
	"sync"
	"time"

	"github.com/albionmarketgg/data-client/internal/dispatch"
	"github.com/albionmarketgg/data-client/internal/photon"
)

// CombatParticipant is one player's contribution within an encounter (mobs are
// excluded — their hits on players surface as those players' DamageTaken).
type CombatParticipant struct {
	Name        string  `json:"name"`
	DamageDealt int64   `json:"damageDealt"`
	DamageTaken int64   `json:"damageTaken"`
	HealingDone int64   `json:"healingDone"`
	DPS         float64 `json:"dps"`
	IsParty     bool    `json:"isParty"`
	IsSelf      bool    `json:"isSelf"`
}

// CombatEncounter is a single fight summary.
type CombatEncounter struct {
	Number          int                 `json:"number"`
	DurationSeconds int64               `json:"durationSeconds"`
	Active          bool                `json:"active"`
	Fame            int64               `json:"fame"`
	Silver          int64               `json:"silver"`
	Participants    []CombatParticipant `json:"participants"`
}

// CombatSnapshot lists recent encounters (most recent first).
type CombatSnapshot struct {
	Encounters []CombatEncounter `json:"encounters"`
}

// CombatSummary aggregates every encounter this session into one per-player table
// (for the solo/group damage view).
type CombatSummary struct {
	DurationSeconds int64               `json:"durationSeconds"`
	Participants    []CombatParticipant `json:"participants"`
}

// DungeonEvent is one resolved dungeon gain: a fame event paired with the
// nearest silver event. Fame with no nearby silver still fires (Silver == 0);
// silver with no nearby fame never fires.
type DungeonEvent struct {
	Fame   int64     `json:"fame"`
	Silver int64     `json:"silver"`
	At     time.Time `json:"at"`
}

const combatIdle = 7 * time.Second

// dungeonPairWindow is how long to wait after a fame event to pair it with the
// closest silver event (silver can land before or after fame, within ~hundreds
// of ms). var, not const, so tests can shorten it.
var dungeonPairWindow = 2 * time.Second

type silverHit struct {
	at     time.Time
	amount int64
	paired bool
}

type entityAgg struct {
	damage  int64 // damage this entity dealt
	taken   int64 // damage this entity received
	healing int64 // healing this entity did
}

// famePoint is one fame gain with its timestamp (for the fame-over-time chart).
type famePoint struct {
	at     time.Time
	amount int64
}

// FamePoint is a chartable cumulative-fame sample (t = ms since the first gain).
type FamePoint struct {
	T          int64 `json:"t"`
	Cumulative int64 `json:"cumulative"`
}

// FameSeries is the fame-over-time series plus headline figures.
type FameSeries struct {
	Points      []FamePoint `json:"points"`
	Total       int64       `json:"total"`
	FamePerHour int64       `json:"famePerHour"`
}

type encounter struct {
	number   int
	started  time.Time
	lastAt   time.Time
	ended    bool
	fame     int64
	silver   int64
	byEntity map[int64]*entityAgg
}

// Combat is a damage-meter style tracker.
type Combat struct {
	party   *Party
	mobName func(index int) (string, bool) // resolve a NewMob index to a name (may be nil)
	selfID  func() (int64, bool)           // local player's object id (may be nil)

	mu         sync.Mutex
	names      map[int64]string
	selfSeen   map[int64]bool   // object ids that have been the local player (survive zone changes)
	memberName map[int64]string // object id -> stable group-player name (self or party)
	encs       []*encounter
	active     *encounter
	nextNum    int
	fameLog    []famePoint // fame gains over time (for the fame chart)
	onChange   func()
	onFeed     func(kind, detail string, count int)

	// dungeon-event pairing: each fame event is held briefly, matched with the
	// nearest silver, then emitted via onDungeon.
	pairMu     sync.Mutex
	pendSilver []*silverHit
	onDungeon  func(DungeonEvent)
}

// OnFeed registers a live-feed callback for notable combat events.
func (c *Combat) OnFeed(fn func(kind, detail string, count int)) { c.onFeed = fn }

// OnDungeonEvent registers a callback fired once per fame event, paired with the
// nearest silver. Silver-only activity never fires it.
func (c *Combat) OnDungeonEvent(fn func(DungeonEvent)) { c.onDungeon = fn }

// NewCombat creates a combat tracker. party, mobName and selfID may all be nil
// (mobs then show as "Mob <index>" and no participant is flagged as self).
func NewCombat(party *Party, mobName func(int) (string, bool), selfID func() (int64, bool)) *Combat {
	return &Combat{party: party, mobName: mobName, selfID: selfID, names: map[int64]string{}, selfSeen: map[int64]bool{}, memberName: map[int64]string{}, nextNum: 1}
}

// OnChange registers a change callback.
func (c *Combat) OnChange(fn func()) { c.onChange = fn }

// Register attaches combat handlers.
func (c *Combat) Register(d *dispatch.Dispatcher) {
	d.OnEvent(photon.EvNewCharacter, c.onNewCharacter)
	d.OnEvent(photon.EvNewMob, c.onNewMob)
	d.OnEvent(photon.EvHealthUpdate, c.onHealth)
	d.OnEvent(photon.EvHealthUpdates, c.onHealthBatch)
	d.OnEvent(photon.EvUpdateFame, c.onFame)
	d.OnEvent(photon.EvTakeSilver, c.onSilver)
}

func (c *Combat) onNewCharacter(m map[byte]any) {
	if obj := i64(m[0]); obj != 0 {
		if n := str(m[1]); n != "" {
			c.mu.Lock()
			c.names[obj] = n
			c.mu.Unlock()
		}
	}
}

func (c *Combat) onNewMob(m map[byte]any) {
	if obj := i64(m[0]); obj != 0 {
		idx := iN(m[1])
		name := "Mob " + itoa(idx)
		if c.mobName != nil {
			if n, ok := c.mobName(idx); ok {
				name = n
			}
		}
		c.mu.Lock()
		if _, ok := c.names[obj]; !ok {
			c.names[obj] = name
		}
		c.mu.Unlock()
	}
}

func (c *Combat) ensureEncounter(now time.Time) *encounter {
	if c.active != nil && now.Sub(c.active.lastAt) <= combatIdle {
		return c.active
	}
	if c.active != nil {
		c.active.ended = true
	}
	e := &encounter{number: c.nextNum, started: now, lastAt: now, byEntity: map[int64]*entityAgg{}}
	c.nextNum++
	c.encs = append(c.encs, e)
	if len(c.encs) > 50 {
		c.encs = c.encs[len(c.encs)-50:]
	}
	c.active = e
	return e
}

func (c *Combat) applyHealth(target, causer int64, change float64, now time.Time) {
	if change == 0 {
		return
	}
	amount := int64(math.Round(math.Abs(change)))
	if amount == 0 {
		return
	}
	c.mu.Lock()
	e := c.ensureEncounter(now)
	e.lastAt = now
	// Remember who is us / our party under this zone's object ids, so the session
	// summary can merge the same player across zone changes (ids are per-zone).
	c.noteGroup(causer)
	c.noteGroup(target)
	agg := e.byEntity[causer]
	if agg == nil {
		agg = &entityAgg{}
		e.byEntity[causer] = agg
	}
	if change < 0 {
		agg.damage += amount // causer dealt damage
		// credit the same hit as damage taken by the target
		ta := e.byEntity[target]
		if ta == nil {
			ta = &entityAgg{}
			e.byEntity[target] = ta
		}
		ta.taken += amount
	} else {
		agg.healing += amount
	}
	c.mu.Unlock()
}

func (c *Combat) onHealth(m map[byte]any) {
	change, _ := toFloat(m[2])
	c.applyHealth(i64(m[0]), i64(m[6]), change, time.Now())
	c.fire()
}

func (c *Combat) onHealthBatch(m map[byte]any) {
	target := i64(m[0])
	changes := toFloatSlice(m[2])
	causers := dispatch.Int64Slice(m[6])
	now := time.Now()
	for i := range changes {
		var causer int64
		if i < len(causers) {
			causer = causers[i]
		}
		c.applyHealth(target, causer, changes[i], now)
	}
	c.fire()
}

// fameIsPremium reports whether the UpdateFame premium-bonus flag (param 5) is
// set. Real packets carry a bool; be tolerant of a numeric 0/1 too.
func fameIsPremium(v any) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return i64(v) != 0
}

func (c *Combat) onFame(m map[byte]any) {
	// Gained fame = zone-multiplied fame (param 2) + the +50% Premium bonus when
	// active (param 5) + satchel bonus (param 10); all fixed-point /10000, matching
	// the in-game number. Premium was previously omitted, which undercounted fame
	// for premium players.
	base := i64(m[2])
	premium := int64(0)
	if base > 0 && fameIsPremium(m[5]) {
		premium = base / 2 // Premium grants +50% on the zone-multiplied fame
	}
	raw := base + premium + i64(m[10])
	fame := raw / 10000
	if fame <= 0 {
		return
	}
	now := time.Now()
	c.mu.Lock()
	e := c.ensureEncounter(now)
	e.fame += fame
	c.fameLog = append(c.fameLog, famePoint{at: now, amount: fame})
	if len(c.fameLog) > 5000 {
		c.fameLog = c.fameLog[len(c.fameLog)-5000:]
	}
	c.mu.Unlock()
	c.fire()
	if c.onFeed != nil {
		c.onFeed("dungeon", "fame +"+itoa(int(fame)), int(fame))
	}
	// hold briefly, then emit paired with the nearest silver
	if c.onDungeon != nil {
		time.AfterFunc(dungeonPairWindow, func() { c.resolveFame(fame, now) })
	}
}

func (c *Combat) onSilver(m map[byte]any) {
	yield := i64(m[3]) / 10000
	if yield <= 0 {
		return
	}
	now := time.Now()
	c.mu.Lock()
	e := c.ensureEncounter(now)
	e.silver += yield
	c.mu.Unlock()
	c.fire()
	if c.onFeed != nil {
		c.onFeed("dungeon", "silver +"+itoa(int(yield)), int(yield))
	}
	// buffer silver so a fame event can pair with it (silver may precede or follow
	// its fame); silver that never pairs is dropped.
	if c.onDungeon != nil {
		c.pairMu.Lock()
		c.pendSilver = append(c.pendSilver, &silverHit{at: now, amount: yield})
		c.prunePendSilver(now)
		c.pairMu.Unlock()
	}
}

// resolveFame matches a held fame event with the closest unpaired silver within
// the pairing window and emits the resulting dungeon event.
func (c *Combat) resolveFame(fame int64, at time.Time) {
	var silver int64
	c.pairMu.Lock()
	best := -1
	bestDelta := dungeonPairWindow + 1
	for i, s := range c.pendSilver {
		if s.paired {
			continue
		}
		d := absDur(s.at.Sub(at))
		if d <= dungeonPairWindow && d < bestDelta {
			best, bestDelta = i, d
		}
	}
	if best >= 0 {
		silver = c.pendSilver[best].amount
		c.pendSilver[best].paired = true
	}
	c.prunePendSilver(at)
	c.pairMu.Unlock()
	if c.onDungeon != nil {
		c.onDungeon(DungeonEvent{Fame: fame, Silver: silver, At: at})
	}
}

// prunePendSilver drops paired silvers and ones too old to pair with any future
// fame. Caller holds pairMu.
func (c *Combat) prunePendSilver(ref time.Time) {
	keep := c.pendSilver[:0]
	for _, s := range c.pendSilver {
		if s.paired || ref.Sub(s.at) > 2*dungeonPairWindow {
			continue
		}
		keep = append(keep, s)
	}
	c.pendSilver = keep
}

func absDur(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}

func (c *Combat) fire() {
	if c.onChange != nil {
		c.onChange()
	}
}

// Reset clears all encounters.
func (c *Combat) Reset() {
	c.mu.Lock()
	c.encs = nil
	c.active = nil
	c.nextNum = 1
	c.fameLog = nil
	c.mu.Unlock()
	c.fire()
}

// FameSeries returns the cumulative fame-over-time series plus total and an
// estimated fame/hour (over the span of fame gains, floored at 60s to tame early
// spikes). Empty when no fame has been gained.
func (c *Combat) FameSeries() FameSeries {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.fameLog) == 0 {
		return FameSeries{}
	}
	start := c.fameLog[0].at
	pts := make([]FamePoint, 0, len(c.fameLog)+1)
	pts = append(pts, FamePoint{T: 0, Cumulative: 0})
	var cum int64
	for _, f := range c.fameLog {
		cum += f.amount
		pts = append(pts, FamePoint{T: f.at.Sub(start).Milliseconds(), Cumulative: cum})
	}
	span := c.fameLog[len(c.fameLog)-1].at.Sub(start)
	if span < time.Minute {
		span = time.Minute // floor so a short burst doesn't read as a huge rate
	}
	perHour := int64(float64(cum) / span.Hours())
	return FameSeries{Points: pts, Total: cum, FamePerHour: perHour}
}

// noteGroup records an object id that currently belongs to the local player or a
// party member, caching its stable name so the summary can merge it across zone
// changes (mobs/others are not cached). Caller holds mu.
func (c *Combat) noteGroup(obj int64) {
	if obj == 0 {
		return
	}
	if _, ok := c.memberName[obj]; ok {
		return // already resolved
	}
	if c.selfID != nil {
		if sid, ok := c.selfID(); ok && sid == obj {
			c.selfSeen[obj] = true
			if c.party != nil {
				if ln := c.party.LocalName(); ln != "" {
					c.memberName[obj] = ln
				}
			}
			return
		}
	}
	if c.party != nil && c.party.IsPartyMember(obj) {
		if n, ok := c.party.NameByObject(obj); ok && n != "" {
			c.memberName[obj] = n
		}
	}
}

// resolveName returns the best display name for an object id, preferring the
// combat/party name caches; falls back to the local player's name for known self
// ids, else "entity:<id>". Caller holds mu.
func (c *Combat) resolveName(obj int64) string {
	if n := c.names[obj]; n != "" {
		return n
	}
	if c.party != nil {
		if n, ok := c.party.NameByObject(obj); ok && n != "" {
			return n
		}
	}
	if n := c.memberName[obj]; n != "" {
		return n
	}
	if c.selfSeen[obj] && c.party != nil {
		if ln := c.party.LocalName(); ln != "" {
			return ln
		}
	}
	return "Unknown" // unresolved object id (not a named player or mob)
}

// buildParticipants turns an entity aggregate into a sorted participant list for
// the per-fight log — all entities, including named mobs. Caller holds mu.
func (c *Combat) buildParticipants(byEntity map[int64]*entityAgg, dur float64, selfID int64, hasSelf bool) []CombatParticipant {
	var out []CombatParticipant
	for obj, agg := range byEntity {
		if agg.damage == 0 && agg.healing == 0 && agg.taken == 0 {
			continue
		}
		p := CombatParticipant{
			Name:        c.resolveName(obj),
			DamageDealt: agg.damage,
			DamageTaken: agg.taken,
			HealingDone: agg.healing,
			DPS:         float64(agg.damage) / dur,
			IsSelf:      (hasSelf && obj == selfID) || c.selfSeen[obj],
		}
		if p.IsSelf || (c.party != nil && c.party.IsPartyMember(obj)) {
			p.IsParty = true
		}
		out = append(out, p)
	}
	sort.Slice(out, func(a, b int) bool { return out[a].DamageDealt > out[b].DamageDealt })
	return out
}

// self returns the local player's object id (0, false if unknown).
func (c *Combat) self() (int64, bool) {
	if c.selfID == nil {
		return 0, false
	}
	return c.selfID()
}

// Snapshot returns recent encounters (most recent first), player participants
// only. Encounters with no player participants are omitted.
func (c *Combat) Snapshot() CombatSnapshot {
	selfID, hasSelf := c.self()
	c.mu.Lock()
	defer c.mu.Unlock()
	var out []CombatEncounter
	for i := len(c.encs) - 1; i >= 0; i-- {
		e := c.encs[i]
		dur := e.lastAt.Sub(e.started).Seconds()
		if dur < 1 {
			dur = 1
		}
		parts := c.buildParticipants(e.byEntity, dur, selfID, hasSelf)
		if len(parts) == 0 {
			continue // nothing happened — don't log an empty encounter
		}
		out = append(out, CombatEncounter{
			Number:          e.number,
			DurationSeconds: int64(dur),
			Active:          e == c.active && !e.ended,
			Fame:            e.fame,
			Silver:          e.silver,
			Participants:    parts,
		})
	}
	return CombatSnapshot{Encounters: out}
}

// SessionSummary aggregates all encounters this session into one per-player table
// (the group totals). It merges by stable player NAME so the same player survives
// zone changes (object ids are per-zone); only the local player + party members
// are included. DPS is over total combat time across encounters.
func (c *Combat) SessionSummary() CombatSummary {
	c.mu.Lock()
	defer c.mu.Unlock()
	localName := ""
	if c.party != nil {
		localName = c.party.LocalName()
	}
	byName := map[string]*entityAgg{}
	var total float64
	for _, e := range c.encs {
		d := e.lastAt.Sub(e.started).Seconds()
		if d < 1 {
			d = 1
		}
		total += d
		for obj, a := range e.byEntity {
			name := c.memberName[obj]
			if name == "" {
				if c.selfSeen[obj] && localName != "" {
					name = localName
				} else {
					continue // not a known group member (mob / other player)
				}
			}
			g := byName[name]
			if g == nil {
				g = &entityAgg{}
				byName[name] = g
			}
			g.damage += a.damage
			g.taken += a.taken
			g.healing += a.healing
		}
	}
	if total < 1 {
		total = 1
	}
	out := make([]CombatParticipant, 0, len(byName))
	for name, a := range byName {
		if a.damage == 0 && a.healing == 0 && a.taken == 0 {
			continue
		}
		out = append(out, CombatParticipant{
			Name:        name,
			DamageDealt: a.damage,
			DamageTaken: a.taken,
			HealingDone: a.healing,
			DPS:         float64(a.damage) / total,
			IsParty:     true,
			IsSelf:      name == localName,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].DamageDealt > out[j].DamageDealt })
	return CombatSummary{DurationSeconds: int64(total), Participants: out}
}

func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	}
	if i, ok := dispatch.Int64(v); ok {
		return float64(i), true
	}
	return 0, false
}

func toFloatSlice(v any) []float64 {
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]float64, 0, len(arr))
	for _, e := range arr {
		if f, ok := toFloat(e); ok {
			out = append(out, f)
		}
	}
	return out
}
