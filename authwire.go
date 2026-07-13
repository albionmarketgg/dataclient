package main

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/albionmarketgg/data-client/internal/auth"
	"github.com/albionmarketgg/data-client/internal/handlers"
	"github.com/albionmarketgg/data-client/internal/store"
	"github.com/albionmarketgg/data-client/internal/trackers"
	"github.com/albionmarketgg/data-client/internal/usersync"
)

// authStore adapts *store.Store to auth.RefreshStore.
type authStore struct{ db *store.Store }

func (a authStore) Save(sa auth.StoredAuth) error {
	if a.db == nil {
		return nil
	}
	return a.db.SaveUserAuth(store.UserAuth{
		UserID: sa.UserID, RefreshToken: sa.RefreshToken, Username: sa.Username, Avatar: sa.Avatar,
	})
}

func (a authStore) Load() (auth.StoredAuth, bool, error) {
	if a.db == nil {
		return auth.StoredAuth{}, false, nil
	}
	u, ok, err := a.db.LoadUserAuth()
	if err != nil || !ok {
		return auth.StoredAuth{}, ok, err
	}
	return auth.StoredAuth{UserID: u.UserID, RefreshToken: u.RefreshToken, Username: u.Username, Avatar: u.Avatar}, true, nil
}

func (a authStore) Clear() error {
	if a.db == nil {
		return nil
	}
	return a.db.ClearUserAuth()
}

// initAuth builds the Discord auth provider and wires it into uploads.
func (a *App) initAuth() {
	openURL := func(u string) error {
		if a.ctx != nil {
			runtime.BrowserOpenURL(a.ctx, u)
			return nil
		}
		return nil
	}
	a.auth = auth.New(a.cfg.AuthBaseURL, authStore{db: a.eng.Store()}, openURL, a.eng.Log)
	a.auth.SetOnPending(func(userCode, url string) {
		if a.ctx != nil {
			runtime.EventsEmit(a.ctx, "authPending", map[string]string{"userCode": userCode, "url": url})
		}
	})
	a.eng.Up.SetTokenProvider(a.tokenProvider)
}

// tokenProvider returns the current bearer token + user id (ok=false if logged out).
func (a *App) tokenProvider() (string, string, bool) {
	if a.auth == nil {
		return "", "", false
	}
	tok, ok := a.auth.Token(context.Background())
	if !ok {
		return "", "", false
	}
	u, _ := a.auth.User()
	return tok, u.ID, true
}

