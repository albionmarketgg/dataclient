// Package state holds the live, observable player/session state shared across
// the engine and surfaced to the UI.
package state

import (
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/albionmarketgg/dataclient/internal/config"
)

// Character is a detected in-game character (from the Join packet).
type Character struct {
	Name       string    `json:"name"`
	Fame       int64     `json:"fame"` // already divided by 10000
	ServerID   int       `json:"serverId"`
	ServerName string    `json:"serverName"`
	At         time.Time `json:"at"`
}

// Location identifies the player's current market location.
type Location struct {
	ID   string
	Name string
}

// Snapshot is an immutable view of the player state for the UI.
type Snapshot struct {
	ServerID         int       `json:"serverId"`
	ServerName       string    `json:"serverName"`
	LocationID       string    `json:"locationId"`
	LocationName     string    `json:"locationName"`
	PlayerName       string    `json:"playerName"`
	InGame           bool      `json:"inGame"`
	HasEncryptedData bool      `json:"hasEncryptedData"`
	LastPacket       time.Time `json:"lastPacket"`
	Listening        bool      `json:"listening"`
}

// State is the concurrency-safe player/session state.
type State struct {
	mu sync.RWMutex

	server       *config.Server
	location     Location
	playerName   string
	userObjectID int64
	hasUserObj   bool
	lastPacket   time.Time
	hasEncrypted bool
	listening    bool
	characters   map[string]Character

	onChange func(Snapshot)
}

// New creates an empty state.
func New() *State { return &State{characters: map[string]Character{}} }

// RecordCharacter stores a detected own-character, keyed by server+name so each
// appears once with its latest fame.
func (s *State) RecordCharacter(name string, fame int64) {
	if name == "" {
		return
	}
	s.mu.Lock()
	srvID, srvName := 0, ""
	if s.server != nil {
		srvID, srvName = s.server.ID, s.server.Name
	}
	if s.characters == nil {
		s.characters = map[string]Character{}
	}
	s.characters[strconv.Itoa(srvID)+"/"+name] = Character{
		Name: name, Fame: fame, ServerID: srvID, ServerName: srvName, At: time.Now(),
	}
	s.mu.Unlock()
	s.notify()
}

// Characters returns detected characters, most recent first.
func (s *State) Characters() []Character {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Character, 0, len(s.characters))
	for _, c := range s.characters {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].At.After(out[j].At) })
	return out
}

// OnChange registers a callback fired (outside the lock) on any state change.
func (s *State) OnChange(fn func(Snapshot)) { s.onChange = fn }

func (s *State) notify() {
	if s.onChange != nil {
		s.onChange(s.snapshotLocked())
	}
}

func (s *State) snapshotLocked() Snapshot {
	snap := Snapshot{
		LocationID:       s.location.ID,
		LocationName:     s.location.Name,
		PlayerName:       s.playerName,
		InGame:           time.Since(s.lastPacket) < 10*time.Second && !s.lastPacket.IsZero(),
		HasEncryptedData: s.hasEncrypted,
		LastPacket:       s.lastPacket,
		Listening:        s.listening,
	}
	if s.server != nil {
		snap.ServerID = s.server.ID
		snap.ServerName = s.server.Name
	}
	return snap
}

// Snapshot returns the current state snapshot.
func (s *State) Snapshot() Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.snapshotLocked()
}

// Server returns the detected server (may be nil).
func (s *State) Server() *config.Server {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.server
}

// SetServer records the detected game server.
func (s *State) SetServer(srv *config.Server) {
	s.mu.Lock()
	changed := s.server == nil || s.server.ID != srv.ID
	s.server = srv
	s.mu.Unlock()
	if changed {
		s.notify()
	}
}

// MarkPacket records that a packet was just observed.
func (s *State) MarkPacket() {
	s.mu.Lock()
	wasOut := time.Since(s.lastPacket) >= 10*time.Second || s.lastPacket.IsZero()
	s.lastPacket = time.Now()
	s.mu.Unlock()
	if wasOut {
		s.notify()
	}
}

// SetEncrypted records whether encrypted market data was seen.
func (s *State) SetEncrypted(v bool) {
	s.mu.Lock()
	changed := s.hasEncrypted != v
	s.hasEncrypted = v
	s.mu.Unlock()
	if changed {
		s.notify()
	}
}

// SetLocation records the player's market location.
func (s *State) SetLocation(loc Location) {
	s.mu.Lock()
	changed := s.location.ID != loc.ID
	s.location = loc
	s.mu.Unlock()
	if changed {
		s.notify()
	}
}

// Location returns the current location.
func (s *State) Location() Location {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.location
}

// SetPlayerName records the player's character name.
func (s *State) SetPlayerName(name string) {
	s.mu.Lock()
	changed := s.playerName != name
	s.playerName = name
	s.mu.Unlock()
	if changed {
		s.notify()
	}
}

// SetUserObjectID records the local player's object id (from Join).
func (s *State) SetUserObjectID(id int64) {
	s.mu.Lock()
	s.userObjectID = id
	s.hasUserObj = true
	s.mu.Unlock()
}

// UserObjectID returns the local player's object id and whether it is known.
func (s *State) UserObjectID() (int64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.userObjectID, s.hasUserObj
}

// SetListening records capture status.
func (s *State) SetListening(v bool) {
	s.mu.Lock()
	s.listening = v
	s.mu.Unlock()
	s.notify()
}

// InGame reports whether a packet was seen recently.
func (s *State) InGame() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return !s.lastPacket.IsZero() && time.Since(s.lastPacket) < 10*time.Second
}

// OkToUpload gates market uploads: location set, in game, server known.
func (s *State) OkToUpload() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.location.ID != "" && s.server != nil &&
		!s.lastPacket.IsZero() && time.Since(s.lastPacket) < 10*time.Second
}
