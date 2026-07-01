package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/niick1231/albionmarket_dataclient/internal/auth"
	"github.com/niick1231/albionmarket_dataclient/internal/autostart"
	"github.com/niick1231/albionmarket_dataclient/internal/config"
	"github.com/niick1231/albionmarket_dataclient/internal/engine"
	"github.com/niick1231/albionmarket_dataclient/internal/handlers"
	"github.com/niick1231/albionmarket_dataclient/internal/items"
	"github.com/niick1231/albionmarket_dataclient/internal/state"
	"github.com/niick1231/albionmarket_dataclient/internal/store"
	"github.com/niick1231/albionmarket_dataclient/internal/trackers"
	"github.com/niick1231/albionmarket_dataclient/internal/update"
	"github.com/niick1231/albionmarket_dataclient/internal/usersync"
)

// App is the Wails-bound application backend.
type App struct {
	ctx      context.Context
	eng      *engine.Engine
	items    *items.Service
	cfg      config.Config
	cfgPath  string
	quitting bool
	auth     *auth.Discord
	usersync *usersync.Syncer
	sessions *usersync.SessionManager

	updater    *update.Checker
	updateMu   sync.Mutex
	lastUpdate update.Result

	lastServerID int // tracks region changes to re-home active sessions
}

// NewApp constructs the app, loading config and wiring the engine.
func NewApp(cfgPath string) *App {
	cfg, _ := config.Load(cfgPath)
	dataDir := filepath.Join(filepath.Dir(cfgPath), "data")
	itemSvc := items.New(cfg.EffectiveItemsURL(), filepath.Join(dataDir, "items.txt"))
	dbPath := filepath.Join(dataDir, "client.db")
	a := &App{
		eng:     engine.New(cfg, itemSvc, dbPath),
		items:   itemSvc,
		cfg:     cfg,
		cfgPath: cfgPath,
	}
	a.initAuth()
	a.initUserSync()
	a.updater = update.New(cfg.AuthBaseURL, version)
	return a
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.eng.OnState(func(s state.Snapshot) { runtime.EventsEmit(ctx, "state", s) })
	a.eng.OnFeed(func(ev handlers.CaptureEvent) { runtime.EventsEmit(ctx, "feed", ev) })
	a.eng.OnLog(func(l engine.LogLine) { runtime.EventsEmit(ctx, "log", l) })
	a.eng.OnCaptureError(func(msg string) {
		// missing Npcap runtime → prompt the user to install it.
		if strings.Contains(strings.ToLower(msg), "wpcap") {
			runtime.EventsEmit(ctx, "npcapMissing", msg)
		}
	})
	a.eng.Log("Albion Market Data Client started.")
	a.applyAutostart()
	go a.auth.Restore(ctx)
	go a.startTray()
	if a.cfg.StartInTray {
		runtime.WindowHide(ctx)
		a.eng.Log("Started minimized to the system tray.")
	}
	go a.items.Load()
	a.eng.Start()
	a.usersync.Start()
	a.sessions.Run()
	go a.updateLoop(ctx)

	// periodic stats push
	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				snap := a.eng.State.Snapshot()
				// region change: re-home active sessions onto the new server.
				if snap.ServerID != 0 && a.lastServerID != 0 && snap.ServerID != a.lastServerID {
					a.sessions.OnRegionChange()
				}
				if snap.ServerID != 0 {
					a.lastServerID = snap.ServerID
				}
				runtime.EventsEmit(ctx, "stats", a.eng.Up.Stats())
				runtime.EventsEmit(ctx, "state", snap)
				runtime.EventsEmit(ctx, "userstats", a.GetUserStats())
			}
		}
	}()
}

// updateLoop checks for a newer build shortly after startup, then daily.
func (a *App) updateLoop(ctx context.Context) {
	t := time.NewTimer(5 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
		a.runUpdateCheck()
		t.Reset(24 * time.Hour)
	}
}

func (a *App) runUpdateCheck() {
	res, err := a.updater.Check(context.Background())
	if err != nil {
		return
	}
	a.updateMu.Lock()
	a.lastUpdate = res
	a.updateMu.Unlock()
	if res.ShouldNotify() && a.ctx != nil {
		runtime.EventsEmit(a.ctx, "update", res)
	}
}

