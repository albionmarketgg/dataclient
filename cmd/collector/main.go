// Command collector is a local server that accepts the data client's uploads,
// decodes them, aggregates live stats, and optionally appends every raw payload
// to a JSONL file. It mirrors the real ingest's PoW handshake and endpoints, so
// pointing the client's "Ingest base URL" at it just works.
//
//	go run ./cmd/collector --addr :8787 --out capture.jsonl
//	# then open http://localhost:8787 in a browser for live stats
package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/albionmarketgg/dataclient/internal/pow"
)

type collector struct {
	mu     sync.Mutex
	issued map[string]string
	out    *os.File

	difficulty int
	started    time.Time

	topicCounts map[string]int
	orders      int
	buys        int
	sells       int
	items       map[string]int
	locations   map[string]int
	goldSamples int
	histories   int
	emv         int
	recent      []rawRecord

	itemsData []byte
	itemsETag string
	itemsTime time.Time
}

type rawRecord struct {
	Time     time.Time       `json:"time"`
	Topic    string          `json:"topic"`
	ServerID string          `json:"serverid"`
	Body     json.RawMessage `json:"body"`
}

type orderUpload struct {
	Orders []struct {
		ItemTypeID  string `json:"ItemTypeId"`
		LocationID  string `json:"LocationId"`
		AuctionType string `json:"AuctionType"`
	} `json:"Orders"`
}

func main() {
	addr := flag.String("addr", ":8787", "listen address")
	out := flag.String("out", "", "optional JSONL file to append every raw payload to")
	difficulty := flag.Int("pow-bits", 8, "PoW difficulty (leading zero bits)")
	itemsSrc := flag.String("items", "https://raw.githubusercontent.com/ao-data/ao-bin-dumps/master/formatted/items.txt",
		"URL or local path to items.txt to serve at /items.txt (empty disables)")
	flag.Parse()

	c := &collector{
		issued:      map[string]string{},
		difficulty:  *difficulty,
		started:     time.Now(),
		topicCounts: map[string]int{},
		items:       map[string]int{},
		locations:   map[string]int{},
	}
	if *out != "" {
		f, err := os.OpenFile(*out, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			log.Fatalf("open --out: %v", err)
		}
		c.out = f
		defer f.Close()
		log.Printf("appending raw payloads to %s", *out)
	}

	if *itemsSrc != "" {
		go c.loadItems(*itemsSrc)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/pow", c.handlePow)
	mux.HandleFunc("/pow/", c.handleUpload)
	mux.HandleFunc("/itemEstimatedMarketValues", c.handleEMV)
	mux.HandleFunc("/items.txt", c.handleItems)
	mux.HandleFunc("/stats", c.handleStats)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "ok") })
	mux.HandleFunc("/", c.handleDashboard)

	log.Printf("Albion Market local collector on %s (pow=%d bits). Open http://localhost%s", *addr, *difficulty, normalizeAddr(*addr))
	log.Fatal(http.ListenAndServe(*addr, mux))
}

func normalizeAddr(a string) string {
	if strings.HasPrefix(a, ":") {
		return a
	}
	return ":" + a
}

func (c *collector) handlePow(w http.ResponseWriter, r *http.Request) {
	kb := make([]byte, 16)
	rand.Read(kb)
	key := hex.EncodeToString(kb)
	wanted := strings.Repeat("0", c.difficulty)
	c.mu.Lock()
	c.issued[key] = wanted
	c.mu.Unlock()
	writeJSON(w, pow.Request{Key: key, Wanted: wanted})
}

func (c *collector) handleUpload(w http.ResponseWriter, r *http.Request) {
	topic := strings.TrimPrefix(r.URL.Path, "/pow/")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	key := r.Form.Get("key")
	c.mu.Lock()
	wanted, ok := c.issued[key]
	if ok {
		delete(c.issued, key)
	}
	c.mu.Unlock()
	if !ok || !pow.Verify(pow.Request{Key: key, Wanted: wanted}, r.Form.Get("solution")) {
		http.Error(w, "invalid pow", http.StatusForbidden)
		return
	}
	c.ingest(topic, r.Form.Get("serverid"), []byte(r.Form.Get("natsmsg")))
	w.WriteHeader(http.StatusOK)
}

