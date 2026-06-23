package trackers

import (
	"math"
	"sort"
	"sync"
	"time"

	"github.com/niick1231/albionmarket_dataclient/internal/dispatch"
	"github.com/niick1231/albionmarket_dataclient/internal/photon"
)

// CombatParticipant is one entity's contribution within an encounter.
type CombatParticipant struct {
	Name        string  `json:"name"`
	DamageDealt int64   `json:"damageDealt"`
	HealingDone int64   `json:"healingDone"`
	DPS         float64 `json:"dps"`
	IsParty     bool    `json:"isParty"`
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
	damage  int64
	healing int64
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
	party *Party

	mu       sync.Mutex
	names    map[int64]string
	encs     []*encounter
	active   *encounter
	nextNum  int
	onChange func()
	onFeed   func(kind, detail string, count int)

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

// NewCombat creates a combat tracker. party may be nil.
func NewCombat(party *Party) *Combat {
	return &Combat{party: party, names: map[int64]string{}, nextNum: 1}
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
		c.mu.Lock()
		if _, ok := c.names[obj]; !ok {
			c.names[obj] = "Mob " + itoa(iN(m[1]))
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
	agg := e.byEntity[causer]
	if agg == nil {
		agg = &entityAgg{}
		e.byEntity[causer] = agg
	}
	if change < 0 {
		agg.damage += amount
	} else {
		agg.healing += amount
	}
	_ = target
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

func (c *Combat) onFame(m map[byte]any) {
	// Gained fame = zone-multiplied fame (param 2) + satchel bonus (param 10);
	// both fixed-point /10000, matching the in-game number.
	raw := i64(m[2]) + i64(m[10])
	fame := raw / 10000
	if fame <= 0 {
		return
	}
	now := time.Now()
	c.mu.Lock()
	e := c.ensureEncounter(now)
	e.fame += fame
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
	c.mu.Unlock()
	c.fire()
}

// Snapshot returns recent encounters, most recent first.
func (c *Combat) Snapshot() CombatSnapshot {
	c.mu.Lock()
	defer c.mu.Unlock()
	var out []CombatEncounter
	for i := len(c.encs) - 1; i >= 0; i-- {
		e := c.encs[i]
		dur := e.lastAt.Sub(e.started).Seconds()
		if dur < 1 {
			dur = 1
		}
		ce := CombatEncounter{
			Number:          e.number,
			DurationSeconds: int64(dur),
			Active:          e == c.active && !e.ended,
			Fame:            e.fame,
			Silver:          e.silver,
		}
		for obj, agg := range e.byEntity {
			if agg.damage == 0 && agg.healing == 0 {
				continue
			}
			name := c.names[obj]
			if name == "" {
				name = "entity:" + itoa(int(obj))
			}
			p := CombatParticipant{
				Name:        name,
				DamageDealt: agg.damage,
				HealingDone: agg.healing,
				DPS:         float64(agg.damage) / dur,
			}
			if c.party != nil {
				p.IsParty = c.party.IsPartyMember(obj)
			}
			ce.Participants = append(ce.Participants, p)
		}
		sort.Slice(ce.Participants, func(a, b int) bool {
			return ce.Participants[a].DamageDealt > ce.Participants[b].DamageDealt
		})
		out = append(out, ce)
	}
	return CombatSnapshot{Encounters: out}
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
