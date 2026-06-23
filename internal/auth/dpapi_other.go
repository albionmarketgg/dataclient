//go:build !windows

package auth

import "encoding/base64"

// protect/unprotect are no-op (base64 only) on non-Windows; keeps the package
// buildable off Windows.
func protect(plaintext string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(plaintext)), nil
}

func unprotect(b64 string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(b64)
	return string(b), err
}
