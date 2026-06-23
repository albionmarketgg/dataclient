package handlers

import (
	"sort"

	"github.com/niick1231/albionmarket_dataclient/internal/dispatch"
	"github.com/niick1231/albionmarket_dataclient/internal/photon"
)

// SpecResolver maps a Destiny Board achievement index to its stable string id.
type SpecResolver interface {
	IDByIndex(i int) (string, bool)
	Version() string
}

// SpecEntry is one Destiny Board node and its mastery level (1-100).
type SpecEntry struct {
	ID    string `json:"id"`
	Level int    `json:"level"`
}

// Specs decodes the FullAchievementInfo packet (the mastery/specialization board
// pushed on entering the world) and emits resolved {id, level} entries for upload.
type Specs struct {
	res    SpecResolver
	onFull func(version string, entries []SpecEntry)
}

// NewSpecs builds the handler. res resolves index->id (may be nil → no output).
func NewSpecs(res SpecResolver) *Specs { return &Specs{res: res} }

// OnFull registers the callback fired once per FullAchievementInfo packet.
func (s *Specs) OnFull(fn func(version string, entries []SpecEntry)) { s.onFull = fn }

// Register attaches the FullAchievementInfo handler.
func (s *Specs) Register(d *dispatch.Dispatcher) {
	d.OnEvent(photon.EvFullAchievementInfo, s.handleFull)
}

// handleFull reads the parallel index/level arrays (param 2 = indices, param 3 =
// levels) plus the maxed-node indices (param 1 → level 100), resolves each index
// to an id, and emits the full board (level > 0 only).
func (s *Specs) handleFull(m map[byte]any) {
	if s.res == nil || s.onFull == nil {
		return
	}
	indices := dispatch.Int64Slice(m[2])
	levels := byteSlice(m[3])
	level100 := dispatch.Int64Slice(m[1])

	byIndex := make(map[int]int, len(indices))
	n := len(indices)
	if len(levels) < n {
		n = len(levels)
	}
	for i := 0; i < n; i++ {
		byIndex[int(indices[i])] = int(levels[i])
	}
	// param 1 lists fully-maxed nodes — force them to 100.
	for _, idx := range level100 {
		byIndex[int(idx)] = 100
	}

	entries := make([]SpecEntry, 0, len(byIndex))
	for idx, lvl := range byIndex {
		if lvl <= 0 {
			continue
		}
		if id, ok := s.res.IDByIndex(idx); ok {
			entries = append(entries, SpecEntry{ID: id, Level: lvl})
		}
	}
	if len(entries) == 0 {
		return
	}
	sort.Slice(entries, func(a, b int) bool { return entries[a].ID < entries[b].ID })
	s.onFull(s.res.Version(), entries)
}

// byteSlice extracts a []byte param. The mastery levels arrive as a raw byte array.
func byteSlice(v any) []byte {
	switch b := v.(type) {
	case []byte:
		return b
	case []any:
		out := make([]byte, 0, len(b))
		for _, e := range b {
			if n, ok := dispatch.Int64(e); ok {
				out = append(out, byte(n))
			}
		}
		return out
	}
	return nil
}
