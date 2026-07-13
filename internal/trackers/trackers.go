// Package trackers implements the gameplay trackers (party, gathering, combat,
// loot). They consume parsed Photon events via the dispatcher and expose JSON
// snapshots for the UI.
//
// These are pragmatic: faithful to the event codes and parameter indices, with
// simplified correlation/timing (noted inline). The market-data pipeline is the
// precise part; these are extras.
package trackers

import (
	"encoding/hex"

	"github.com/albionmarketgg/dataclient/internal/dispatch"
	"github.com/albionmarketgg/dataclient/internal/photon"
)

// ItemInfo resolves item names and values for the trackers.
type ItemInfo interface {
	UniqueName(id int) (string, bool)
	DisplayName(id int) (string, bool)
	EMV(id, quality int) (int64, bool)
}

// nopItems is used when no item service is supplied.
type nopItems struct{}

func (nopItems) UniqueName(int) (string, bool)   { return "", false }
func (nopItems) DisplayName(int) (string, bool)  { return "", false }
func (nopItems) EMV(int, int) (int64, bool)      { return 0, false }

func itemName(info ItemInfo, id int) string {
	if n, ok := info.DisplayName(id); ok && n != "" {
		return n
	}
	if n, ok := info.UniqueName(id); ok && n != "" {
		return n
	}
	return "item:" + itoa(id)
}

func itemValue(info ItemInfo, id, quality int) int64 {
	v, _ := info.EMV(id, quality)
	return v
}

// itemUnique returns an item's canonical unique name (e.g. "T5_MAIN_MACE@2",
// enchant baked in), or "" when unresolved.
func itemUnique(info ItemInfo, id int) string {
	if n, ok := info.UniqueName(id); ok {
		return n
	}
	return ""
}

// guidKey converts a Photon GUID parameter (CustomData/bytes/string) to a stable
// string key.
func guidKey(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return hex.EncodeToString(x)
	case photon.CustomData:
		return hex.EncodeToString(x.Data)
	}
	return ""
}

// guidKeys converts a GUID array parameter to string keys.
func guidKeys(v any) []string {
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(arr))
	for _, e := range arr {
		if k := guidKey(e); k != "" {
			out = append(out, k)
		}
	}
	return out
}

func i64(v any) int64 { n, _ := dispatch.Int64(v); return n }
func iN(v any) int    { n, _ := dispatch.Int(v); return n }
func str(v any) string {
	s, _ := dispatch.String(v)
	return s
}

func itoa(i int) string {
	// local helper to avoid importing strconv everywhere
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var b [20]byte
	p := len(b)
	for i > 0 {
		p--
		b[p] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		p--
		b[p] = '-'
	}
	return string(b[p:])
}
