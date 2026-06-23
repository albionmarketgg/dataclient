// Package items resolves Albion item ids to unique/display names. Data is loaded
// from our ingest/CDN base (configurable); it is optional and nil-safe.
package items

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Service maps numeric item ids to names.
type Service struct {
	itemsURL  string
	cachePath string // local cache of items.txt ("" disables caching)

	mu      sync.RWMutex
	unique  map[int]string
	display map[int]string
	emv     map[emvKey]int64
	loaded  atomic.Bool
}

type emvKey struct {
	id      int
	quality int
}

// New creates a Service that loads item data from itemsURL (a full URL to an
// items.txt in the public ao-bin-dumps format). cachePath, if set, caches the file
// for conditional (only-if-newer) refresh. Empty itemsURL disables loading.
func New(itemsURL, cachePath string) *Service {
	return &Service{
		itemsURL:  strings.TrimSpace(itemsURL),
		cachePath: strings.TrimSpace(cachePath),
		unique:    map[int]string{},
		display:   map[int]string{},
		emv:       map[emvKey]int64{},
	}
}

// SetEMV records an estimated market value for an item+quality.
func (s *Service) SetEMV(id, quality int, emv int64) {
	if emv <= 0 {
		return
	}
	s.mu.Lock()
	s.emv[emvKey{id, quality}] = emv
	s.mu.Unlock()
}

// EMV returns the estimated market value for an item+quality, if known.
func (s *Service) EMV(id, quality int) (int64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.emv[emvKey{id, quality}]
	if !ok && quality != 1 {
		v, ok = s.emv[emvKey{id, 1}]
	}
	return v, ok
}

// Loaded reports whether item data has been loaded.
func (s *Service) Loaded() bool { return s.loaded.Load() }

// UniqueName returns the unique item name for an id (implements handlers.ItemNamer).
func (s *Service) UniqueName(id int) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.unique[id]
	return n, ok
}

// DisplayName returns the localized display name for an id.
func (s *Service) DisplayName(id int) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.display[id]
	return n, ok
}

// Load fetches item data conditionally (If-None-Match / If-Modified-Since): on 304
// or network failure it uses the local cache; on 200 it downloads, caches, and
// loads. Best-effort — on total failure the service stays empty.
func (s *Service) Load() {
	if s.itemsURL == "" {
		return
	}
	etag, modified := s.readCacheMeta()

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(http.MethodGet, s.itemsURL, nil)
	if err == nil {
		if etag != "" {
			req.Header.Set("If-None-Match", etag)
		}
		if modified != "" {
			req.Header.Set("If-Modified-Since", modified)
		}
		resp, derr := client.Do(req)
		if derr == nil {
			defer resp.Body.Close()
			switch resp.StatusCode {
			case http.StatusOK:
				body, rerr := io.ReadAll(resp.Body)
				if rerr == nil && s.parse(bytes.NewReader(body)) {
					s.writeCache(body, resp.Header.Get("ETag"), resp.Header.Get("Last-Modified"))
					return
				}
			case http.StatusNotModified:
				// fall through to cache load below
			}
		}
	}
	// network failed or 304: load from cache if present.
	if s.cachePath != "" {
		if f, ferr := os.Open(s.cachePath); ferr == nil {
			defer f.Close()
			s.parse(f)
		}
	}
}

func (s *Service) parse(r io.Reader) bool {
	uniq := map[int]string{}
	disp := map[int]string{}
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		// format: "<index>: <UniqueName> : <Localized Name>"
		parts := strings.SplitN(sc.Text(), ":", 3)
		if len(parts) < 2 {
			continue
		}
		id, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			continue
		}
		uniq[id] = strings.TrimSpace(parts[1])
		if len(parts) == 3 {
			disp[id] = strings.TrimSpace(parts[2])
		}
	}
	if len(uniq) == 0 {
		return false
	}
	s.mu.Lock()
	s.unique = uniq
	s.display = disp
	s.mu.Unlock()
	s.loaded.Store(true)
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
