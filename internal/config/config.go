// Package config holds runtime configuration. Point IngestBaseURL at your own
// Albion Market ingest service.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Server identifies an Albion game server (region) by its source-IP prefixes.
type Server struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	HostIPs  []string `json:"hostIps"`
}

// DefaultServers are the public Albion game-server IP prefixes used to label
// captured data by region.
func DefaultServers() []Server {
	return []Server{
		{ID: 1, Name: "Americas", HostIPs: []string{"5.188.125", "85.234.70"}},
		{ID: 2, Name: "Asia", HostIPs: []string{"5.45.187"}},
		{ID: 3, Name: "Europe", HostIPs: []string{"193.169.238"}},
	}
}

// Config is the application configuration.
type Config struct {
	// IngestBaseURL is our market-data ingest endpoint.
	IngestBaseURL string `json:"ingestBaseUrl"`
	// AuthBaseURL is our auth origin that brokers "Login with Discord" (the client
	// holds no Discord secret). Empty = auth disabled (anonymous uploads, login hidden).
	AuthBaseURL string `json:"authBaseUrl"`
	// RequirePoW toggles the proof-of-work handshake before uploads.
	RequirePoW bool `json:"requirePow"`
	// PacketFilter is the libpcap BPF filter for Albion's Photon UDP traffic.
	PacketFilter string `json:"packetFilter"`

	// StartInTray hides the window to the system tray on launch.
	StartInTray bool `json:"startInTray"`
	// CloseToTray hides to tray on window close instead of quitting.
	CloseToTray bool `json:"closeToTray"`
	// Per-type opt-in for syncing the user's own captured data to their account.
	// Each uploads only when logged in. Default ON.
	UploadTrades    bool `json:"uploadTrades"`
	UploadMails     bool `json:"uploadMails"`
	UploadGathering bool `json:"uploadGathering"`
	UploadCombat    bool `json:"uploadCombat"`
	UploadLoot      bool `json:"uploadLoot"`
	UploadParty     bool `json:"uploadParty"`
	// UploadSpecs syncs the character's Destiny Board (mastery/specialization
	// levels). Login-gated.
	UploadSpecs bool `json:"uploadSpecs"`

	// ItemsURL is the source for item id->name data. Empty means "<IngestBaseURL>/items.txt".
	ItemsURL string `json:"itemsUrl"`
	// AchievementsURL is the source for Destiny Board index->id data (achievements.xml).
	// Empty means "<IngestBaseURL>/achievements.xml".
	AchievementsURL string `json:"achievementsUrl"`
	// CaptureDevice optionally restricts capture to one device (by description
	// substring). Empty captures on all devices.
	CaptureDevice string `json:"captureDevice"`

	// Ingest topic/path segments.
	MarketOrdersTopic    string `json:"marketOrdersTopic"`
	MarketHistoriesTopic string `json:"marketHistoriesTopic"`
	GoldPricesTopic      string `json:"goldPricesTopic"`

	// Capture/idle tuning.
	NetworkStartDelaySecs int `json:"networkStartDelaySecs"`
	IdleMinutes           int `json:"idleMinutes"`
	IdleCheckMinutes      int `json:"idleCheckMinutes"`

	Servers []Server `json:"servers"`
}

// Default returns configuration pointing at the Albion Market services.
func Default() Config {
	return Config{
		IngestBaseURL:         "https://albionmarket.gg",
		AuthBaseURL:           "https://albionmarket.gg",
		RequirePoW:            true,
		PacketFilter:          "udp port 5056",
		StartInTray:           false,
		CloseToTray:           true,
		UploadTrades:          true,
		UploadMails:           true,
		UploadGathering:       true,
		UploadCombat:          true,
		UploadLoot:            true,
		UploadParty:           true,
		UploadSpecs:           true,
		MarketOrdersTopic:     "marketorders.ingest",
		MarketHistoriesTopic:  "markethistories.ingest",
		GoldPricesTopic:       "goldprices.ingest",
		NetworkStartDelaySecs: 3,
		IdleMinutes:           20,
		IdleCheckMinutes:      5,
		Servers:               DefaultServers(),
	}
}

// Load reads config from path, falling back to defaults for missing fields.
func Load(path string) (Config, error) {
	cfg := Default()
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return cfg, err
	}
	if len(cfg.Servers) == 0 {
		cfg.Servers = DefaultServers()
	}
	// Managed endpoints (not user-editable): always force to in-code defaults so a
	// stale config value can never repoint the client.
	def := Default()
	cfg.IngestBaseURL = def.IngestBaseURL
	cfg.AuthBaseURL = def.AuthBaseURL
	cfg.ItemsURL = def.ItemsURL
	cfg.PacketFilter = def.PacketFilter
	return cfg, nil
}

// EffectiveItemsURL returns the configured items URL, defaulting to
// "<IngestBaseURL>/items.txt" when ItemsURL is empty.
func (c Config) EffectiveItemsURL() string {
	if c.ItemsURL != "" {
		return c.ItemsURL
	}
	return c.IngestBaseURL + "/items.txt"
}

// EffectiveAchievementsURL returns the configured achievements.xml URL, defaulting
// to "<IngestBaseURL>/achievements.xml" (our mirror) when empty.
func (c Config) EffectiveAchievementsURL() string {
	if c.AchievementsURL != "" {
		return c.AchievementsURL
	}
	return c.IngestBaseURL + "/achievements.xml"
}

// Save writes config to path (pretty-printed).
func (c Config) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
