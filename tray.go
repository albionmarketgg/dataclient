package main

import (
	"context"
	_ "embed"

	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed build/icon.ico
var trayIcon []byte

// startTray installs the system-tray icon and menu; the app keeps capturing
// while hidden to the tray.
func (a *App) startTray() {
	systray.Run(func() {
		systray.SetIcon(trayIcon)
		systray.SetTitle("Albion Market")
		systray.SetTooltip("Albion Market — Data Client")

		mShow := systray.AddMenuItem("Show Albion Market", "Open the window")
		mCapture := systray.AddMenuItem("Toggle Capture", "Start/stop capture")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Exit the data client")

		mShow.Click(func() { a.showWindow() })
		mCapture.Click(func() { a.ToggleCapture(!a.eng.Capturing()) })
		mQuit.Click(func() { a.quit() })

		// Left-click the tray icon to restore; right-click opens the menu.
		systray.SetOnClick(func(menu systray.IMenu) { a.showWindow() })
		systray.SetOnRClick(func(menu systray.IMenu) { menu.ShowMenu() })
	}, func() {})
}

func (a *App) showWindow() {
	if a.ctx == nil {
		return
	}
	runtime.WindowShow(a.ctx)
	runtime.WindowUnminimise(a.ctx)
}

// onSecondInstance runs in the first instance when a second copy is launched:
// bring our window to the front instead of starting a duplicate.
func (a *App) onSecondInstance(_ options.SecondInstanceData) {
	a.showWindow()
	if a.ctx != nil {
		runtime.WindowSetAlwaysOnTop(a.ctx, true)
		runtime.WindowSetAlwaysOnTop(a.ctx, false)
	}
}

// quitTray removes the system-tray icon cleanly (avoids a lingering ghost icon).
func (a *App) quitTray() { systray.Quit() }

func (a *App) quit() {
	a.quitting = true
	if a.ctx != nil {
		runtime.Quit(a.ctx)
	}
}

// HideToTray hides the window to the system tray (bound, callable from the UI).
func (a *App) HideToTray() {
	if a.ctx != nil {
		runtime.WindowHide(a.ctx)
	}
}

// QuitApp fully exits the application (bound).
func (a *App) QuitApp() { a.quit() }

// beforeClose hides to tray instead of quitting, unless the user chose to quit
// or CloseToTray is disabled.
func (a *App) beforeClose(_ context.Context) bool {
	if a.quitting || !a.cfg.CloseToTray {
		return false // allow close
	}
	a.HideToTray()
	return true // prevent close; stay in tray
}
