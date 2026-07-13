// Command inspector is a local developer tool that captures Albion's Photon UDP
// traffic, parses every event/request/response (including unknown codes the main
// client drops), and serves a live packet log on a local web page.
//
//	go run ./cmd/inspector --addr :8799 --out packets.jsonl
//	# then open http://localhost:8799
//
// Requires Npcap; captures both directions of udp port 5056.
package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/albionmarketgg/dataclient/internal/capture"
	"github.com/albionmarketgg/dataclient/internal/config"
	"github.com/albionmarketgg/dataclient/internal/photon"
	"github.com/albionmarketgg/dataclient/internal/state"
)

type nopHandler struct{}

func (nopHandler) HandleEvent(photon.EventCode, map[byte]any)                  {}
func (nopHandler) HandleRequest(photon.OperationCode, map[byte]any)            {}
func (nopHandler) HandleResponse(photon.OperationCode, int16, string, map[byte]any) {}

type record struct {
	Seq        int            `json:"seq"`
	Time       string         `json:"time"`
	Type       string         `json:"type"` // event|request|response
	Code       int            `json:"code"` // routing code (param 252/253) or leading byte
	Name       string         `json:"name"`
	Known      bool           `json:"known"`
	ReturnCode int            `json:"returnCode"`
	NumParams  int            `json:"numParams"`
	Params     map[string]any `json:"params"`
}

type inspector struct {
	mu     sync.Mutex
	seq    int
	recs   []record
	counts map[string]int
	server string
	out    *jsonlWriter
}

const ringCap = 8000

func main() {
	addr := flag.String("addr", ":8799", "web UI listen address")
	out := flag.String("out", "", "optional JSONL file to append every packet to")
	flag.Parse()

	ins := &inspector{counts: map[string]int{}}
	if *out != "" {
		w, err := newJSONL(*out)
		if err != nil {
			log.Fatalf("open --out: %v", err)
		}
		ins.out = w
		defer w.Close()
		log.Printf("appending packets to %s", *out)
	}

	cfg := config.Default()
	st := state.New()
	st.OnChange(func(s state.Snapshot) {
		ins.mu.Lock()
		ins.server = s.ServerName
		ins.mu.Unlock()
	})
	parser := photon.NewParser(nopHandler{})
	parser.SetInspector(ins.onMessage)

	listener := capture.New(cfg, st, parser, func(m string) { log.Println(m) })
	if err := listener.Start(); err != nil {
		log.Printf("capture: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/recent", ins.handleRecent)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "ok") })
	mux.HandleFunc("/", handleIndex)

	log.Printf("packet inspector on %s — open http://localhost%s  (filter: %s)", *addr, normAddr(*addr), cfg.PacketFilter)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

func normAddr(a string) string {
	if len(a) > 0 && a[0] == ':' {
		return a
	}
	return ":" + a
}

func (ins *inspector) onMessage(typ string, leadingCode byte, returnCode int16, params map[byte]any) {
	var code int
	var name string
	var known bool
	if typ == "event" {
		c, ok := codeOf(params, 252)
		if !ok {
			c = int(leadingCode)
		}
		code = c
		ec := photon.EventCode(c)
		name = ec.String()
		known = photon.IsKnownEventCode(ec)
	} else {
		c, ok := codeOf(params, 253)
		if !ok {
			c = int(leadingCode)
		}
		code = c
		oc := photon.OperationCode(c)
		name = oc.String()
		known = photon.IsKnownOperationCode(oc)
	}

	ins.mu.Lock()
	ins.seq++
	rec := record{
		Seq: ins.seq, Time: time.Now().Format("15:04:05.000"), Type: typ,
		Code: code, Name: name, Known: known, ReturnCode: int(returnCode),
		NumParams: len(params), Params: renderParams(params),
	}
	ins.recs = append(ins.recs, rec)
	if len(ins.recs) > ringCap {
		ins.recs = ins.recs[len(ins.recs)-ringCap:]
	}
	ins.counts[typ]++
	if !known {
		ins.counts["unknown"]++
	}
	ins.mu.Unlock()

	if ins.out != nil {
		ins.out.Write(rec)
	}
}

func (ins *inspector) handleRecent(w http.ResponseWriter, r *http.Request) {
	since, _ := strconv.Atoi(r.URL.Query().Get("since"))
	ins.mu.Lock()
	var recs []record
	for _, rec := range ins.recs {
		if rec.Seq > since {
			recs = append(recs, rec)
		}
	}
	// cap the first big load
	if since == 0 && len(recs) > 1500 {
		recs = recs[len(recs)-1500:]
	}
	resp := map[string]any{"seq": ins.seq, "server": ins.server, "counts": ins.counts, "records": recs}
	ins.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func codeOf(params map[byte]any, idx byte) (int, bool) {
	v, ok := params[idx]
	if !ok {
		return 0, false
	}
	switch n := v.(type) {
	case int16:
		return int(n), true
	case int32:
		return int(n), true
	case int64:
		return int(n), true
	case byte:
		return int(n), true
	case uint16:
		return int(n), true
	}
	return 0, false
}

// renderParams renders a parameter table to a JSON-safe, typed structure so each
// value's exact wire type stays visible.
func renderParams(params map[byte]any) map[string]any {
	out := make(map[string]any, len(params))
	for k, v := range params {
		out[strconv.Itoa(int(k))] = renderVal(v)
	}
	return out
}

func typed(t string, v any) map[string]any { return map[string]any{"t": t, "v": v} }

func renderVal(v any) any {
	switch x := v.(type) {
	case nil:
		return typed("null", nil)
	case bool:
		return typed("bool", x)
	case string:
		return typed("string", x)
	case byte:
		return typed("byte", x)
	case int16:
		return typed("int16", x)
	case int32:
		return typed("int32", x)
	case int64:
		return typed("int64", x)
	case uint16:
		return typed("uint16", x)
	case float32:
		return typed("float32", x)
	case float64:
		return typed("float64", x)
	case []byte:
		return typed("bytes["+strconv.Itoa(len(x))+"]", hex.EncodeToString(x))
	case []bool:
		return typed("bool["+strconv.Itoa(len(x))+"]", x)
	case photon.CustomData:
		return typed("custom", map[string]any{"type": x.TypeCode, "hex": hex.EncodeToString(x.Data)})
	case []photon.CustomData:
		arr := make([]any, len(x))
		for i, e := range x {
			arr[i] = renderVal(e)
		}
		return typed("custom["+strconv.Itoa(len(x))+"]", arr)
	case []any:
		arr := make([]any, len(x))
		for i, e := range x {
			arr[i] = renderVal(e)
		}
		return typed("array["+strconv.Itoa(len(x))+"]", arr)
	case map[any]any:
		m := make(map[string]any, len(x))
		for kk, vv := range x {
			m[fmt.Sprint(kk)] = renderVal(vv)
		}
		return typed("dict", m)
	default:
		return typed(fmt.Sprintf("%T", v), fmt.Sprintf("%v", v))
	}
}

// ---- JSONL writer ----

type jsonlWriter struct {
	mu sync.Mutex
	f  *os.File
}

func newJSONL(path string) (*jsonlWriter, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	return &jsonlWriter{f: f}, nil
}

func (w *jsonlWriter) Write(rec record) {
	b, err := json.Marshal(rec)
	if err != nil {
		return
	}
	w.mu.Lock()
	w.f.Write(append(b, '\n'))
	w.mu.Unlock()
}
func (w *jsonlWriter) Close() error { return w.f.Close() }