// OpenNpcapDownload opens the official Npcap download page in the browser.
func (a *App) OpenNpcapDownload() {
	if a.ctx != nil {
		runtime.BrowserOpenURL(a.ctx, "https://npcap.com/#download")
	}
}

// RetryCapture re-attempts packet capture (after the user installs Npcap).
func (a *App) RetryCapture() { a.eng.StartCapture() }

// GetUpdateInfo returns the latest version-check result (empty until the first
// check completes).
func (a *App) GetUpdateInfo() update.Result {
	a.updateMu.Lock()
	defer a.updateMu.Unlock()
	return a.lastUpdate
}

// OpenUpdateDownload opens the update download page in the browser.
func (a *App) OpenUpdateDownload() {
	a.updateMu.Lock()
	dl := a.lastUpdate.DownloadURL
	a.updateMu.Unlock()
	switch {
	case dl == "":
		dl = a.cfg.AuthBaseURL
	case strings.HasPrefix(dl, "http"):
		// absolute already
	default:
		dl = strings.TrimRight(a.cfg.AuthBaseURL, "/") + "/" + strings.TrimLeft(dl, "/")
	}
	if a.ctx != nil {
		runtime.BrowserOpenURL(a.ctx, dl)
	}
}

func (a *App) shutdown(ctx context.Context) {
	a.quitTray()
	// close all active sessions cleanly before we go (best-effort handshake).
	if a.sessions != nil {
		a.sessions.EndAll("")
		a.sessions.Stop()
	}
	if a.usersync != nil {
		a.usersync.Stop()
	}
	a.eng.Stop()
}

// ---- bound methods (callable from the frontend) ----

// GetSnapshot returns the current player/session state.
func (a *App) GetSnapshot() state.Snapshot { return a.eng.State.Snapshot() }

// GetStats returns upload counters.
func (a *App) GetStats() any { return a.eng.Up.Stats() }

// GetFeed returns recent capture events.
func (a *App) GetFeed() []handlers.CaptureEvent { return a.eng.Feed() }

// GetLogs returns recent log lines.
func (a *App) GetLogs() []engine.LogLine { return a.eng.Logs() }

// GetDevices lists available capture devices.
func (a *App) GetDevices() []string {
	devs, err := captureDevices()
	if err != nil {
		return []string{"(error: " + err.Error() + ")"}
	}
	return devs
}

// IsCapturing reports capture status.
func (a *App) IsCapturing() bool { return a.eng.Capturing() }

// ToggleCapture starts or stops the capture listener.
func (a *App) ToggleCapture(on bool) {
	if on {
		a.eng.StartCapture()
	} else {
		a.eng.StopCapture()
	}
}

// GetAwakenedItems returns the current awakened-inventory snapshot.
func (a *App) GetAwakenedItems() trackers.AwakenedSnapshot { return a.eng.Awakened.Snapshot() }

// AwakenedSellItem is one item to list, with its own price (0 == offers only).
type AwakenedSellItem struct {
	Key   string `json:"key"`
	Price int64  `json:"price"`
}

// AwakenedSellRequest lists one or more awakened items, each at its own price.
type AwakenedSellRequest struct {
	Items []AwakenedSellItem `json:"items"`
}

// RemoveAwakened hides an item from the Awakened Inventory table.
func (a *App) RemoveAwakened(key string) { a.eng.Awakened.Remove(key) }

