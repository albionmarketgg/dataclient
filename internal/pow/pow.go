// Package pow implements the proof-of-work upload handshake. Hash input is
// "aod^" + 16-hex counter + "^" + key, fed to SHA-256; difficulty is N leading
// zero bits, where N == len(Request.Wanted).
package pow

import (
	"crypto/sha256"
)

// Request is the challenge issued by the ingest server.
type Request struct {
	Key    string `json:"key"`
	Wanted string `json:"wanted"` // length = required leading zero bits
}

const hexDigits = "0123456789abcdef"

// writeCounterHex writes value as 16 lowercase hex chars, most-significant first.
func writeCounterHex(dst []byte, value uint64) {
	for i := 15; i >= 0; i -= 2 {
		b := byte(value & 0xff)
		value >>= 8
		dst[i-1] = hexDigits[b>>4]
		dst[i] = hexDigits[b&0xf]
	}
}

// incrementHexAsciiInPlace increments a lowercase-hex ASCII counter in place.
func incrementHexAsciiInPlace(s []byte) {
	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		switch {
		case c == 'f':
			s[i] = '0' // carry
		case c == '9':
			s[i] = 'a'
			return
		default:
			s[i] = c + 1
			return
		}
	}
}

// leadingZeroBits reports whether hash starts with at least n zero bits.
func leadingZeroBits(hash []byte, n int) bool {
	full := n / 8
	for i := 0; i < full; i++ {
		if hash[i] != 0 {
			return false
		}
	}
	if rem := n % 8; rem != 0 {
		mask := byte(0xff << (8 - rem))
		if hash[full]&mask != 0 {
			return false
		}
	}
	return true
}

// Verify reports whether solution (16 hex chars) satisfies the challenge.
func Verify(req Request, solution string) bool {
	if len(solution) != 16 {
		return false
	}
	input := make([]byte, 0, 4+16+1+len(req.Key))
	input = append(input, "aod^"...)
	input = append(input, solution...)
	input = append(input, '^')
	input = append(input, req.Key...)
	sum := sha256.Sum256(input)
	return leadingZeroBits(sum[:], len(req.Wanted))
}

// Solve finds a counter solving the challenge and returns it as 16 hex chars.
func Solve(req Request) string {
	difficulty := len(req.Wanted)
	prefix := []byte("aod^")
	suffix := append([]byte{'^'}, req.Key...)
	input := make([]byte, 0, len(prefix)+16+len(suffix))
	input = append(input, prefix...)
	counterStart := len(input)
	input = append(input, make([]byte, 16)...)
	input = append(input, suffix...)
	counter := input[counterStart : counterStart+16]
	writeCounterHex(counter, 0)

	for {
		sum := sha256.Sum256(input)
		if leadingZeroBits(sum[:], difficulty) {
			return string(counter)
		}
		incrementHexAsciiInPlace(counter)
	}
}
