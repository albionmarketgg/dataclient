// Shared wire types for the Albion Market data-client -> ingest service contract.
// Mirror this file into the refactor's `packages/shared` when building the real
// ingest endpoint. Keep it in sync with internal/market/models.go.
//
// See docs/INGEST_CONTRACT.md for transport (PoW handshake) details.

/** Region id used in the `serverid` form field / `serverId` JSON field. */
export type ServerId = 1 | 2 | 3; // 1=Americas, 2=Asia, 3=Europe

export type AuctionType = "offer" | "request";

/** Timescale for market history (matches the in-game enum). */
export enum Timescale {
  Day = 0,
  Week = 1,
  Month = 2,
}

/** A single market buy/sell order. PascalCase matches the in-game serialization. */
export interface Order {
  Id: number;
  ItemTypeId: string;       // e.g. "T4_BAG"
  ItemGroupTypeId: string;
  LocationId: string;       // market location id as string
  QualityLevel: number;     // 1..5
  EnchantmentLevel: number; // 0..4
  UnitPriceSilver: number;  // RAW silver (x10000); divide by 10000 to display
  Amount: number;
  AuctionType: AuctionType;
  Expires: string;          // ISO timestamp
}

/** POST /pow/marketorders.ingest (form field `natsmsg`). */
export interface MarketOrdersUpload {
  Orders: Order[];
}

export interface MarketHistoryPoint {
  ItemAmount: number;
  SilverAmount: number;
  Timestamp: number;
}

/** POST /pow/markethistories.ingest (form field `natsmsg`). */
export interface MarketHistoriesUpload {
  AlbionId: number;
  LocationId: string;
  QualityLevel: number;
  Timescale: Timescale;
  MarketHistories: MarketHistoryPoint[];
}

/** POST /pow/goldprices.ingest (form field `natsmsg`). */
export interface GoldPriceUpload {
  Prices: number[];
  Timestamps: number[];
}

export interface EstimatedValueEntry {
  itemUniqueName: string;
  emv: number;   // already divided by 10000 (display silver)
  quality: number;
  day: string;   // "YYYY-MM-DD" (UTC)
}

/** POST /itemEstimatedMarketValues (application/json). */
export interface EstimatedValueUpload {
  serverId: ServerId;
  items: EstimatedValueEntry[];
}

/** GET /pow response. */
export interface PowChallenge {
  key: string;
  wanted: string; // length = required leading zero bits
}

/** Fields posted with a PoW-gated upload (x-www-form-urlencoded). */
export interface PowUploadForm {
  key: string;
  solution: string;   // 16 hex chars
  serverid: string;   // ServerId as string
  natsmsg: string;    // JSON body (one of the *Upload types above)
  identifier: string; // client idempotency id
}

export const INGEST_TOPICS = {
  marketOrders: "marketorders.ingest",
  marketHistories: "markethistories.ingest",
  goldPrices: "goldprices.ingest",
} as const;
