package usersync

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestSessionHandshake(t *testing.T) {
	var mu sync.Mutex
	bodies := map[string]string{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		mu.Lock()
		bodies[r.URL.Path] = string(b)
		mu.Unlock()
		if r.URL.Path == "/user/session/start" {
			w.Write([]byte(`{"ok":true,"sessionId":777}`))
			return
		}
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()
	body := func(p string) string { mu.Lock(); defer mu.Unlock(); return bodies[p] }

	token := func() (string, string, bool) { return "jwt-x", "u1", true }
	sm := NewSessionManager(srv.URL, token, func() int { return 2 }, func() string { return "Hero" }, nil)

	// start -> opens server session, records the assigned id
	sm.Start("dungeon")
	deadline := time.Now().Add(2 * time.Second)
	for sm.CurrentID("dungeon") == 0 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if sm.CurrentID("dungeon") != 777 {
		t.Fatalf("sessionId not recorded: got %d", sm.CurrentID("dungeon"))
	}
	if b := body("/user/session/start"); !contains(b, `"kind":"dungeon"`) || !contains(b, `"serverId":2`) || !contains(b, `"playerName":"Hero"`) {
		t.Fatalf("bad start body: %q", b)
	}
	if !sm.Active("dungeon") || !sm.Any() {
		t.Fatal("session should be active")
	}

	// a dungeon upload now carries the sessionId
	s := New(srv.URL, nil, token, func(string) bool { return true }, Snapshots{
		ServerID:   func() int { return 2 },
		PlayerName: func() string { return "Hero" },
	}, nil)
	s.SetSessions(sm)
	s.UploadDungeonEvent(365, 1200, time.Now())
	if b := body("/user/dungeon"); !contains(b, `"sessionId":777`) {
		t.Fatalf("upload not tagged with sessionId: %q", b)
	}

	// end -> closes server session and clears local state
	sm.End("dungeon")
	deadline = time.Now().Add(2 * time.Second)
	for body("/user/session/end") == "" && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if b := body("/user/session/end"); !contains(b, `"sessionId":777`) {
		t.Fatalf("end not sent for the session: %q", b)
	}
	if sm.Active("dungeon") || sm.CurrentID("dungeon") != 0 {
		t.Fatal("session should be cleared after end")
	}
}

// If the user stops before /start returns an id (or it returned none), End must
// still close the session server-side via end-all-by-kind — never silently drop.
func TestSessionStopBeforeStartIdEndsViaEndAll(t *testing.T) {
	var mu sync.Mutex
	bodies := map[string]string{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		mu.Lock()
		bodies[r.URL.Path] = string(b)
		mu.Unlock()
		// /start returns ok but NO sessionId (simulates slow/unknown id)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()
	body := func(p string) string { mu.Lock(); defer mu.Unlock(); return bodies[p] }

	token := func() (string, string, bool) { return "jwt-x", "u1", true }
	sm := NewSessionManager(srv.URL, token, func() int { return 1 }, func() string { return "Hero" }, nil)

	sm.Start("gathering")
	sm.End("gathering")

	deadline := time.Now().Add(2 * time.Second)
	for body("/user/session/end-all") == "" && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if b := body("/user/session/end-all"); !contains(b, `"kind":"gathering"`) {
		t.Fatalf("stop without an id did not fall back to end-all-by-kind: %q", b)
	}
	if sm.Active("gathering") {
		t.Fatal("session should be inactive after end")
	}
}

func TestSessionStartLoggedOutNoServerCall(t *testing.T) {
	var mu sync.Mutex
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		hits++
		mu.Unlock()
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	loggedOut := func() (string, string, bool) { return "", "", false }
	sm := NewSessionManager(srv.URL, loggedOut, func() int { return 1 }, func() string { return "" }, nil)
	sm.Start("loot")
	time.Sleep(150 * time.Millisecond)

	// active locally for the UI, but no server session id and no network call
	if !sm.Active("loot") {
		t.Fatal("should be locally active even logged out")
	}
	if sm.CurrentID("loot") != 0 {
		t.Fatal("should have no server id when logged out")
	}
	mu.Lock()
	defer mu.Unlock()
	if hits != 0 {
		t.Fatalf("made %d server calls while logged out", hits)
	}
}
