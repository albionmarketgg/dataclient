// Package handlers parses captured game messages into domain models and routes
// them to the uploader and the live UI feed.
package handlers

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/albionmarketgg/data-client/internal/dispatch"
	"github.com/albionmarketgg/data-client/internal/market"
	"github.com/albionmarketgg/data-client/internal/photon"
	"github.com/albionmarketgg/data-client/internal/state"
	"github.com/albionmarketgg/data-client/internal/upload"
)

// ItemNamer resolves an item's unique name from its numeric id.
type ItemNamer interface {
	UniqueName(id int) (string, bool)
}

// CaptureEvent is a single entry in the live UI feed.
type CaptureEvent struct {
	Time   time.Time `json:"time"`
	Kind   string    `json:"kind"`   // "offers","requests","history","gold","emv"
	Detail string    `json:"detail"` // short human description
	Count  int       `json:"count"`
}

// FeedFunc receives capture events for the UI.
type FeedFunc func(CaptureEvent)

type historyInfo struct {
	quality    uint8
	timescale  market.Timescale
	albionID   uint32
	locationID string
}

// Market wires market-data handlers onto a dispatcher.
type Market struct {
	st       *state.State
	up       *upload.Uploader
	namer    ItemNamer
	feed     FeedFunc
	tradeObs func([]market.Order)

	emvSink func(id, quality int, emv int64)

	mu          sync.Mutex
	historyRing [8192]*historyInfo

	// EMV from the EstimatedMarketValueUpdate event and per-item New*Item events,
	// deduped and flushed in batches on a ticker.
	emvMu       sync.Mutex
	emvPending  map[emvKey]int64  // (id,quality) -> latest emv (display silver)
	emvUploaded map[string]bool   // fingerprint dedup so a value uploads once/day
	emvCancel   context.CancelFunc
}

type emvKey struct{ id, quality int }

// SetEMVSink registers a local store for estimated market values seen on the wire.
func (m *Market) SetEMVSink(fn func(id, quality int, emv int64)) { m.emvSink = fn }

// NewMarket builds the market handler set.
func NewMarket(st *state.State, up *upload.Uploader, namer ItemNamer, feed FeedFunc) *Market {
	if feed == nil {
		feed = func(CaptureEvent) {}
	}
	return &Market{
		st: st, up: up, namer: namer, feed: feed,
		emvPending:  map[emvKey]int64{},
		emvUploaded: map[string]bool{},
	}
}

// OnOrders registers an observer for parsed orders (e.g. trade cache).
func (m *Market) OnOrders(fn func([]market.Order)) { m.tradeObs = fn }

// Register attaches all market handlers to the dispatcher.
func (m *Market) Register(d *dispatch.Dispatcher) {
	d.OnResponse(photon.OpAuctionGetOffers, func(_ int16, _ string, p map[byte]any) {
		m.handleOrders(p[0], market.AuctionOffer, "offers")
	})
	d.OnResponse(photon.OpAuctionGetRequests, func(_ int16, _ string, p map[byte]any) {
		m.handleOrders(p[0], market.AuctionRequest, "requests")
	})
	d.OnResponse(photon.OpAuctionGetLoadoutOffers, func(_ int16, _ string, p map[byte]any) {
		m.handleOrders(p[1], market.AuctionOffer, "offers")
	})
	d.OnRequest(photon.OpAuctionGetItemAverageStats, m.handleHistoryRequest)
	d.OnResponse(photon.OpAuctionGetItemAverageStats, m.handleHistoryResponse)
	d.OnResponse(photon.OpGoldMarketGetAverageInfo, m.handleGold)
	d.OnEvent(photon.EvEstimatedMarketValueUpdate, m.handleEMV)

	// Every New*Item event carries that item's EMV at param[4] (/10000) and id at param[1].
	for _, code := range []photon.EventCode{
		photon.EvNewSimpleItem, photon.EvNewEquipmentItem, photon.EvNewFurnitureItem,
		photon.EvNewJournalItem, photon.EvNewLaborerItem, photon.EvNewKillTrophyItem,
		photon.EvNewSiegeBannerItem,
	} {
		d.OnEvent(code, m.handleItemEMV)
	}
}