// SellOnAlbionMarket lists the given awakened items for sale on albionmarket.gg,
// each at its own price (0 == offers only). Works for a single item or a bulk
// list. Login-gated. Returns "" on success, else an error message.
//
// TODO: confirm the endpoint + payload with the web agent (see the sell API outline).
func (a *App) SellOnAlbionMarket(req AwakenedSellRequest) string {
	if len(req.Items) == 0 {
		return "no items selected"
	}
	jwt, _, ok := a.tokenProvider()
	if !ok {
		return "please sign in first"
	}
	byKey := map[string]trackers.AwakenedItem{}
	for _, it := range a.eng.Awakened.Snapshot().Items {
		byKey[it.Key] = it
	}
	// Per-item payload: itemId + serverId + quality + strain + traits + price.
	// Seller identity comes from the JWT + serverId (backend checks the user has a
	// verified character there); attuned-to name and display strings are dropped.
	items := make([]map[string]any, 0, len(req.Items))
	for _, sel := range req.Items {
		it, found := byKey[sel.Key]
		if !found {
			continue
		}
		traits := make([]map[string]any, 0, len(it.Traits))
		for _, tr := range it.Traits {
			traits = append(traits, map[string]any{"id": tr.ID, "value": tr.Value})
		}
		items = append(items, map[string]any{
			"itemId":   it.ItemID,
			"serverId": it.ServerID,
			"quality":  it.Quality,
			"strain":   it.Strain,
			"traits":   traits,
			"price":    sel.Price,
		})
	}
	body := map[string]any{"items": items}
	if err := a.postJSON(jwt, "/user/awakened/list", body); err != nil {
		return err.Error()
	}
	a.eng.Log(fmt.Sprintf("Listed %d awakened item(s) for sale.", len(items)))
	return ""
}

// postJSON is a small authenticated JSON POST helper for bound methods.
func (a *App) postJSON(jwt, path string, body any) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(a.cfg.IngestBaseURL, "/")+path, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("User-Agent", "AlbionMarketDataClient")
	resp, err := (&http.Client{Timeout: 20 * time.Second}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("server returned %d", resp.StatusCode)
	}
	return nil
}

// GetConfig returns the current configuration.
func (a *App) GetConfig() config.Config { return a.cfg }

// SaveConfig persists a new configuration (takes effect on restart for capture).
func (a *App) SaveConfig(cfg config.Config) string {
	a.cfg = cfg
	if err := cfg.Save(a.cfgPath); err != nil {
		return err.Error()
	}
	a.applyAutostart()
	a.eng.Log("Configuration saved. Restart to apply capture changes.")
	return ""
}

// applyAutostart syncs the Windows "start at login" registry entry to the
// current StartWithWindows setting (points at this executable).
func (a *App) applyAutostart() {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	_ = autostart.Set(a.cfg.StartWithWindows, exe)
}

// ItemsLoaded reports whether the item-name database has loaded.
func (a *App) ItemsLoaded() bool { return a.items.Loaded() }

// GetMails returns recent persisted marketplace mails.
func (a *App) GetMails() []store.Mail { return a.eng.Mails(1000) }

// GetTrades returns recent persisted trades.
func (a *App) GetTrades() []store.Trade { return a.eng.Trades(1000) }

// GetParty returns the current party snapshot.
func (a *App) GetParty() trackers.PartySnapshot { return a.eng.Party.Snapshot() }

// GetGathering returns the gathering session snapshot.
func (a *App) GetGathering() trackers.GatherSnapshot { return a.eng.Gathering.Snapshot() }

// GetCombat returns the combat snapshot.
func (a *App) GetCombat() trackers.CombatSnapshot { return a.eng.Combat.Snapshot() }

// GetLoot returns the loot snapshot.
func (a *App) GetLoot() trackers.LootSnapshot { return a.eng.Loot.Snapshot() }

// GetUserStats returns per-kind personal-data upload counts for this session.
func (a *App) GetUserStats() usersync.Stats {
	if a.usersync == nil {
		return usersync.Stats{}
	}
	return a.usersync.Stats()
}

// GetDetectedCharacters returns characters detected from Join packets.
func (a *App) GetDetectedCharacters() []state.Character { return a.eng.State.Characters() }

// SubmitVerification verifies a detected character against the account.
func (a *App) SubmitVerification(serverID int, name string, fame int64) auth.VerifyResult {
	if a.auth == nil || !a.auth.Configured() {
		return auth.VerifyResult{Status: "error", Message: "Sign-in is not available."}
	}
	res, err := a.auth.Verify(context.Background(), serverID, name, fame)
	if err != nil {
		return auth.VerifyResult{Status: "error", Message: err.Error()}
	}
	return res
}

// ResetGathering clears the gathering session.
func (a *App) ResetGathering() { a.eng.Gathering.Reset() }

// ResetCombat clears combat encounters.
func (a *App) ResetCombat() { a.eng.Combat.Reset() }

// ClearLoot clears the loot log.
func (a *App) ClearLoot() { a.eng.Loot.Clear() }
