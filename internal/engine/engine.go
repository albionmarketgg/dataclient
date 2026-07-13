// Package engine wires the capture → parse → dispatch → upload pipeline and
// exposes a small surface for the UI layer.
package engine

import (
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/albionmarketgg/dataclient/internal/capture"
	"github.com/albionmarketgg/dataclient/internal/config"
	"github.com/albionmarketgg/dataclient/internal/dispatch"
	"github.com/albionmarketgg/dataclient/internal/handlers"
	"github.com/albionmarketgg/dataclient/internal/locations"
	"github.com/albionmarketgg/dataclient/internal/market"
	"github.com/albionmarketgg/dataclient/internal/mobs"
	"github.com/albionmarketgg/dataclient/internal/photon"
	"github.com/albionmarketgg/dataclient/internal/specs"
	"github.com/albionmarketgg/dataclient/internal/state"
	"github.com/albionmarketgg/dataclient/internal/store"
	"github.com/albionmarketgg/dataclient/internal/trackers"
	"github.com/albionmarketgg/dataclient/internal/upload"
)

// LogLine is a timestamped log entry for the UI.
type LogLine struct {
	Time time.Time `json:"time"`
	Text string    `json:"text"`
}

// Engine owns the full data pipeline.
type Engine struct {
	Cfg   config.Config
	State *state.State
	Up    *upload.Uploader

	disp      *dispatch.Dispatcher
	parser    *photon.Parser
	listener  *capture.Listener
	market    *handlers.Market
	mailTrade *handlers.MailTrade
	store     *store.Store

	Party     *trackers.Party
	Gathering *trackers.Gathering
	Combat    *trackers.Combat
	Loot      *trackers.Loot
	Awakened  *trackers.Awakened
	Specs     *handlers.Specs

	specsSvc *specs.Service
	mobsSvc  *mobs.Service

	feedMu sync.Mutex
	feed   []handlers.CaptureEvent

	logMu sync.Mutex
	logs  []LogLine

	packetsSeen atomic.Int64

	onFeed     func(handlers.CaptureEvent)
	onLog      func(LogLine)
	onState    func(state.Snapshot)
	onCaptureErr func(string)
}

// ItemDB resolves item names + values and stores EMVs seen on the wire.
// (*items.Service satisfies this; pass nil to disable item resolution.)
type ItemDB interface {
	UniqueName(id int) (string, bool)
	DisplayName(id int) (string, bool)
	EMV(id, quality int) (int64, bool)
	SetEMV(id, quality int, emv int64)
}

// New builds an Engine from config. itemDB may be nil. dbPath is the SQLite file
// for mail/trade persistence (pass "" to disable persistence).
func New(cfg config.Config, itemDB ItemDB, dbPath string) *Engine {
	e := &Engine{
		Cfg:   cfg,
		State: state.New(),
	}
	e.Up = upload.New(cfg, e.log)
	e.disp = dispatch.New()
	e.disp.OnAny(func() { e.packetsSeen.Add(1) })
	e.parser = photon.NewParser(e.disp)

	var namer handlers.ItemNamer
	var info trackers.ItemInfo
	if itemDB != nil {
		namer = itemDB
		info = itemDB
	}

	e.market = handlers.NewMarket(e.State, e.Up, namer, e.pushFeed)
	if itemDB != nil {
		e.market.SetEMVSink(itemDB.SetEMV)
	}
	e.market.Register(e.disp)

	// Join handler: establishes player name + location; required for the upload gate.
	handlers.NewSession(e.State, e.pushFeed).Register(e.disp)

	e.Party = trackers.NewParty()
	e.Gathering = trackers.NewGathering(info, e.State)
	e.Combat = trackers.NewCombat(e.Party,
		func(index int) (string, bool) {
			if e.mobsSvc != nil {
				return e.mobsSvc.Name(index)
			}
			return "", false
		},
		func() (int64, bool) { return e.State.UserObjectID() })
	e.Loot = trackers.NewLoot(info,
		func() string { return e.State.Snapshot().PlayerName },
		func() string { // the single OTHER party member's name (for their chest loot)
			others := ""
			n := 0
			for _, m := range e.Party.Snapshot().Members {
				if !m.IsLocal && m.Name != "" {
					others = m.Name
					n++
				}
			}
			if n == 1 {
				return others
			}
			return ""
		})
	e.Awakened = trackers.NewAwakened(info,
		func() int { return e.State.Snapshot().ServerID },
		func() string { return e.State.Snapshot().PlayerName },
		func() string {
			if name, ok := locations.Resolve(e.State.Snapshot().LocationID); ok {
				return name
			}
			return ""
		})
	e.Party.Register(e.disp)
	e.Gathering.Register(e.disp)
	e.Combat.Register(e.disp)
	e.Loot.Register(e.disp)
	e.Awakened.Register(e.disp)

	// Destiny Board reader: resolves the FullAchievementInfo packet to {id, level}
	// via the achievements.xml dump.
	cachePath := ""
	if dbPath != "" {
		cachePath = filepath.Join(filepath.Dir(dbPath), "achievements.xml")
	}
	e.specsSvc = specs.New(cfg.EffectiveAchievementsURL(), cachePath)
	e.Specs = handlers.NewSpecs(e.specsSvc)
	e.Specs.Register(e.disp)

	// Mob-name resolver for the damage meter (index -> name). Fetched + cached like
	// achievements.xml; nil-safe until loaded (mobs show as "Mob <index>").
	mobsCache := ""
	if dbPath != "" {
		mobsCache = filepath.Join(filepath.Dir(dbPath), "mobs.json")
	}
	e.mobsSvc = mobs.New(cfg.EffectiveMobsURL(), mobsCache)

	// surface gameplay-tracker activity in the live feed
	trackerFeed := func(kind, detail string, count int) {
		e.pushFeed(handlers.CaptureEvent{Time: time.Now(), Kind: kind, Detail: detail, Count: count})
	}
	e.Gathering.OnFeed(trackerFeed)
	e.Combat.OnFeed(trackerFeed)
	e.Loot.OnFeed(trackerFeed)
	e.Party.OnFeed(trackerFeed)

	if dbPath != "" {
		if db, err := store.Open(dbPath); err != nil {
			e.log("Persistence disabled (db open failed): " + err.Error())
		} else {
			e.store = db
			e.mailTrade = handlers.NewMailTrade(e.State, db, namer, e.pushFeed)
			e.mailTrade.Register(e.disp)
			e.market.OnOrders(func(orders []market.Order) { e.mailTrade.CacheOrders(orders) })
		}
	}

	e.listener = capture.New(cfg, e.State, e.parser, e.log)

	e.State.OnChange(func(s state.Snapshot) {
		e.Up.SetServerID(s.ServerID)
		if e.onState != nil {
			e.onState(s)
		}
	})
	return e
}