func (m *Market) handleOrders(raw any, at market.AuctionType, kind string) {
	if !m.st.OkToUpload() {
		return
	}
	var strs []string
	if kind == "offers" || kind == "requests" {
		strs = dispatch.Strings(raw)
		if strs == nil {
			strs = dispatch.StringMatrix(raw)
		}
	} else {
		strs = dispatch.StringMatrix(raw)
	}
	loc := m.st.Location()
	var orders []market.Order
	for _, s := range strs {
		var o market.Order
		if err := json.Unmarshal([]byte(s), &o); err != nil {
			continue
		}
		if o.LocationID == "" {
			o.LocationID = loc.ID
		}
		if o.AuctionType == "" {
			o.AuctionType = at
		}
		orders = append(orders, o)
	}
	if len(orders) == 0 {
		return
	}
	m.st.SetEncrypted(false)
	if m.tradeObs != nil {
		m.tradeObs(orders)
	}
	m.up.EnqueueMarket(market.Upload{Orders: orders})
	m.feed(CaptureEvent{Time: time.Now(), Kind: kind, Detail: string(at) + " orders", Count: len(orders)})
}

func (m *Market) handleHistoryRequest(p map[byte]any) {
	albionID, _ := dispatch.Int(p[1])
	quality, _ := dispatch.Int(p[2])
	ts, _ := dispatch.Int(p[3])
	msgID, ok := dispatch.Uint64(p[255])
	if !ok {
		return
	}
	loc := m.st.Location()
	m.mu.Lock()
	m.historyRing[msgID%8192] = &historyInfo{
		quality:    uint8(quality),
		timescale:  market.Timescale(ts),
		albionID:   uint32(albionID),
		locationID: loc.ID,
	}
	m.mu.Unlock()
}

func (m *Market) handleHistoryResponse(_ int16, _ string, p map[byte]any) {
	if !m.st.OkToUpload() {
		return
	}
	msgID, ok := dispatch.Uint64(p[255])
	if !ok {
		return
	}
	m.mu.Lock()
	info := m.historyRing[msgID%8192]
	m.historyRing[msgID%8192] = nil
	m.mu.Unlock()
	if info == nil {
		return
	}
	amounts := dispatch.Int64Slice(p[0])
	silver := dispatch.Uint64Slice(p[1])
	stamps := dispatch.Uint64Slice(p[2])
	n := len(amounts)
	if len(silver) < n {
		n = len(silver)
	}
	if len(stamps) < n {
		n = len(stamps)
	}
	hist := make([]market.History, 0, n)
	for i := 0; i < n; i++ {
		amt := amounts[i]
		if amt < 0 { // byte-wraparound fixup
			amt = 256 + amt
		}
		hist = append(hist, market.History{
			ItemAmount:   uint64(amt),
			SilverAmount: silver[i],
			Timestamp:    stamps[i],
		})
	}
	if len(hist) == 0 {
		return
	}
	m.up.EnqueueHistories(market.HistoriesUpload{
		AlbionID:        info.albionID,
		LocationID:      info.locationID,
		QualityLevel:    info.quality,
		Timescale:       info.timescale,
		MarketHistories: hist,
	})
	m.feed(CaptureEvent{Time: time.Now(), Kind: "history", Detail: "item history", Count: len(hist)})
}

func (m *Market) handleGold(_ int16, _ string, p map[byte]any) {
	prices := dispatch.Uint32Slice(p[0])
	stamps := dispatch.Int64Slice(p[1])
	if len(prices) == 0 || len(prices) != len(stamps) {
		return
	}
	m.up.EnqueueGold(market.GoldPriceUpload{Prices: prices, Timestamps: stamps})
	m.feed(CaptureEvent{Time: time.Now(), Kind: "gold", Detail: "gold prices", Count: len(prices)})
}

