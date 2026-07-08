// Package locations resolves Albion location ids to human names for display only
// (the player's current cluster, trade/mail locations, awakened item location).
// Market/trade *data* keys off the numeric id separately, so these names never
// affect uploaded data. Coverage is the curated set of towns, markets, banks,
// portals, Brecilien districts and personal/guild islands; unknown ids fall back
// to the raw value so the UI still shows something.
package locations

import "strings"

// builtin maps location ids to display names. Numeric keys are the 4-digit padded
// form (Resolve pads numeric ids before lookup, so "301" and "0301" both match).
var builtin = map[string]string{
	// City clusters — the player's actual standing location in town (from Join).
	"0000": "Thetford",
	"1000": "Lymhurst",
	"2000": "Bridgewatch",
	"3003": "Caerleon",
	"3004": "Martlock",
	"4000": "Fort Sterling",
	"5000": "Brecilien",
	"5001": "Brecilien",

	// Market instances — the auction-house zones that carry city price data.
	"0007": "Thetford Market",
	"1002": "Lymhurst Market",
	"2004": "Bridgewatch Market",
	"3005": "Caerleon Market",
	"3008": "Martlock Market",
	"4002": "Fort Sterling Market",
	"5003": "Brecilien Market",
	"3013": "Caerleon Market", // second Caerleon auction-house instance
	"3013-Auction2": "Caerleon Market",

	// Banks.
	"0006": "Thetford Bank",
	"1001": "Lymhurst Bank",
	"2003": "Bridgewatch Bank",
	"3006": "Caerleon Bank",
	"3007": "Martlock Bank",
	"4001": "Fort Sterling Bank",
	"5002": "Brecilien Bank",

	// Portal towns.
	"0301": "Thetford Portal",
	"1301": "Lymhurst Portal",
	"2301": "Bridgewatch Portal",
	"3301": "Martlock Portal",
	"4301": "Fort Sterling Portal",

	// Brecilien districts.
	"5100": "Brecilien Weald",
	"5101": "Brecilien Copse",
	"5102": "Brecilien Grove",
	"5103": "Brecilien Thicket",
}

// Resolve returns a display name for a raw location id. It matches personal/guild
// islands by template prefix, then tries the raw id, the segment before an '@',
// and the zero-padded numeric form. Returns the raw value (and false) when unknown
// so the UI still shows something.
func Resolve(raw string) (string, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", false
	}
	// Personal/guild islands arrive as templates like
	// "ISLAND-PLAYER-CAERLEON-0001a" / "ISLAND-GUILD-CAERLEON-0001b" — collapse the
	// whole family to a single label regardless of city and instance suffix.
	switch up := strings.ToUpper(raw); {
	case strings.HasPrefix(up, "ISLAND-PLAYER"):
		return "Player Island", true
	case strings.HasPrefix(up, "ISLAND-GUILD"):
		return "Guild Island", true
	}
	// Exact match first (covers non-numeric ids like "3013-Auction2").
	if n, ok := builtin[raw]; ok {
		return n, true
	}
	// Strip an "@instance" suffix.
	base := raw
	if at := strings.IndexByte(base, '@'); at > 0 {
		base = base[:at]
		if n, ok := builtin[base]; ok {
			return n, true
		}
	}
	// Numeric ids: pad to the 4-digit map key form ("301" -> "0301").
	if isNumeric(base) {
		if n, ok := builtin[pad4(base)]; ok {
			return n, true
		}
	}
	return raw, false
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// pad4 left-pads a numeric id with zeros to 4 digits (longer ids are returned
// unchanged).
func pad4(s string) string {
	if len(s) >= 4 {
		return s
	}
	return strings.Repeat("0", 4-len(s)) + s
}