func (c *collector) handleEMV(w http.ResponseWriter, r *http.Request) {
	b, _ := io.ReadAll(r.Body)
	c.ingest("itemEstimatedMarketValues", "", b)
	w.WriteHeader(http.StatusOK)
}

func (c *collector) ingest(topic, serverID string, body []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.topicCounts[topic]++

	switch topic {
	case "marketorders.ingest":
		var ou orderUpload
		if json.Unmarshal(body, &ou) == nil {
			for _, o := range ou.Orders {
				c.orders++
				c.items[o.ItemTypeID]++
				c.locations[o.LocationID]++
				switch o.AuctionType {
				case "offer":
					c.sells++
				case "request":
					c.buys++
				}
			}
		}
	case "goldprices.ingest":
		var g struct {
			Prices []int64 `json:"Prices"`
		}
		if json.Unmarshal(body, &g) == nil {
			c.goldSamples += len(g.Prices)
		}
	case "markethistories.ingest":
		c.histories++
	case "itemEstimatedMarketValues":
		var e struct {
			Items []json.RawMessage `json:"items"`
		}
		if json.Unmarshal(body, &e) == nil {
			c.emv += len(e.Items)
		}
	}

	rec := rawRecord{Time: time.Now(), Topic: topic, ServerID: serverID, Body: json.RawMessage(body)}
	c.recent = append(c.recent, rec)
	if len(c.recent) > 100 {
		c.recent = c.recent[len(c.recent)-100:]
	}
	if c.out != nil {
		if line, err := json.Marshal(rec); err == nil {
			c.out.Write(append(line, '\n'))
		}
	}
	log.Printf("[%-26s] server=%s bytes=%d", topic, serverID, len(body))
}

type statsView struct {
	Uptime    string
	Topics    map[string]int
	Orders    int
	Buys      int
	Sells     int
	Gold      int
	Histories int
	EMV       int
	TopItems  []kv
	TopLocs   []kv
	Recent    []rawRecord
}

type kv struct {
	Key   string
	Count int
}

func (c *collector) snapshot() statsView {
	c.mu.Lock()
	defer c.mu.Unlock()
	tc := map[string]int{}
	for k, v := range c.topicCounts {
		tc[k] = v
	}
	recent := make([]rawRecord, len(c.recent))
	copy(recent, c.recent)
	// reverse for newest-first
	for i, j := 0, len(recent)-1; i < j; i, j = i+1, j-1 {
		recent[i], recent[j] = recent[j], recent[i]
	}
	return statsView{
		Uptime:    time.Since(c.started).Round(time.Second).String(),
		Topics:    tc,
		Orders:    c.orders,
		Buys:      c.buys,
		Sells:     c.sells,
		Gold:      c.goldSamples,
		Histories: c.histories,
		EMV:       c.emv,
		TopItems:  topN(c.items, 15),
		TopLocs:   topN(c.locations, 10),
		Recent:    recent,
	}
}

func topN(m map[string]int, n int) []kv {
	out := make([]kv, 0, len(m))
	for k, v := range m {
		out = append(out, kv{k, v})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Count > out[j].Count })
	if len(out) > n {
		out = out[:n]
	}
	return out
}

func (c *collector) loadItems(src string) {
	var data []byte
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		resp, err := http.Get(src)
		if err != nil {
			log.Printf("items: fetch failed: %v", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Printf("items: fetch status %d", resp.StatusCode)
			return
		}
		data, _ = io.ReadAll(resp.Body)
	} else {
		b, err := os.ReadFile(src)
		if err != nil {
			log.Printf("items: read failed: %v", err)
			return
		}
		data = b
	}
	sum := sha256.Sum256(data)
	c.mu.Lock()
	c.itemsData = data
	c.itemsETag = `"` + hex.EncodeToString(sum[:8]) + `"`
	c.itemsTime = time.Now().UTC()
	c.mu.Unlock()
	log.Printf("items: serving /items.txt (%d bytes, etag %s)", len(data), c.itemsETag)
}

