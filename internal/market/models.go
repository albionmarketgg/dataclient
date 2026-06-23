// Package market defines the market-data domain models and their wire shapes.
package market

// AuctionType is serialized as the string the game/ingest expects.
type AuctionType string

const (
	AuctionOffer   AuctionType = "offer"
	AuctionRequest AuctionType = "request"
	AuctionUnknown AuctionType = "unknown"
)

// Order is a single market buy/sell order. JSON keys are PascalCase to match the
// in-game serialization (orders arrive as JSON strings in auction responses).
type Order struct {
	ID              uint64      `json:"Id"`
	ItemTypeID      string      `json:"ItemTypeId"`
	ItemGroupTypeID string      `json:"ItemGroupTypeId"`
	LocationID      string      `json:"LocationId"`
	QualityLevel    uint8       `json:"QualityLevel"`
	EnchantmentLevel uint8      `json:"EnchantmentLevel"`
	UnitPriceSilver uint64      `json:"UnitPriceSilver"`
	Amount          uint32      `json:"Amount"`
	AuctionType     AuctionType `json:"AuctionType"`
	Expires         string      `json:"Expires"`
}

// Upload is the market-orders upload payload.
type Upload struct {
	Orders []Order `json:"Orders"`
}

// Timescale matches the in-game history timescale enum.
type Timescale int

const (
	TimescaleDay   Timescale = 0
	TimescaleWeek  Timescale = 1
	TimescaleMonth Timescale = 2
)

// History is a single market-history data point.
type History struct {
	ItemAmount  uint64 `json:"ItemAmount"`
	SilverAmount uint64 `json:"SilverAmount"`
	Timestamp   uint64 `json:"Timestamp"`
}

// HistoriesUpload is the market-history upload payload.
type HistoriesUpload struct {
	AlbionID        uint32    `json:"AlbionId"`
	LocationID      string    `json:"LocationId"`
	QualityLevel    uint8     `json:"QualityLevel"`
	Timescale       Timescale `json:"Timescale"`
	MarketHistories []History `json:"MarketHistories"`
}

// GoldPriceUpload is the gold-price upload payload.
type GoldPriceUpload struct {
	Prices     []uint32 `json:"Prices"`
	Timestamps []int64  `json:"Timestamps"`
}

// EstimatedValueEntry is a single item EMV.
type EstimatedValueEntry struct {
	ItemUniqueName string `json:"itemUniqueName"`
	EMV            int64  `json:"emv"`
	Quality        int    `json:"quality"`
	Day            string `json:"day"`
}

// EstimatedValueUpload is the EMV batch upload payload (camelCase ingest shape).
type EstimatedValueUpload struct {
	ServerID int                   `json:"serverId"`
	Items    []EstimatedValueEntry `json:"items"`
}
