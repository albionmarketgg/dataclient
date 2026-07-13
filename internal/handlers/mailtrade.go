package handlers

import (
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/albionmarketgg/dataclient/internal/dispatch"
	"github.com/albionmarketgg/dataclient/internal/market"
	"github.com/albionmarketgg/dataclient/internal/photon"
	"github.com/albionmarketgg/dataclient/internal/state"
	"github.com/albionmarketgg/dataclient/internal/store"
)

// mail types (positional, mirror the marketplace mail summary codes).
const (
	mtUnknown = 0
	mtSellFin = 1 // MARKETPLACE_SELLORDER_FINISHED_SUMMARY
	mtBuyFin  = 2 // MARKETPLACE_BUYORDER_FINISHED_SUMMARY
	mtBuyExp  = 3 // MARKETPLACE_BUYORDER_EXPIRED_SUMMARY
	mtSellExp = 4 // MARKETPLACE_SELLORDER_EXPIRED_SUMMARY
	mtBMSell  = 5 // BLACKMARKET_SELLORDER_EXPIRED_SUMMARY
)

var mailTypeByName = map[string]int{
	"MARKETPLACE_SELLORDER_FINISHED_SUMMARY": mtSellFin,
	"MARKETPLACE_BUYORDER_FINISHED_SUMMARY":  mtBuyFin,
	"MARKETPLACE_BUYORDER_EXPIRED_SUMMARY":   mtBuyExp,
	"MARKETPLACE_SELLORDER_EXPIRED_SUMMARY":  mtSellExp,
	"BLACKMARKET_SELLORDER_EXPIRED_SUMMARY":  mtBMSell,
}

// auction type ints: unknown/offer/request.
const (
	atUnknown = 0
	atOffer   = 1
	atRequest = 2
)

// trade operation / type ints.
const (
	opBuy  = 0
	opSell = 1

	ttInstant = 0
	ttOrder   = 1
)

func mailAuctionType(t int) int {
	switch t {
	case mtSellFin, mtSellExp, mtBMSell:
		return atOffer
	case mtBuyFin, mtBuyExp:
		return atRequest
	default:
		return atUnknown
	}
}

// MailTrade captures marketplace mails (realized prices) and instant trades.
type MailTrade struct {
	st    *state.State
	store *store.Store
	namer ItemNamer
	feed  FeedFunc

	mu          sync.Mutex
	orderCache  []market.Order // FIFO, max 500
	orderIndex  map[uint64]struct{}
	unconfirmed *store.Trade
}

// NewMailTrade builds the mail/trade handler set.
func NewMailTrade(st *state.State, db *store.Store, namer ItemNamer, feed FeedFunc) *MailTrade {
	if feed == nil {
		feed = func(CaptureEvent) {}
	}
	return &MailTrade{st: st, store: db, namer: namer, feed: feed, orderIndex: map[uint64]struct{}{}}
}

// Register attaches mail/trade handlers to the dispatcher.
func (mt *MailTrade) Register(d *dispatch.Dispatcher) {
	d.OnResponse(photon.OpGetMailInfos, mt.handleMailInfos)
	d.OnResponse(photon.OpReadMail, mt.handleReadMail)
	d.OnRequest(photon.OpAuctionBuyOffer, mt.handleBuyRequest)
	d.OnResponse(photon.OpAuctionBuyOffer, mt.handleTradeConfirm)
	d.OnRequest(photon.OpAuctionSellSpecificItemRequest, mt.handleSellRequest)
	d.OnResponse(photon.OpAuctionSellSpecificItemRequest, mt.handleTradeConfirm)
}

// CacheOrders adds browsed market orders to the FIFO cache for trade correlation.
func (mt *MailTrade) CacheOrders(orders []market.Order) {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	for _, o := range orders {
		if _, ok := mt.orderIndex[o.ID]; ok {
			continue
		}
		if len(mt.orderCache) >= 500 {
			old := mt.orderCache[0]
			mt.orderCache = mt.orderCache[1:]
			delete(mt.orderIndex, old.ID)
		}
		mt.orderCache = append(mt.orderCache, o)
		mt.orderIndex[o.ID] = struct{}{}
	}
}

