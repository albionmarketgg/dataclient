package dispatch

// Parameter-extraction helpers. Photon values arrive as Go types produced by the
// Protocol18 deserializer (int16/int32/int64/byte/uint16/string/[]any/...).

// Strings extracts a []string from a param that is a []any of strings or []string.
func Strings(v any) []string {
	switch x := v.(type) {
	case []string:
		return x
	case []any:
		out := make([]string, 0, len(x))
		for _, e := range x {
			if s, ok := e.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

// StringMatrix flattens a [][]string (jagged) param (e.g. loadout offers).
func StringMatrix(v any) []string {
	arr, ok := v.([]any)
	if !ok {
		return Strings(v)
	}
	var out []string
	for _, row := range arr {
		out = append(out, Strings(row)...)
	}
	return out
}

// Int64 coerces any integer param to int64.
func Int64(v any) (int64, bool) {
	switch n := v.(type) {
	case int64:
		return n, true
	case int32:
		return int64(n), true
	case int16:
		return int64(n), true
	case byte:
		return int64(n), true
	case uint16:
		return int64(n), true
	case uint32:
		return int64(n), true
	case uint64:
		return int64(n), true
	case int:
		return int64(n), true
	}
	return 0, false
}

// Uint64 coerces any integer param to uint64.
func Uint64(v any) (uint64, bool) {
	if n, ok := Int64(v); ok {
		return uint64(n), true
	}
	return 0, false
}

// Int coerces to int.
func Int(v any) (int, bool) {
	if n, ok := Int64(v); ok {
		return int(n), true
	}
	return 0, false
}

// Int64Slice extracts a slice of int64 from an integer-array param. Photon
// encodes an integer array whose values all fit in a byte as a ByteArray, which
// deserializes to []byte — common for low-volume data (e.g. price-history sale
// counts with zero/small days), so both shapes must be accepted.
func Int64Slice(v any) []int64 {
	if b, ok := v.([]byte); ok {
		out := make([]int64, len(b))
		for i, e := range b {
			out[i] = int64(e)
		}
		return out
	}
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]int64, 0, len(arr))
	for _, e := range arr {
		if n, ok := Int64(e); ok {
			out = append(out, n)
		}
	}
	return out
}

// Uint64Slice extracts a slice of uint64 from an integer-array param (accepts
// the ByteArray shape like Int64Slice).
func Uint64Slice(v any) []uint64 {
	if b, ok := v.([]byte); ok {
		out := make([]uint64, len(b))
		for i, e := range b {
			out[i] = uint64(e)
		}
		return out
	}
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]uint64, 0, len(arr))
	for _, e := range arr {
		if n, ok := Uint64(e); ok {
			out = append(out, n)
		}
	}
	return out
}

// Uint32Slice extracts a slice of uint32 from an integer-array param.
func Uint32Slice(v any) []uint32 {
	s := Int64Slice(v)
	out := make([]uint32, len(s))
	for i, n := range s {
		out[i] = uint32(n)
	}
	return out
}

// String extracts a string param.
func String(v any) (string, bool) {
	s, ok := v.(string)
	return s, ok
}

// ItemQuality reads the quality (1-5) from a New*Item parameter table, tolerating
// both wire layouts: originally quality sat at [6]; a 2026-07 game patch inserted
// a second value field at [5] on equipment items, which shifted the crafter name
// (a string) to [6] and quality to [7]. The string at [6] discriminates the two.
// Unknown/absent quality defaults to 1 (the in-game minimum).
func ItemQuality(p map[byte]any) int {
	if q, ok := Int(p[6]); ok && q > 0 {
		return q
	}
	if _, isString := p[6].(string); isString {
		if q, ok := Int(p[7]); ok && q > 0 {
			return q
		}
	}
	return 1
}
