package handlers

import "testing"

func TestParseMailBody(t *testing.T) {
	cases := []struct {
		name        string
		typ         int
		body        string
		partial     int
		total       int
		itemID      string
		totalSilver int64
		unit        float64
	}{
		{"sell_finished", mtSellFin, "1|T5_2H_SHAPESHIFTER_SET3@1|1549840000|1549840000", 1, 1, "T5_2H_SHAPESHIFTER_SET3@1", 154984, 154984},
		{"buy_finished", mtBuyFin, "10|T7_ALCHEMY_RARE_ENT|11000100000|1100010000", 10, 10, "T7_ALCHEMY_RARE_ENT", 1100010, 110001},
		{"buy_expired", mtBuyExp, "23|100|65450000000|T5_ALCHEMY_RARE_PANTHER|", 23, 100, "T5_ALCHEMY_RARE_PANTHER", 1955000, 85000},
		{"sell_expired", mtSellExp, "0|39|0|T7_JOURNAL_HUNTER_FULL|", 0, 39, "T7_JOURNAL_HUNTER_FULL", 0, 0},
		{"bm_sell_expired", mtBMSell, "6|53|4420680000|T6_OFF_HORN_KEEPER@1|", 6, 53, "T6_OFF_HORN_KEEPER@1", 442068, 73678},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			partial, total, itemID, totalSilver, totalTaxes, unit, ok := parseMailBody(c.typ, c.body)
			if !ok {
				t.Fatalf("parse failed")
			}
			if partial != c.partial || total != c.total || itemID != c.itemID {
				t.Fatalf("got partial=%d total=%d item=%q", partial, total, itemID)
			}
			if totalSilver != c.totalSilver {
				t.Fatalf("totalSilver got %d want %d", totalSilver, c.totalSilver)
			}
			if unit != c.unit {
				t.Fatalf("unit got %v want %v", unit, c.unit)
			}
			if totalTaxes != 0 {
				t.Fatalf("totalTaxes got %d want 0", totalTaxes)
			}
		})
	}
}

func TestParseMailBodyBadInput(t *testing.T) {
	if _, _, _, _, _, _, ok := parseMailBody(mtSellFin, "garbage|not|numbers|x"); ok {
		t.Fatal("expected parse failure on non-numeric input")
	}
}

func TestNormalizeUnitSilver(t *testing.T) {
	// Use binary-exact values so the test reflects true double rounding (matching
	// C# Math.Round(double, 2, AwayFromZero); 0.125 is exactly representable).
	cases := []struct{ in, want float64 }{
		{0.125, 0.13}, {-0.125, -0.13}, {100.0, 100.0}, {0.004, 0.0}, {73678.0, 73678.0},
	}
	for _, c := range cases {
		if got := normalizeUnitSilver(c.in); got != c.want {
			t.Fatalf("normalize(%v)=%v want %v", c.in, got, c.want)
		}
	}
}

func TestTicksToTime(t *testing.T) {
	// 2024-01-01T00:00:00Z in .NET ticks = 638_396_640_000_000_000
	got := ticksToTime(638396640000000000)
	if got.UTC().Format("2006-01-02") != "2024-01-01" {
		t.Fatalf("ticksToTime gave %v", got.UTC())
	}
}
