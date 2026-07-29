package protomap

import "github.com/albionmarketgg/dataclient/internal/photon"

// packetDefault is the compiled-in layout of one packet: its wire code and the
// param position of every semantic field the client consumes. These tables are
// the single source of truth that the remote protocol-map.json is diffed
// against — a remote entry that matches is a no-op; one that differs produces a
// wire→compiled translation applied before dispatch (see map.go), so handlers
// keep reading their compiled positions untouched.
//
// Names must stay in sync with the endpoint seed (PROTOCOL_MAP_ENDPOINT.md §5).
type packetDefault struct {
	code   int16
	params map[string]byte
}

var defaultEvents = map[string]packetDefault{
	// market / items
	"EstimatedMarketValueUpdate": {int16(photon.EvEstimatedMarketValueUpdate), map[string]byte{
		"itemIds": 0, "emvs": 1, "explicitIds": 2, "explicitQualities": 3, "explicitEmvs": 4,
	}},
	// NewEquipmentItem positions reflect the post-2026-07 wire layout.
	"NewEquipmentItem": {int16(photon.EvNewEquipmentItem), map[string]byte{
		"objectId": 0, "itemIndex": 1, "quantity": 2, "emv": 4, "quality": 7, "isAwakened": 11,
	}},
	"NewSimpleItem":      {int16(photon.EvNewSimpleItem), simpleItemParams()},
	"NewSiegeBannerItem": {int16(photon.EvNewSiegeBannerItem), simpleItemParams()},
	"NewFurnitureItem":   {int16(photon.EvNewFurnitureItem), simpleItemParams()},
	"NewKillTrophyItem":  {int16(photon.EvNewKillTrophyItem), simpleItemParams()},
	"NewJournalItem":     {int16(photon.EvNewJournalItem), simpleItemParams()},
	"NewLaborerItem":     {int16(photon.EvNewLaborerItem), simpleItemParams()},

	// awakened
	"NewEquipmentItemLegendarySoul": {int16(photon.EvNewEquipmentItemLegendarySoul), map[string]byte{
		"objectId": 0, "soulId": 1, "era": 3, "attunedToMe": 4, "attunedTo": 5,
		"strain": 6, "attunement": 7, "traitIds": 8, "traitRolls": 9, "attunementSpent": 12,
	}},

	// combat
	"HealthUpdate":  {int16(photon.EvHealthUpdate), map[string]byte{"targetId": 0, "change": 2, "causerId": 6}},
	"HealthUpdates": {int16(photon.EvHealthUpdates), map[string]byte{"targetId": 0, "changes": 2, "causers": 6}},
	"NewCharacter":  {int16(photon.EvNewCharacter), map[string]byte{"objectId": 0, "characterName": 1, "playerGuid": 7}},
	"NewMob":        {int16(photon.EvNewMob), map[string]byte{"objectId": 0, "mobTypeIndex": 1}},

	// fame / dungeon / gathering
	"UpdateFame": {int16(photon.EvUpdateFame), map[string]byte{
		"totalFame": 1, "gainedFame": 2, "zoneFame": 3, "multiplier": 4, "isPremium": 5,
		"bagInsightIndex": 8, "satchelFame": 10, "eventBonusFactor": 17,
	}},
	"TakeSilver": {int16(photon.EvTakeSilver), map[string]byte{"silver": 3}},
	"HarvestFinished": {int16(photon.EvHarvestFinished), map[string]byte{
		"gathererId": 0, "itemIndex": 4, "standardAmount": 5, "bonusAmount": 6, "premiumBonusAmount": 7,
	}},
	"RewardGranted": {int16(photon.EvRewardGranted), map[string]byte{"itemIndex": 1, "amount": 3}},

	// loot
	"OtherGrabbedLoot": {int16(photon.EvOtherGrabbedLoot), map[string]byte{
		"source": 1, "looterName": 2, "isSilver": 3, "itemIndex": 4, "amount": 5,
	}},
	"NewLoot":                   {int16(photon.EvNewLoot), map[string]byte{"objectId": 0, "sourceName": 3}},
	"NewLootChest":              {int16(photon.EvNewLootChest), map[string]byte{"objectId": 0, "chestName": 3}},
	"LootChestOpened":           {int16(photon.EvLootChestOpened), map[string]byte{"objectId": 0}},
	"AttachItemContainer":       {int16(photon.EvAttachItemContainer), map[string]byte{"objectId": 0, "containerGuid": 1, "slotItemIds": 3}},
	"DetachItemContainer":       {int16(photon.EvDetachItemContainer), map[string]byte{"containerGuid": 0}},
	"InventoryPutItem":          {int16(photon.EvInventoryPutItem), map[string]byte{"itemObjectId": 0}},
	"InventoryDeleteItem":       {int16(photon.EvInventoryDeleteItem), map[string]byte{"itemObjectId": 0}},
	"PartyLootItemTypesRemoved": {int16(photon.EvPartyLootItemTypesRemoved), map[string]byte{"itemTypeIds": 1, "isSilverFlags": 3}},

	// party
	"PartyJoined":               {int16(photon.EvPartyJoined), map[string]byte{"memberGuids": 5, "memberNames": 6}},
	"PartyPlayerJoined":         {int16(photon.EvPartyPlayerJoined), map[string]byte{"guid": 1, "name": 2}},
	"PartyPlayerLeft":           {int16(photon.EvPartyPlayerLeft), map[string]byte{"guid": 1}},
	"PartyDisbanded":            {int16(photon.EvPartyDisbanded), map[string]byte{}},
	"PartyOnClusterPartyJoined": {int16(photon.EvPartyOnClusterPartyJoined), map[string]byte{"memberGuids": 0}},
	"PartySetRoleFlag":          {int16(photon.EvPartySetRoleFlag), map[string]byte{"guid": 1}},

	// specs
	"FullAchievementInfo": {int16(photon.EvFullAchievementInfo), map[string]byte{
		"level100Indices": 1, "achievementIndices": 2, "levels": 3,
	}},
}

