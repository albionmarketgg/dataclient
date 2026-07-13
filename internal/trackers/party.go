package trackers

import (
	"sort"
	"sync"

	"github.com/albionmarketgg/dataclient/internal/dispatch"
	"github.com/albionmarketgg/dataclient/internal/photon"
)

// PartyMember is one party member in a snapshot.
type PartyMember struct {
	ObjectID int64  `json:"objectId"`
	Name     string `json:"name"`
	IsLocal  bool   `json:"isLocal"`
}

// PartySnapshot is the current party state.
type PartySnapshot struct {
	Members []PartyMember `json:"members"`
}

type partyEntry struct {
	objectID int64
	name     string
}

// Party tracks the player's party membership.
type Party struct {
	mu        sync.RWMutex
	members   map[string]*partyEntry // by guid key
	known     map[string]*partyEntry // sticky cache
	localGuid string
	localName string
	localObj  int64
	onChange  func()
	onFeed    func(kind, detail string, count int)
}

// OnFeed registers a live-feed callback for party changes.
func (p *Party) OnFeed(fn func(kind, detail string, count int)) { p.onFeed = fn }

// NewParty creates a party tracker.
func NewParty() *Party {
	return &Party{members: map[string]*partyEntry{}, known: map[string]*partyEntry{}}
}

// OnChange registers a change callback.
func (p *Party) OnChange(fn func()) { p.onChange = fn }

// Register attaches party handlers.
func (p *Party) Register(d *dispatch.Dispatcher) {
	d.OnEvent(photon.EvNewCharacter, p.onNewCharacter)
	d.OnEvent(photon.EvPartyJoined, p.onPartyJoined)
	d.OnEvent(photon.EvPartyPlayerJoined, p.onPlayerJoined)
	d.OnEvent(photon.EvPartyPlayerLeft, p.onPlayerLeft)
	d.OnEvent(photon.EvPartyDisbanded, p.onDisbanded)
	d.OnEvent(photon.EvPartyOnClusterPartyJoined, p.onClusterJoined)
	d.OnEvent(photon.EvPartySetRoleFlag, p.onSetRole)
	d.OnResponse(photon.OpJoin, p.onJoin)
}

func (p *Party) notify() {
	if p.onChange != nil {
		p.onChange()
	}
}

func (p *Party) onJoin(_ int16, _ string, m map[byte]any) {
	p.mu.Lock()
	p.localObj = i64(m[0])
	p.localGuid = guidKey(m[1])
	p.localName = str(m[2])
	p.mu.Unlock()
	p.notify()
}

func (p *Party) onNewCharacter(m map[byte]any) {
	guid := guidKey(m[7])
	name := str(m[1])
	obj := i64(m[0])
	if guid == "" || name == "" {
		return
	}
	p.mu.Lock()
	e := &partyEntry{objectID: obj, name: name}
	p.known[guid] = e
	if _, inParty := p.members[guid]; inParty {
		p.members[guid] = e
	}
	p.mu.Unlock()
}

func (p *Party) onPartyJoined(m map[byte]any) {
	guids := guidKeys(m[5])
	names := dispatch.Strings(m[6])
	p.mu.Lock()
	p.members = map[string]*partyEntry{}
	for i := 0; i < len(guids) && i < len(names); i++ {
		if guids[i] == "" || names[i] == "" {
			continue
		}
		p.members[guids[i]] = &partyEntry{name: names[i]}
		p.known[guids[i]] = p.members[guids[i]]
	}
	p.mu.Unlock()
	p.notify()
}

func (p *Party) onPlayerJoined(m map[byte]any) {
	guid := guidKey(m[1])
	name := str(m[2])
	if guid == "" || name == "" {
		return
	}
	p.mu.Lock()
	p.members[guid] = &partyEntry{name: name}
	p.known[guid] = p.members[guid]
	p.mu.Unlock()
	if p.onFeed != nil {
		p.onFeed("party", name+" joined the party", 1)
	}
	p.notify()
}

func (p *Party) onPlayerLeft(m map[byte]any) {
	guid := guidKey(m[1])
	if guid == "" {
		return
	}
	p.mu.Lock()
	delete(p.members, guid)
	p.mu.Unlock()
	p.notify()
}

func (p *Party) onDisbanded(_ map[byte]any) {
	p.mu.Lock()
	p.members = map[string]*partyEntry{}
	p.mu.Unlock()
	p.notify()
}

func (p *Party) onClusterJoined(m map[byte]any) {
	guids := guidKeys(m[0])
	p.mu.Lock()
	for _, g := range guids {
		if g == "" || g == p.localGuid {
			continue
		}
		if _, ok := p.members[g]; !ok {
			if k, ok := p.known[g]; ok {
				p.members[g] = k
			} else {
				p.members[g] = &partyEntry{}
			}
		}
	}
	p.mu.Unlock()
	p.notify()
}

func (p *Party) onSetRole(m map[byte]any) {
	guid := guidKey(m[1])
	if guid == "" || guid == p.localGuid {
		return
	}
	p.mu.Lock()
	if _, ok := p.members[guid]; !ok {
		if k, ok := p.known[guid]; ok {
			p.members[guid] = k
		} else {
			p.members[guid] = &partyEntry{}
		}
	}
	p.mu.Unlock()
	p.notify()
}

// LocalName returns the local player's character name (empty until Join).
func (p *Party) LocalName() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.localName
}

// NameByObject resolves an object id to a known player name — the local player or
// a party member whose character we've linked to an object id. ok=false if unknown.
func (p *Party) NameByObject(objectID int64) (string, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if objectID != 0 && objectID == p.localObj && p.localName != "" {
		return p.localName, true
	}
	for _, e := range p.members {
		if e.objectID == objectID && e.name != "" {
			return e.name, true
		}
	}
	return "", false
}

// IsPartyMember reports whether an object id belongs to a tracked party member
// or the local player.
func (p *Party) IsPartyMember(objectID int64) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if objectID == p.localObj {
		return true
	}
	for _, e := range p.members {
		if e.objectID == objectID {
			return true
		}
	}
	return false
}

// Snapshot returns the current party state (local player first, then by name).
func (p *Party) Snapshot() PartySnapshot {
	p.mu.RLock()
	defer p.mu.RUnlock()
	var out []PartyMember
	if p.localName != "" {
		out = append(out, PartyMember{ObjectID: p.localObj, Name: p.localName, IsLocal: true})
	}
	var rest []PartyMember
	for _, e := range p.members {
		if e.name == "" || e.name == p.localName {
			continue
		}
		rest = append(rest, PartyMember{ObjectID: e.objectID, Name: e.name})
	}
	sort.Slice(rest, func(i, j int) bool { return rest[i].Name < rest[j].Name })
	return PartySnapshot{Members: append(out, rest...)}
}
