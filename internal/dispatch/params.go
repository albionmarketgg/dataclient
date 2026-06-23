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

// Int64Slice extracts a slice of int64 from an integer-array param.
func Int64Slice(v any) []int64 {
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

// Uint64Slice extracts a slice of uint64 from an integer-array param.
func Uint64Slice(v any) []uint64 {
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