func simpleItemParams() map[string]byte {
	return map[string]byte{"objectId": 0, "itemIndex": 1, "amount": 2, "emv": 4, "quality": 6}
}

// Operations have distinct request/response layouts (positions may collide with
// different meanings, e.g. AuctionGetItemAverageStats), so their defaults are
// kept per direction.
var defaultOpRequests = map[string]packetDefault{
	"InventoryMoveItem": {int16(photon.OpInventoryMoveItem), map[string]byte{"sourceContainerGuid": 1, "destinationContainerGuid": 4}},
	"AuctionBuyOffer":   {int16(photon.OpAuctionBuyOffer), map[string]byte{"amount": 1, "orderId": 2}},
	"AuctionGetItemAverageStats": {int16(photon.OpAuctionGetItemAverageStats), map[string]byte{
		"itemId": 1, "quality": 2, "timestamp": 3, "messageId": 255,
	}},
	"AuctionSellSpecificItemRequest": {int16(photon.OpAuctionSellSpecificItemRequest), map[string]byte{"orderId": 1, "amount": 4}},
}

var defaultOpResponses = map[string]packetDefault{
	"Join": {int16(photon.OpJoin), map[string]byte{
		"objectId": 0, "guid": 1, "characterName": 2, "locationId": 8, "fame": 35,
	}},
	"AuctionGetOffers":   {int16(photon.OpAuctionGetOffers), map[string]byte{"orders": 0}},
	"AuctionGetRequests": {int16(photon.OpAuctionGetRequests), map[string]byte{"orders": 0}},
	"AuctionGetItemAverageStats": {int16(photon.OpAuctionGetItemAverageStats), map[string]byte{
		"amounts": 0, "silver": 1, "timestamps": 2, "messageId": 255,
	}},
	"GetMailInfos": {int16(photon.OpGetMailInfos), map[string]byte{
		"mailIds": 3, "locations": 7, "mailTypes": 11, "receivedTimestamps": 12,
	}},
	"ReadMail":                 {int16(photon.OpReadMail), map[string]byte{"mailId": 0, "body": 1}},
	"GoldMarketGetAverageInfo": {int16(photon.OpGoldMarketGetAverageInfo), map[string]byte{"goldPrices": 0, "timestamps": 1}},
	"AuctionGetLoadoutOffers":  {int16(photon.OpAuctionGetLoadoutOffers), map[string]byte{"orders": 1}},
}