func (c *collector) handleItems(w http.ResponseWriter, r *http.Request) {
	c.mu.Lock()
	data, etag, mod := c.itemsData, c.itemsETag, c.itemsTime
	c.mu.Unlock()
	if data == nil {
		http.Error(w, "items.txt not loaded yet", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("ETag", etag)
	w.Header().Set("Last-Modified", mod.Format(http.TimeFormat))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if match := r.Header.Get("If-None-Match"); match != "" && match == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Write(data)
}

func (c *collector) handleStats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, c.snapshot())
}

func (c *collector) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	dashTmpl.Execute(w, c.snapshot())
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

var dashTmpl = template.Must(template.New("d").Funcs(template.FuncMap{
	"pretty": func(b json.RawMessage) string {
		var out bytes.Buffer
		if json.Indent(&out, b, "", "  ") != nil {
			return string(b)
		}
		return out.String()
	},
}).Parse(`<!doctype html><html><head><meta charset="utf-8"><meta http-equiv="refresh" content="2">
<title>Albion Market — Local Collector</title>
<style>
 body{background:#161616;color:#fff;font:14px Roboto,system-ui,sans-serif;margin:0;padding:24px}
 h1{font-size:18px;margin:0 0 4px} .sub{color:#6f757c;margin-bottom:18px}
 .cards{display:grid;grid-template-columns:repeat(auto-fit,minmax(150px,1fr));gap:12px;margin-bottom:18px}
 .card{background:#25282b;border:1px solid #383c41;border-radius:10px;padding:14px}
 .card .l{font-size:11px;color:#6f757c;text-transform:uppercase} .card .v{font-size:26px;font-weight:600;margin-top:4px}
 .accent{color:#5865f2}.good{color:#43b581}.bad{color:#f04747}
 .cols{display:grid;grid-template-columns:1fr 1fr;gap:18px}
 .panel{background:#25282b;border:1px solid #383c41;border-radius:10px;overflow:hidden;margin-bottom:18px}
 .panel h3{margin:0;padding:11px 14px;font-size:13px;border-bottom:1px solid #2c2f33}
 table{width:100%;border-collapse:collapse} td,th{padding:6px 14px;border-bottom:1px solid #2c2f33;text-align:left;font-size:13px}
 th{color:#6f757c;font-size:11px;text-transform:uppercase} pre{margin:0;padding:10px 14px;font-size:11px;color:#a0a0a0;white-space:pre-wrap;border-bottom:1px solid #2c2f33}
</style></head><body>
<h1>Albion Market — Local Collector</h1>
<div class="sub">Live capture stats · uptime {{.Uptime}} · auto-refresh 2s</div>
<div class="cards">
 <div class="card"><div class="l">Orders</div><div class="v accent">{{.Orders}}</div></div>
 <div class="card"><div class="l">Sell offers</div><div class="v bad">{{.Sells}}</div></div>
 <div class="card"><div class="l">Buy requests</div><div class="v good">{{.Buys}}</div></div>
 <div class="card"><div class="l">Gold samples</div><div class="v">{{.Gold}}</div></div>
 <div class="card"><div class="l">History uploads</div><div class="v">{{.Histories}}</div></div>
 <div class="card"><div class="l">EMV items</div><div class="v">{{.EMV}}</div></div>
</div>
<div class="cols">
 <div class="panel"><h3>Top items</h3><table><tr><th>Item</th><th>Orders</th></tr>
  {{range .TopItems}}<tr><td>{{.Key}}</td><td>{{.Count}}</td></tr>{{end}}</table></div>
 <div class="panel"><h3>Top locations</h3><table><tr><th>Location id</th><th>Orders</th></tr>
  {{range .TopLocs}}<tr><td>{{.Key}}</td><td>{{.Count}}</td></tr>{{end}}</table></div>
</div>
<div class="panel"><h3>Recent raw payloads (exactly what we capture)</h3>
 {{range .Recent}}<pre><b style="color:#9fb0ff">{{.Topic}}</b> server={{.ServerID}} @ {{.Time.Format "15:04:05"}}
{{pretty .Body}}</pre>{{end}}
</div>
</body></html>`))
