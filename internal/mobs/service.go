// Package mobs resolves an Albion mob index (the NewMob packet's param[1]) to a
// display name, from a processed_mobs.json dump ([{mobId, uniqueName, en}]) where
// mobId equals that packet index. Loading is optional and nil-safe: until the file
// is fetched, Name returns ("", false) and callers fall back to "Mob <index>".
//
// Mirrors internal/specs.Service (fetch + conditional cache), JSON instead of XML.
package mobs

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Service maps mob index -> display name.
type Service struct {
	url       string
	cachePath string

	mu     sync.RWMutex
	byID   map[int]string
	loaded atomic.Bool
}

// mobEntry is one row of processed_mobs.json.
type mobEntry struct {
	MobID      int    `json:"mobId"`
	UniqueName string `json:"uniqueName"`
	En         string `json:"en"`
}

// New creates a Service that loads processed_mobs.json from url, caching to
// cachePath. Empty url disables loading.
func New(url, cachePath string) *Service {
	return &Service{url: strings.TrimSpace(url), cachePath: strings.TrimSpace(cachePath), byID: map[int]string{}}
}

// Ready reports whether the mob map is loaded.
func (s *Service) Ready() bool { return s != nil && s.loaded.Load() }

// Name returns the display name for a mob index (ok=false if unknown / unloaded).
func (s *Service) Name(index int) (string, bool) {
	if s == nil {
		return "", false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.byID[index]
	return n, ok && n != ""
}

// Load fetches processed_mobs.json (conditional on the cached copy) and builds the
// index->name map. Best-effort: on any failure it falls back to the local cache.
func (s *Service) Load() {
	if s == nil || s.url == "" {
		return
	}
	etag, modified := s.readCacheMeta()
	client := &http.Client{Timeout: 60 * time.Second}
	if req, err := http.NewRequest(http.MethodGet, s.url, nil); err == nil {
		if etag != "" {
			req.Header.Set("If-None-Match", etag)
		}
		if modified != "" {
			req.Header.Set("If-Modified-Since", modified)
		}
		if resp, derr := client.Do(req); derr == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				if body, rerr := io.ReadAll(resp.Body); rerr == nil && s.parse(bytes.NewReader(body)) {
					s.writeCache(body, resp.Header.Get("ETag"), resp.Header.Get("Last-Modified"))
					return
				}
			}
		}
	}
	if s.cachePath != "" {
		if f, ferr := os.Open(s.cachePath); ferr == nil {
			defer f.Close()
			s.parse(f)
		}
	}
}

// parse decodes the JSON array into the index->name map (en, falling back to
// uniqueName). Returns false on decode error / empty input.
func (s *Service) parse(r io.Reader) bool {
	var entries []mobEntry
	if err := json.NewDecoder(r).Decode(&entries); err != nil || len(entries) == 0 {
		return false
	}
	m := make(map[int]string, len(entries))
	for _, e := range entries {
		name := e.En
		if name == "" {
			name = e.UniqueName
		}
		if name != "" {
			m[e.MobID] = name
		}
	}
	if len(m) == 0 {
		return false
	}
	s.mu.Lock()
	s.byID = m
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
