//go:build windows

// Package autostart manages launching the client at Windows login via the
// per-user Run key (no admin needed).
package autostart

import "golang.org/x/sys/windows/registry"

const (
	runKey    = `Software\Microsoft\Windows\CurrentVersion\Run`
	valueName = "AlbionMarketDataClient"
)

// Set adds (enabled) or removes (disabled) the Run entry pointing at exePath.
func Set(enabled bool, exePath string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	if enabled {
		return k.SetStringValue(valueName, `"`+exePath+`"`)
	}
	_ = k.DeleteValue(valueName) // best-effort; ignore "not present"
	return nil
}
