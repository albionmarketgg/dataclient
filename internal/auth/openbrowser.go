package auth

import (
	"os/exec"
	"runtime"
)

// openInBrowser opens url in the default browser (fallback when Wails'
// BrowserOpenURL isn't injected).
func openInBrowser(url string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return exec.Command("xdg-open", url).Start()
	}
}
