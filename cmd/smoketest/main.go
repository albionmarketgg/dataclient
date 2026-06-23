// Command smoketest runs the capture-free pipeline (synthetic packets → parse →
// dispatch → PoW → upload) against a running ingest endpoint and prints the
// upload stats, verifying the data path end-to-end.
package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/niick1231/albionmarket_dataclient/internal/config"
	"github.com/niick1231/albionmarket_dataclient/internal/demo"
	"github.com/niick1231/albionmarket_dataclient/internal/engine"
)

func main() {
	ingest := flag.String("ingest", "http://localhost:8787", "ingest base URL")
	dur := flag.Duration("duration", 6*time.Second, "demo duration")
	flag.Parse()

	cfg := config.Default()
	cfg.IngestBaseURL = *ingest
	cfg.NetworkStartDelaySecs = 0

	e := engine.New(cfg, nil, "")
	e.OnLog(func(l engine.LogLine) { fmt.Printf("%s  %s\n", l.Time.Format("15:04:05"), l.Text) })
	e.Up.Start()
	defer e.Up.Stop()

	fmt.Printf("Smoke test → %s for %s\n", *ingest, *dur)
	demo.Run(e, *dur)
	time.Sleep(1500 * time.Millisecond) // let queued uploads drain

	s := e.Up.Stats()
	fmt.Println("---- upload stats ----")
	fmt.Printf("market orders : %d\n", s.MarketOrders)
	fmt.Printf("market history: %d\n", s.MarketHistory)
	fmt.Printf("gold prices   : %d\n", s.GoldPrices)
	fmt.Printf("emv items     : %d\n", s.EMV)
	fmt.Printf("queued        : %d\n", s.Queued)
	fmt.Printf("failed        : %d\n", s.Failed)
	if s.MarketOrders+s.GoldPrices == 0 {
		fmt.Println("RESULT: no uploads recorded — is the ingest server running?")
		return
	}
	fmt.Println("RESULT: pipeline OK ✓")
}
