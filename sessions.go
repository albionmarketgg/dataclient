package main

// Session kinds. "dungeon" scopes fame + damage + loot for a PvE run; "pvp" is a
// separate damage-only session for open-world / arena fights; "gathering" is its
// own thing. Loot no longer has its own session — it rides the dungeon session.
var sessionKinds = []string{"gathering", "dungeon", "pvp"}

// StartSession starts a capture session for a kind (gathering/dungeon/pvp); the
// manager opens a server session and tags uploads with its id while active.
func (a *App) StartSession(kind string) {
	// dungeon/pvp both scope the damage meter (+ fame on dungeon): begin fresh.
	if kind == "dungeon" || kind == "pvp" {
		a.eng.Combat.Reset()
	}
	a.sessions.Start(kind)
	a.eng.Log("Started " + kind + " session.")
}

// StopSession ends a capture session for a kind (clean handshake end).
func (a *App) StopSession(kind string) {
	a.sessions.End(kind)
	a.eng.Log("Stopped " + kind + " session.")
}

// GetSessions returns each kind's session start time as unix-ms (0 = inactive).
func (a *App) GetSessions() map[string]int64 {
	out := make(map[string]int64, len(sessionKinds))
	for _, k := range sessionKinds {
		if t, ok := a.sessions.StartedAt(k); ok {
			out[k] = t.UnixMilli()
		} else {
			out[k] = 0
		}
	}
	return out
}

func (a *App) sessionActive(kind string) bool { return a.sessions.Active(kind) }

func (a *App) anySessionActive() bool { return a.sessions.Any() }