func (mt *MailTrade) lookupOrder(id uint64) (market.Order, bool) {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	for _, o := range mt.orderCache {
		if o.ID == id {
			return o, true
		}
	}
	return market.Order{}, false
}

// ---- mail ----

func (mt *MailTrade) handleMailInfos(_ int16, _ string, p map[byte]any) {
	srv := mt.st.Server()
	if srv == nil || mt.st.Snapshot().PlayerName == "" {
		return
	}
	player := mt.st.Snapshot().PlayerName
	ids := dispatch.Int64Slice(p[3])
	locs := dispatch.Strings(p[7])
	types := dispatch.Strings(p[11])
	recv := dispatch.Int64Slice(p[12])
	n := len(ids)
	added := 0
	for i := 0; i < n; i++ {
		typ := mtUnknown
		if i < len(types) {
			typ = mailTypeByName[strings.ToUpper(strings.TrimSpace(types[i]))]
		}
		if typ == mtUnknown {
			continue
		}
		raw := ""
		if i < len(locs) {
			raw = locs[i]
		}
		var received time.Time
		if i < len(recv) {
			received = ticksToTime(recv[i])
		}
		m := store.Mail{
			ID:             ids[i],
			AlbionServerID: srv.ID,
			AuctionType:    mailAuctionType(typ),
			Type:           typ,
			RawLocationID:  raw,
			LocationID:     resolveLocation(raw),
			PlayerName:     player,
			Received:       received,
		}
		if isNew, err := mt.store.InsertMailInfo(m); err == nil && isNew {
			added++
		}
	}
	if added > 0 {
		mt.feed(CaptureEvent{Time: time.Now(), Kind: "mail", Detail: "mail summaries", Count: added})
	}
}

func (mt *MailTrade) handleReadMail(_ int16, _ string, p map[byte]any) {
	id, ok := dispatch.Int64(p[0])
	if !ok {
		return
	}
	body, ok := dispatch.String(p[1])
	if !ok {
		return
	}
	exists, isSet, err := mt.store.MailIsSet(id)
	if err != nil || !exists || isSet {
		return
	}
	mail, ok, err := mt.store.GetMail(id)
	if err != nil || !ok {
		return
	}
	partial, total, itemID, totalSilver, totalTaxes, unit, parsed := parseMailBody(mail.Type, body)
	if !parsed {
		return
	}
	if err := mt.store.SetMailData(id, partial, total, itemID, totalSilver, totalTaxes, unit); err != nil {
		return
	}
	mt.feed(CaptureEvent{Time: time.Now(), Kind: "mail", Detail: "mail price data", Count: 1})

	// new mail body data → an Order-type trade (offer→Sell, request→Buy)
	op := opBuy
	if mail.AuctionType == atOffer {
		op = opSell
	}
	mt.store.InsertTrade(store.Trade{
		ID:             uuid.NewString(),
		AlbionServerID: mail.AlbionServerID,
		Amount:         partial,
		DateTime:       mail.Received,
		ItemID:         itemID,
		LocationID:     mail.LocationID,
		Operation:      op,
		PlayerName:     mail.PlayerName,
		QualityLevel:   0,
		RawLocationID:  mail.RawLocationID,
		Type:           ttOrder,
		UnitSilver:     normalizeUnitSilver(unit),
	})
}

// parseMailBody parses a marketplace mail body. Returns ok=false on any parse error.
func parseMailBody(typ int, body string) (partial, total int, itemID string, totalSilver, totalTaxes int64, unit float64, ok bool) {
	parts := strings.Split(body, "|")
	defer func() {
		if r := recover(); r != nil {
			partial, total, itemID, totalSilver, totalTaxes, unit, ok = 0, 0, "", 0, 0, 0, false
		}
	}()
	atoi := func(i int) int { return mustAtoi(parts[i]) }
	atol := func(i int) int64 { return mustAtoi64(parts[i]) }

	switch typ {
	case mtSellFin, mtBuyFin: // AMOUNT|ITEM|TOTAL|UNIT
		partial = atoi(0)
		total = partial
		itemID = parts[1]
		totalSilver = atol(2) / 10000
		unit = normalizeUnitSilver(float64(atol(3)) / 10000.0)
	case mtBuyExp: // BOUGHT|TOTAL|REFUND|ITEM|
		partial = atoi(0)
		total = atoi(1)
		refund := float64(atol(2)) / 10000.0
		itemID = parts[3]
		remaining := total - partial
		if remaining > 0 {
			unit = normalizeUnitSilver(refund / float64(remaining))
		}
		totalSilver = int64(roundAway(unit * float64(partial)))
	case mtSellExp, mtBMSell: // SOLD|TOTAL|TOTAL_SILVER|ITEM|
		partial = atoi(0)
		total = atoi(1)
		itemID = parts[3]
		totalSilver = atol(2) / 10000
		if partial != 0 {
			unit = normalizeUnitSilver(float64(totalSilver) / float64(partial))
		}
	default:
		return 0, 0, "", 0, 0, 0, false
	}
	return partial, total, itemID, totalSilver, totalTaxes, unit, true
}