// initUserSync wires the login-gated, consent-gated personal-data uploader.
func (a *App) initUserSync() {
	a.usersync = usersync.New(
		a.cfg.IngestBaseURL,
		a.eng.Store(),
		a.tokenProvider,
		func(kind string) bool {
			switch kind {
			// ledger data: always-on (consent only)
			case "trades":
				return a.cfg.UploadTrades
			case "mails":
				return a.cfg.UploadMails
			// session data: consent AND the kind's session is running
			case "gathering":
				return a.cfg.UploadGathering && a.sessionActive("gathering")
			// dungeon (fame/silver events) rides the dungeon session
			case "dungeon":
				return a.cfg.UploadCombat && a.sessionActive("dungeon")
			// loot rides the dungeon session now (no separate loot session)
			case "loot":
				return a.cfg.UploadLoot && a.sessionActive("dungeon")
			// party rides alongside any active session (used to group combat/loot)
			case "party":
				return a.cfg.UploadParty && a.anySessionActive()
			// character specs: login + consent only (no session, no verification)
			case "specs":
				return a.cfg.UploadSpecs
			// awakened weapons: login + consent only
			case "awakened":
				return a.cfg.UploadAwakened
			// damage meter: rides EITHER the dungeon or the pvp session + its consent
			case "damage":
				return a.cfg.UploadCombat && (a.sessionActive("dungeon") || a.sessionActive("pvp"))
			}
			return false
		},
		usersync.Snapshots{
			Gathering: func() any { return a.eng.Gathering.Snapshot() },
			Loot:      func() any { return a.eng.Loot.Snapshot() },
			// party uploads as { roster: ["Name", ...] } per the backend spec
			Party: func() any {
				snap := a.eng.Party.Snapshot()
				roster := make([]string, 0, len(snap.Members))
				for _, m := range snap.Members {
					if m.Name != "" {
						roster = append(roster, m.Name)
					}
				}
				return map[string]any{"roster": roster}
			},
			Awakened: func() any { return a.eng.Awakened.SyncBody() },
			// damage upload = the party's cumulative group totals (players only),
			// pushed every few seconds so the site can graph the session.
			Combat:     func() any { return a.eng.Combat.SessionSummary() },
			ServerID:   func() int { return a.eng.State.Snapshot().ServerID },
			PlayerName: func() string { return a.eng.State.Snapshot().PlayerName },
		},
		a.eng.Log,
	)
	// session handshake: explicit start/end/heartbeat/end-all + sessionId tagging.
	a.sessions = usersync.NewSessionManager(
		a.cfg.IngestBaseURL,
		a.tokenProvider,
		func() int { return a.eng.State.Snapshot().ServerID },
		func() string { return a.eng.State.Snapshot().PlayerName },
		a.eng.Log,
	)
	a.usersync.SetSessions(a.sessions)
	// stamp uploads with the running client version (for the backend's version gate).
	a.usersync.SetClientVersion(version)
	a.sessions.SetClientVersion(version)
	// the awakened sync response echoes backend-computed trait values — apply them
	// onto the tracked items so the UI paints real numbers over its skeletons.
	a.usersync.SetOnAwakened(func(items []usersync.ResolvedAwakened) {
		for _, it := range items {
			traits := make([]trackers.ResolvedTrait, 0, len(it.Traits))
			for _, tr := range it.Traits {
				traits = append(traits, trackers.ResolvedTrait{ID: tr.ID, Value: tr.Value, Percent: tr.Percent})
			}
			a.eng.Awakened.ApplyValues(it.SoulID, traits)
		}
	})
	// dungeon data is event-driven: upload each fame (paired with its nearest
	// silver) the moment it resolves, not on the snapshot tick.
	a.eng.Combat.OnDungeonEvent(func(ev trackers.DungeonEvent) {
		go a.usersync.UploadDungeonEvent(ev.Fame, ev.Silver, ev.At)
	})
	// character specs: the game pushes the full Destiny Board on entering the
	// world; upload it (resolved to ids) for the current in-game character.
	// awakened inventory: event-driven — a change schedules a debounced (~2s) sync.
	a.eng.Awakened.OnChange(func() { a.usersync.TriggerAwakened() })
	a.eng.Specs.OnFull(func(version string, entries []handlers.SpecEntry) {
		snap := a.eng.State.Snapshot()
		out := make([]usersync.SpecEntry, len(entries))
		for i, e := range entries {
			out[i] = usersync.SpecEntry{ID: e.ID, Level: e.Level}
		}
		go a.usersync.UploadSpecs(snap.ServerID, snap.PlayerName, version, out)
	})
}

// ---- bound auth methods ----

// AuthEnabled reports whether Discord login is configured.
func (a *App) AuthEnabled() bool { return a.auth != nil && a.auth.Configured() }

// GetUser returns the logged-in user (empty id when logged out).
func (a *App) GetUser() auth.User {
	if a.auth == nil {
		return auth.User{}
	}
	u, _ := a.auth.User()
	return u
}

// Login starts the Discord device-login flow in the background; progress
// arrives via the "authPending" / "auth" / "authError" events.
func (a *App) Login() string {
	if a.auth == nil || !a.auth.Configured() {
		return "auth not configured"
	}
	go func() {
		if err := a.auth.Login(context.Background()); err != nil {
			a.eng.Log("Discord login failed: " + err.Error())
			if a.ctx != nil {
				runtime.EventsEmit(a.ctx, "authError", err.Error())
			}
			return
		}
		if a.ctx != nil {
			runtime.EventsEmit(a.ctx, "auth", a.GetUser())
		}
	}()
	return ""
}

// Logout signs out of Discord.
func (a *App) Logout() {
	if a.auth == nil {
		return
	}
	a.auth.Logout(context.Background())
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "auth", auth.User{})
	}
}
