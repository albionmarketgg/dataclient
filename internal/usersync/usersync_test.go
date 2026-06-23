package usersync

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/niick1231/albionmarket_dataclient/internal/store"
)

type capture struct {
	mu     sync.Mutex
	bodies map[string]string
	auth   map[string]string
}

func newCapture() *capture { return &capture{bodies: map[string]string{}, auth: map[string]string{}} }

func (c *capture) server() *httptest.Server {
	h := func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		c.mu.Lock()
		c.bodies[r.URL.Path] = string(b)
		c.auth[r.URL.Path] = r.Header.Get("Authorization")
		c.mu.Unlock()
		w.Write([]byte(`{"ok":true}`))
	}
	mux := http.NewServeMux()
	for _, p := range []string{"/user/trades", "/user/mails", "/user/dungeon", "/user/loot", "/user/party", "/user/gathering"} {
		mux.HandleFunc(p, h)
	}
	return httptest.NewServer(mux)
}

func (c *capture) body(p string) string { c.mu.Lock(); defer c.mu.Unlock(); return c.bodies[p] }

func newStore(t *testing.T) *store.Store {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestSyncUploadsAndMarksSynced(t *testing.T) {
	cap := newCapture()
	srv := cap.server()
	defer srv.Close()
	db := newStore(t)
	defer db.Close()

	// seed a trade + a mail
	db.InsertTrade(store.Trade{ID: "t-guid-1", AlbionServerID: 1, Amount: 3, DateTime: time.Now(), ItemID: "T4_BAG", Operation: 0, Type: 0, UnitSilver: 1234, PlayerName: "Hero", RawLocationID: "3005"})
	db.InsertMailInfo(store.Mail{ID: 42, AlbionServerID: 1, AuctionType: 1, RawLocationID: "3005", PlayerName: "Hero", Received: time.Now(), Type: 1})
	db.SetMailData(42, 5, 5, "T5_BAG", 154984, 0, 30996.8) // marks IsSet=1

	token := func() (string, string, bool) { return "jwt-abc", "u1", true }
	enabled := func(string) bool { return true }
	snaps := Snapshots{
		Loot:      func() any { return map[string]any{"records": []any{}} },
		Party:     func() any { return map[string]any{"members": []any{}} },
		Gathering: func() any { return map[string]any{"active": false} },
	}
	s := New(srv.URL, db, token, enabled, snaps, nil)
	s.syncOnce(context.Background())

	if b := cap.body("/user/trades"); b == "" || !contains(b, "T4_BAG") || !contains(b, "t-guid-1") {
		t.Fatalf("trades not uploaded: %q", b)
	}
	if b := cap.body("/user/mails"); b == "" || !contains(b, "\"id\":42") {
		t.Fatalf("mails not uploaded: %q", b)
	}
	if cap.auth["/user/trades"] != "Bearer jwt-abc" {
		t.Fatalf("missing bearer: %q", cap.auth["/user/trades"])
	}
	// trades/mails should now be marked synced
	if tr, _ := db.UnsyncedTrades(10); len(tr) != 0 {
		t.Fatalf("trades not marked synced: %d", len(tr))
	}
	if ml, _ := db.UnsyncedMails(10); len(ml) != 0 {
		t.Fatalf("mails not marked synced: %d", len(ml))
	}
	// snapshot trackers uploaded
	if cap.body("/user/loot") == "" {
		t.Fatal("loot snapshot not uploaded")
	}

	// second pass: nothing new -> trades/mails not re-sent (no unsynced), and the
	// unchanged loot snapshot is skipped (signature dedup)
	cap.mu.Lock()
	cap.bodies = map[string]string{}
	cap.mu.Unlock()
	s.syncOnce(context.Background())
	if cap.body("/user/loot") != "" {
		t.Fatal("unchanged loot snapshot re-sent (should dedup)")
	}
}

func TestUploadDungeonEvent(t *testing.T) {
	cap := newCapture()
	srv := cap.server()
	defer srv.Close()

	token := func() (string, string, bool) { return "jwt-abc", "u1", true }
	snaps := Snapshots{
		ServerID:   func() int { return 1 },
		PlayerName: func() string { return "Hero" },
	}
	s := New(srv.URL, nil, token, func(string) bool { return true }, snaps, nil)

	// fame paired with silver -> uploaded with both, plus serverId + playerName
	s.UploadDungeonEvent(365, 1200, time.Now())
	b := cap.body("/user/dungeon")
	if b == "" || !contains(b, "\"fame\":365") || !contains(b, "\"silver\":1200") {
		t.Fatalf("dungeon event not uploaded with fame+silver: %q", b)
	}
	if !contains(b, "\"playerName\":\"Hero\"") || !contains(b, "\"serverId\":1") {
		t.Fatalf("dungeon event missing attribution: %q", b)
	}
	if s.Stats().Dungeon != 1 {
		t.Fatalf("dungeon counter not incremented: %d", s.Stats().Dungeon)
	}

	// silver-only (fame<=0) must never be sent
	cap.mu.Lock()
	cap.bodies = map[string]string{}
	cap.mu.Unlock()
	s.UploadDungeonEvent(0, 5000, time.Now())
	if cap.body("/user/dungeon") != "" {
		t.Fatal("silver-only dungeon event was uploaded (should be dropped)")
	}

	// consent off -> nothing
	off := New(srv.URL, nil, token, func(string) bool { return false }, snaps, nil)
	off.UploadDungeonEvent(100, 0, time.Now())
	if cap.body("/user/dungeon") != "" {
		t.Fatal("uploaded dungeon event without consent")
	}
}

func TestSyncGatedByLoginAndConsent(t *testing.T) {
	cap := newCapture()
	srv := cap.server()
	defer srv.Close()
	db := newStore(t)
	defer db.Close()
	db.InsertTrade(store.Trade{ID: "t1", AlbionServerID: 1, ItemID: "T4_BAG", DateTime: time.Now()})

	snaps := Snapshots{Gathering: func() any { return map[string]any{} }}

	// logged out
	loggedOut := New(srv.URL, db, func() (string, string, bool) { return "", "", false }, func(string) bool { return true }, snaps, nil)
	loggedOut.syncOnce(context.Background())
	if cap.body("/user/trades") != "" {
		t.Fatal("uploaded while logged out")
	}

	// logged in but consent off
	noConsent := New(srv.URL, db, func() (string, string, bool) { return "jwt", "u1", true }, func(string) bool { return false }, snaps, nil)
	noConsent.syncOnce(context.Background())
	if cap.body("/user/trades") != "" {
		t.Fatal("uploaded without consent")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || indexOf(s, sub) >= 0)
}
func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

var _ = json.Marshal