// ---- instant trades ----

func (mt *MailTrade) handleBuyRequest(p map[byte]any) {
	amount, _ := dispatch.Int(p[1])
	orderID, ok := dispatch.Uint64(p[2])
	if !ok {
		return
	}
	mt.stageInstantTrade(orderID, amount)
}

func (mt *MailTrade) handleSellRequest(p map[byte]any) {
	orderID, ok := dispatch.Uint64(p[1])
	if !ok {
		return
	}
	amount, _ := dispatch.Int(p[4])
	mt.stageInstantTrade(orderID, amount)
}

func (mt *MailTrade) stageInstantTrade(orderID uint64, amount int) {
	if !mt.st.OkToUpload() {
		return
	}
	order, ok := mt.lookupOrder(orderID)
	if !ok {
		return
	}
	op := opSell
	if order.AuctionType == market.AuctionOffer {
		op = opBuy // instant-buy from a sell offer
	}
	srv := mt.st.Server()
	serverID := 0
	if srv != nil {
		serverID = srv.ID
	}
	unit := float64(order.UnitPriceSilver) / 10000.0
	t := store.Trade{
		ID:             uuid.NewString(),
		AlbionServerID: serverID,
		Amount:         amount,
		DateTime:       time.Now().UTC(),
		ItemID:         order.ItemTypeID,
		LocationID:     resolveLocation(order.LocationID),
		Operation:      op,
		PlayerName:     mt.st.Snapshot().PlayerName,
		QualityLevel:   int(order.QualityLevel),
		RawLocationID:  order.LocationID,
		Type:           ttInstant,
		UnitSilver:     normalizeUnitSilver(unit),
	}
	mt.mu.Lock()
	mt.unconfirmed = &t
	mt.mu.Unlock()
}

func (mt *MailTrade) handleTradeConfirm(_ int16, _ string, _ map[byte]any) {
	mt.mu.Lock()
	t := mt.unconfirmed
	mt.unconfirmed = nil
	mt.mu.Unlock()
	if t == nil {
		return
	}
	if err := mt.store.InsertTrade(*t); err == nil {
		mt.feed(CaptureEvent{Time: time.Now(), Kind: "trade", Detail: "instant trade", Count: t.Amount})
	}
}

// ---- helpers ----

func mustAtoi(s string) int {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		panic(err)
	}
	return n
}

func mustAtoi64(s string) int64 {
	n, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		panic(err)
	}
	return n
}

// normalizeUnitSilver rounds half away from zero to 2 decimal places.
func normalizeUnitSilver(x float64) float64 {
	return roundAway(x*100) / 100
}

func roundAway(x float64) float64 {
	if x < 0 {
		return -math.Floor(-x + 0.5)
	}
	return math.Floor(x + 0.5)
}

// ticksToTime converts .NET DateTime ticks (100ns since year 1) to time.Time.
func ticksToTime(ticks int64) time.Time {
	const ticksToUnixEpoch = 621355968000000000
	nanos := (ticks - ticksToUnixEpoch) * 100
	return time.Unix(0, nanos).UTC()
}

// resolveLocation best-effort maps a raw location id string to a numeric id.
// Numeric raw ids map directly; anything else yields 0.
func resolveLocation(raw string) int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0
	}
	if n, err := strconv.Atoi(raw); err == nil {
		return n
	}
	return 0
}
