// Package locations resolves Albion market location ids to human names. The main
// market cities are built in; for full coverage set a LocationsURL (world.json) in
// config.
package locations

import "strings"

// builtin maps the common market location ids to display names.
var builtin = map[string]string{
	"7":    "Thetford",
	"0007": "Thetford",
	"1002": "Lymhurst",
	"2004": "Bridgewatch",
	"3008": "Martlock",
	"4002": "Fort Sterling",
	"3005": "Caerleon",
	"3003": "Caerleon Black Market",
	"5003": "Brecilien",
	// portal towns
	"301":  "Thetford Portal",
	"1301": "Lymhurst Portal",
	"2301": "Bridgewatch Portal",
	"3301": "Martlock Portal",
	"4301": "Fort Sterling Portal",
}

// Resolve returns a display name for a raw location id, trying the raw value, the
// segment before an '@', then a zero-trimmed form. Returns the raw value (and
// false) when unknown so the UI still shows something.
func Resolve(raw string) (string, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", false
	}
	if n, ok := builtin[raw]; ok {
		return n, true
	}
	if at := strings.IndexByte(raw, '@'); at > 0 {
		if n, ok := builtin[raw[:at]]; ok {
			return n, true
		}
	}
	trimmed := strings.TrimLeft(raw, "0")
	if trimmed == "" {
		trimmed = "0"
	}
	if n, ok := builtin[trimmed]; ok {
		return n, true
	}
	return raw, false
}
