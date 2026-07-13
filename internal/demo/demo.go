// Package demo injects synthetic market traffic into the engine so the pipeline
// and UI can be exercised without the game running.
package demo

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/albionmarketgg/data-client/internal/config"
	"github.com/albionmarketgg/data-client/internal/engine"
	"github.com/albionmarketgg/data-client/internal/photon"
	"github.com/albionmarketgg/data-client/internal/phototest"
	"github.com/albionmarketgg/data-client/internal/state"
)

var items = []string{"T4_BAG", "T5_2H_BOW", "T6_HEAD_CLOTH_SET1", "T4_MOUNT_HORSE", "T7_OFF_SHIELD", "T8_MAIN_SWORD"}
var cities = []state.Location{
	{ID: "3005", Name: "Caerleon"},
	{ID: "0007", Name: "Thetford"},
	{ID: "1002", Name: "Lymhurst"},
	{ID: "4002", Name: "Bridgewatch"},
}

const (
	demoPlayerObj = int64(1001)
	demoMobObj    = int64(2001)
)

// Run injects synthetic traffic for the given duration.
func Run(e *engine.Engine, duration time.Duration) {
	e.Log(fmt.Sprintf("Demo mode: injecting synthetic traffic for %s.", duration))
	e.State.SetServer(&config.Server{ID: 1, Name: "Americas"})

	// Establish player + entities (Join, NewCharacter, NewMob, party).
	e.Inject(phototest.ResponsePacket(byte(photon.OpJoin), 0, map[byte]any{
		253: int16(photon.OpJoin), 0: demoPlayerObj, 1: "demo-hero-guid", 2: "DemoHero", 8: "3005",
	}))
	e.Inject(phototest.EventPacket(byte(photon.EvNewCharacter), map[byte]any{
		252: int16(photon.EvNewCharacter), 0: demoPlayerObj, 1: "DemoHero", 7: "demo-hero-guid",
	}))
	e.Inject(phototest.EventPacket(byte(photon.EvNewMob), map[byte]any{
		252: int16(photon.EvNewMob), 0: demoMobObj, 1: int32(7),
	}))
	e.Inject(phototest.EventPacket(byte(photon.EvPartyJoined), map[byte]any{
		252: int16(photon.EvPartyJoined), 5: []string{"demo-hero-guid", "demo-ally-guid"}, 6: []string{"DemoHero", "DemoAlly"},
	}))

	deadline := time.Now().Add(duration)
	for time.Now().Before(deadline) {
		city := cities[rand.Intn(len(cities))]
		e.State.SetLocation(city)
		e.State.MarkPacket()
		switch rand.Intn(7) {
		case 0, 1:
			Offers(e, city.ID, rand.Intn(2) == 0)
		case 2:
			Gold(e)
		case 3:
			EMV(e)
		case 4:
			Harvest(e)
		case 5:
			CombatHit(e)
		case 6:
			LootGrab(e)
		}
		time.Sleep(500 * time.Millisecond)
	}
	e.Log("Demo mode finished.")
}

// Harvest injects a gathering harvest event.
func Harvest(e *engine.Engine) {
	itemID := rand.Intn(4000) + 1
	e.Inject(phototest.EventPacket(byte(photon.EvHarvestFinished), map[byte]any{
		252: int16(photon.EvHarvestFinished), 0: demoPlayerObj,
		4: int32(itemID), 5: int32(rand.Intn(5) + 1), 6: int32(rand.Intn(3)), 7: int32(0),
	}))
}

// CombatHit injects a damage health-update (player hits mob).
func CombatHit(e *engine.Engine) {
	dmg := int64(rand.Intn(900) + 100)
	e.Inject(phototest.EventPacket(byte(photon.EvHealthUpdate), map[byte]any{
		252: int16(photon.EvHealthUpdate), 0: demoMobObj, 2: -dmg, 3: int64(5000), 6: demoPlayerObj,
	}))
	if rand.Intn(3) == 0 {
		e.Inject(phototest.EventPacket(byte(photon.EvUpdateFame), map[byte]any{
			252: int16(photon.EvUpdateFame), 2: int64((rand.Intn(500) + 50) * 10000),
		}))
	}
}

// LootGrab injects an OtherGrabbedLoot event.
func LootGrab(e *engine.Engine) {
	itemID := rand.Intn(4000) + 1
	players := []string{"DemoAlly", "RandomPlayer", "DemoHero"}
	// EventData code byte is a placeholder; routing uses param 252 (code > 255).
	e.Inject(phototest.EventPacket(0, map[byte]any{
		252: int16(photon.EvOtherGrabbedLoot), 1: "Mob Camp", 2: players[rand.Intn(len(players))],
		3: false, 4: int32(itemID), 5: int64(rand.Intn(10) + 1),
	}))
}

// Offers injects a market offers/requests response.
func Offers(e *engine.Engine, locID string, isOffer bool) {
	at := "offer"
	op := photon.OpAuctionGetOffers
	if !isOffer {
		at = "request"
		op = photon.OpAuctionGetRequests
	}
	var orders []string
	for i := 0; i < rand.Intn(4)+1; i++ {
		item := items[rand.Intn(len(items))]
		price := (rand.Intn(50000) + 100) * 10000
		amt := rand.Intn(20) + 1
		orders = append(orders, fmt.Sprintf(
			`{"Id":%d,"ItemTypeId":"%s","LocationId":"%s","QualityLevel":%d,"EnchantmentLevel":%d,"UnitPriceSilver":%d,"Amount":%d,"AuctionType":"%s","Expires":"2026-07-01T00:00:00"}`,
			rand.Int63(), item, locID, rand.Intn(5)+1, rand.Intn(4), price, amt, at))
	}
	pkt := phototest.ResponsePacket(byte(op), 0, map[byte]any{253: int16(op), 0: orders})
	e.Inject(pkt)
}

// Gold injects a gold-price response.
func Gold(e *engine.Engine) {
	n := rand.Intn(5) + 2
	prices := make([]int32, n)
	stamps := make([]int64, n)
	base := int32(4400 + rand.Intn(300))
	now := time.Now().Unix()
	for i := 0; i < n; i++ {
		prices[i] = base + int32(rand.Intn(40)-20)
		stamps[i] = now - int64((n-i)*60)
	}
	pkt := phototest.ResponsePacket(byte(photon.OpGoldMarketGetAverageInfo), 0, map[byte]any{
		253: int16(photon.OpGoldMarketGetAverageInfo), 0: prices, 1: stamps,
	})
	e.Inject(pkt)
}

// EMV injects an estimated-market-value event.
func EMV(e *engine.Engine) {
	n := rand.Intn(4) + 1
	ids := make([]int64, n)
	emvs := make([]int64, n)
	for i := 0; i < n; i++ {
		ids[i] = int64(rand.Intn(5000) + 1)
		emvs[i] = int64((rand.Intn(100000) + 1000) * 10000)
	}
	pkt := phototest.EventPacket(7, map[byte]any{
		252: int16(photon.EvEstimatedMarketValueUpdate), 0: ids, 1: emvs,
	})
	e.Inject(pkt)
}
