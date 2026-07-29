package protomap

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/albionmarketgg/dataclient/internal/photon"
)

// Service fetches protocol-map.json (mirrors mobs/specs: conditional GET with a
// local cache) and serves the active Remap to the parser. Nil-safe and
// identity-by-default: until a payload with actual differences loads, packets
// pass through untouched.
type Service struct {
	url       string
	cachePath string
	logf      func(string)

	cur atomic.Pointer[Remap]
}

// New creates the service. Empty url disables fetching (compiled defaults only).
func New(url, cachePath string, logf func(string)) *Service {
	if logf == nil {
		logf = func(string) {}
	}
	return &Service{url: strings.TrimSpace(url), cachePath: strings.TrimSpace(cachePath), logf: logf}
}

// Event/Request/Response implement photon.Remapper by delegating to the active
// remap (nil = identity). Hot path: one atomic load per packet.
func (s *Service) Event(c photon.EventCode, p map[byte]any) (photon.EventCode, map[byte]any) {
	if s == nil {
		return c, p
	}
	return s.cur.Load().Event(c, p)
}

func (s *Service) Request(c photon.OperationCode, p map[byte]any) (photon.OperationCode, map[byte]any) {
	if s == nil {
		return c, p
	}
	return s.cur.Load().Request(c, p)
}

func (s *Service) Response(c photon.OperationCode, p map[byte]any) (photon.OperationCode, map[byte]any) {
	if s == nil {
		return c, p
	}
	return s.cur.Load().Response(c, p)
}

// Load fetches the map (conditional on the cached copy) and activates it.
// Best-effort: on fetch failure it falls back to the cached file; on invalid
// payload it keeps whatever is active (ultimately the compiled defaults).
func (s *Service) Load() {
	if s == nil || s.url == "" {
		return
	}
	etag, modified := s.readCacheMeta()
	client := &http.Client{Timeout: 30 * time.Second}
	if req, err := http.NewRequest(http.MethodGet, s.url, nil); err == nil {
		if etag != "" {
			req.Header.Set("If-None-Match", etag)
		}
		if modified != "" {
			req.Header.Set("If-Modified-Since", modified)
		}
		if resp, derr := client.Do(req); derr == nil {
			defer resp.Body.Close()
			switch {
			case resp.StatusCode == http.StatusOK:
				if body, rerr := io.ReadAll(resp.Body); rerr == nil && s.activate(body, "remote") {
					s.writeCache(body, resp.Header.Get("ETag"), resp.Header.Get("Last-Modified"))
					return
				}
			case resp.StatusCode == http.StatusNotModified:
				// cached copy is current; fall through to (re)load it below.
			}
		}
	}
	if s.cachePath != "" {
		if body, ferr := os.ReadFile(s.cachePath); ferr == nil {
			s.activate(body, "cache")
		}
	}
}

// Poll runs Load immediately and then on every tick, for the app's lifetime.
func (s *Service) Poll(interval time.Duration) {
	s.Load()
	for range time.Tick(interval) {
		s.Load()
	}
}

// activate parses and swaps in a remap; logs when it actually translates.
func (s *Service) activate(body []byte, origin string) bool {
	r := parseRemap(bytes.TrimSpace(body))
	if r == nil {
		s.logf("Protocol map (" + origin + ") rejected; using compiled defaults.")
		return false
	}
	prev := s.cur.Load()
	s.cur.Store(r)
	if !r.identity() && (prev == nil || prev.identity()) {
		s.logf("Protocol map (" + origin + ") active: remote layout overrides applied.")
	}
	return true
}

func (s *Service) readCacheMeta() (etag, modified string) {
	if s.cachePath == "" {
		return "", ""
	}
	if b, err := os.ReadFile(s.cachePath + ".meta"); err == nil {
		lines := strings.SplitN(string(b), "\n", 2)
		if len(lines) > 0 {
			etag = strings.TrimSpace(lines[0])
		}
		if len(lines) > 1 {
			modified = strings.TrimSpace(lines[1])
		}
	}
	return etag, modified
}

func (s *Service) writeCache(body []byte, etag, modified string) {
	if s.cachePath == "" {
		return
	}
	if err := os.MkdirAll(filepath.Dir(s.cachePath), 0o755); err != nil {
		return
	}
	_ = os.WriteFile(s.cachePath, body, 0o644)
	_ = os.WriteFile(s.cachePath+".meta", []byte(etag+"\n"+modified), 0o644)
}
