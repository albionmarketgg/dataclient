// Command albionmarket-dataclient is the Albion Market companion desktop app.
// It passively captures Albion Online market traffic and uploads it to your
// local Albion Market ingest service. No real/public backend is referenced.
package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// version is set at build time via -ldflags "-X main.version=…".
var version = "dev"

func main() {
	app := NewApp(configPath())

	err := wails.Run(&options.App{
		Title:            "Albion Market — Data Client",
		Width:            1200,
		Height:           800,
		MinWidth:         960,
		MinHeight:        620,
		AssetServer:      &assetserver.Options{Assets: assets},
		BackgroundColour: &options.RGBA{R: 22, G: 22, B: 22, A: 255},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		OnBeforeClose:    app.beforeClose,
		// Single-instance: a second launch shows the running window, then exits.
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId:               "albionmarket-dataclient-1a2b3c",
			OnSecondInstanceLaunch: app.onSecondInstance,
		},
		Bind: []interface{}{app},
	})
	if err != nil {
		log.Fatal(err)
	}
}

// configPath returns the per-user config file location.
func configPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = "."
	}
	d := filepath.Join(dir, "AlbionMarketDataClient")
	_ = os.MkdirAll(d, 0o755)
	return filepath.Join(d, "config.json")
}