// handleEMV captures the bulk EstimatedMarketValueUpdate event (event 464).
func (m *Market) handleEMV(p map[byte]any) {
	normalIDs := dispatch.Int64Slice(p[0])
	normalEMV := dispatch.Int64Slice(p[1])
	for i := 0; i < len(normalIDs) && i < len(normalEMV); i++ {
		m.recordEMV(int(normalIDs[i]), 1, normalEMV[i])
	}
	explicitIDs := dispatch.Int64Slice(p[2])
	explicitQual := dispatch.Int64Slice(p[3])
	explicitEMV := dispatch.Int64Slice(p[4])
	for i := 0; i < len(explicitIDs) && i < len(explicitQual) && i < len(explicitEMV); i++ {
		m.recordEMV(int(explicitIDs[i]), int(explicitQual[i]), explicitEMV[i])
	}
}

// handleItemEMV captures the EMV carried by a single New*Item event:
// id = param[1], quality = param[6] (default 1), emv = param[4] (/10000).
func (m *Market) handleItemEMV(p map[byte]any) {
	id, _ := dispatch.Int(p[1])
	emv, _ := dispatch.Int64(p[4])
	quality := 1
	if q, ok := dispatch.Int(p[6]); ok && q > 0 {
		quality = q
	}
	m.recordEMV(id, quality, emv)
}

// recordEMV stores an estimated value (raw, /10000) into the local store and the
// pending upload batch (deduped by item+quality, latest value wins).
func (m *Market) recordEMV(id, quality int, emvRaw int64) {
	emv := emvRaw / 10000
	if id <= 0 || emv <= 0 {
		return
	}
	if quality <= 0 {
		quality = 1
	}
	if m.emvSink != nil {
		m.emvSink(id, quality, emv)
	}
	m.emvMu.Lock()
	m.emvPending[emvKey{id, quality}] = emv
	m.emvMu.Unlock()
}

// flushEMV uploads the pending EMV batch (deduped so a value uploads once/day).
func (m *Market) flushEMV() {
	srv := m.st.Server()
	if srv == nil || m.namer == nil {
		return
	}
	m.emvMu.Lock()
	pending := m.emvPending
	m.emvPending = map[emvKey]int64{}
	m.emvMu.Unlock()
	if len(pending) == 0 {
		return
	}

	day := time.Now().UTC().Format("2006-01-02")
	var items []market.EstimatedValueEntry
	m.emvMu.Lock()
	if len(m.emvUploaded) > 50000 {
		m.emvUploaded = map[string]bool{}
	}
	for k, emv := range pending {
		fp := strconv.Itoa(srv.ID) + "/" + strconv.Itoa(k.id) + "/" + strconv.Itoa(k.quality) + "/" + day + "/" + strconv.FormatInt(emv, 10)
		if m.emvUploaded[fp] {
			continue
		}
		name, ok := m.namer.UniqueName(k.id)
		if !ok || name == "" {
			continue
		}
		m.emvUploaded[fp] = true
		items = append(items, market.EstimatedValueEntry{ItemUniqueName: name, EMV: emv, Quality: k.quality, Day: day})
	}
	m.emvMu.Unlock()
	if len(items) == 0 {
		return
	}
	for i := 0; i < len(items); i += 500 {
		end := i + 500
		if end > len(items) {
			end = len(items)
		}
		m.up.UploadEMV(market.EstimatedValueUpload{ServerID: srv.ID, Items: items[i:end]})
	}
	m.feed(CaptureEvent{Time: time.Now(), Kind: "emv", Detail: "estimated values", Count: len(items)})
}

// Start begins the periodic EMV batch flush.
func (m *Market) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	m.emvCancel = cancel
	go func() {
		t := time.NewTicker(8 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				m.flushEMV()
				return
			case <-t.C:
				m.flushEMV()
			}
		}
	}()
}

// Stop ends the EMV flush loop (final flush attempted).
func (m *Market) Stop() {
	if m.emvCancel != nil {
		m.emvCancel()
	}
}
