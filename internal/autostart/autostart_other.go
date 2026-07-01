//go:build !windows

package autostart

// Set is a no-op off Windows.
func Set(enabled bool, exePath string) error { return nil }