// Dispatcher exposes the dispatcher for registering additional handlers.
func (e *Engine) Dispatcher() *dispatch.Dispatcher { return e.disp }

// Store exposes the persistence layer (may be nil when persistence is disabled).
func (e *Engine) Store() *store.Store { return e.store }

// Market exposes the market handler set (e.g. to attach a trade observer).
func (e *Engine) Market() *handlers.Market { return e.market }

// OnFeed/OnLog/OnState register UI callbacks.
func (e *Engine) OnFeed(fn func(handlers.CaptureEvent)) { e.onFeed = fn }
func (e *Engine) OnLog(fn func(LogLine))                { e.onLog = fn }
func (e *Engine) OnState(fn func(state.Snapshot))       { e.onState = fn }

// OnCaptureError registers a callback for capture start failures (e.g. a missing
// Npcap runtime). The argument is the raw error string.
func (e *Engine) OnCaptureError(fn func(string)) { e.onCaptureErr = fn }

// startListener starts capture and reports any error to the log + callback.
func (e *Engine) startListener() {
	if err := e.listener.Start(); err != nil {
		e.log("Capture error: " + err.Error())
		if e.onCaptureErr != nil {
			e.onCaptureErr(err.Error())
		}
	}
}

// Start begins uploading and capturing.
func (e *Engine) Start() error {
	e.Up.Start()
	e.market.Start()
	if e.specsSvc != nil {
		go e.specsSvc.Load()
	}
	if e.mobsSvc != nil {
		go e.mobsSvc.Load()
	}
	go e.startListener()
	return nil
}

// Stop halts capture and uploads.
func (e *Engine) Stop() {
	e.listener.Stop()
	e.market.Stop()
	e.Up.Stop()
	if e.store != nil {
		e.store.Close()
	}
}

// Mails returns recent persisted marketplace mails (empty if persistence off).
func (e *Engine) Mails(limit int) []store.Mail {
	if e.store == nil {
		return nil
	}
	m, err := e.store.ListMails(limit)
	if err != nil {
		e.log("ListMails error: " + err.Error())
		return nil
	}
	return m
}

// Trades returns recent persisted trades (empty if persistence off).
func (e *Engine) Trades(limit int) []store.Trade {
	if e.store == nil {
		return nil
	}
	t, err := e.store.ListTrades(limit)
	if err != nil {
		e.log("ListTrades error: " + err.Error())
		return nil
	}
	return t
}

// StartCapture/StopCapture toggle just the listener.
func (e *Engine) StartCapture() { go e.startListener() }
func (e *Engine) StopCapture() { e.listener.Stop() }

// Capturing reports listener status.
func (e *Engine) Capturing() bool { return e.listener.Running() }

// Inject feeds a raw UDP payload into the parser (used by tests/replay).
func (e *Engine) Inject(payload []byte) photon.PacketStatus {
	e.State.MarkPacket()
	return e.parser.ReceivePacket(payload)
}

// PacketsSeen returns the count of dispatched messages.
func (e *Engine) PacketsSeen() int64 { return e.packetsSeen.Load() }

func (e *Engine) pushFeed(ev handlers.CaptureEvent) {
	e.feedMu.Lock()
	e.feed = append(e.feed, ev)
	if len(e.feed) > 500 {
		e.feed = e.feed[len(e.feed)-500:]
	}
	e.feedMu.Unlock()
	if e.onFeed != nil {
		e.onFeed(ev)
	}
}

// Feed returns the most recent capture events (newest last).
func (e *Engine) Feed() []handlers.CaptureEvent {
	e.feedMu.Lock()
	defer e.feedMu.Unlock()
	out := make([]handlers.CaptureEvent, len(e.feed))
	copy(out, e.feed)
	return out
}

func (e *Engine) log(text string) {
	line := LogLine{Time: time.Now(), Text: text}
	e.logMu.Lock()
	e.logs = append(e.logs, line)
	if len(e.logs) > 1000 {
		e.logs = e.logs[len(e.logs)-1000:]
	}
	e.logMu.Unlock()
	if e.onLog != nil {
		e.onLog(line)
	}
}

// Log appends a line to the engine log (exported for the app layer).
func (e *Engine) Log(text string) { e.log(text) }

// Logs returns recent log lines.
func (e *Engine) Logs() []LogLine {
	e.logMu.Lock()
	defer e.logMu.Unlock()
	out := make([]LogLine, len(e.logs))
	copy(out, e.logs)
	return out
}
