// Package specs resolves Destiny Board achievement indices to stable string ids
// (e.g. "CRAFT_QUARTERSTAFFS") from the public ao-bin-dumps achievements.xml. The
// index is a 0-based counter over the root's direct <achievement>/
// <templateachievement> children in document order, matching the in-game
// FullAchievementInfo packet. Loading is optional and nil-safe.
package specs

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Service maps achievement indices to ids.
type Service struct {
	url       string
	cachePath string

	mu      sync.RWMutex
	ids     []string // index -> id
	version string
	loaded  atomic.Bool
}

// New creates a Service that loads achievements.xml from url, caching to
// cachePath. Empty url disables loading.
func New(url, cachePath string) *Service {
	return &Service{url: strings.TrimSpace(url), cachePath: strings.TrimSpace(cachePath)}
}

// Ready reports whether the index→id map is loaded.
func (s *Service) Ready() bool { return s.loaded.Load() }

// Version returns the achievements.xml Version attribute (empty until loaded).
func (s *Service) Version() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.version
}

// IDByIndex returns the achievement id for an index (ok=false if out of range).
func (s *Service) IDByIndex(i int) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if i < 0 || i >= len(s.ids) {
		return "", false
	}
	id := s.ids[i]
	return id, id != ""
}

// Load fetches achievements.xml (conditional on the cached copy) and builds the
// index→id map. Best-effort: on any failure it falls back to the local cache.
func (s *Service) Load() {
	if s.url == "" {
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

// parse walks the XML, recording the id of every direct <achievement>/
// <templateachievement> child of the root, in order (matching the game index).
func (s *Service) parse(r io.Reader) bool {
	dec := xml.NewDecoder(r)
	var ids []string
	var version string
	depth := 0
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if depth == 0 {
				// root <achievements Version="...">
				for _, a := range t.Attr {
					if strings.EqualFold(a.Name.Local, "Version") {
						version = a.Value
					}
				}
			} else if depth == 1 {
				name := strings.ToLower(t.Name.Local)
				if name == "achievement" || name == "templateachievement" {
					id := ""
					for _, a := range t.Attr {
						if strings.EqualFold(a.Name.Local, "id") {
							id = a.Value
							break
						}
					}
					if id != "" {
						ids = append(ids, id)
					}
				}
			}
			depth++
		case xml.EndElement:
			depth--
		}
	}
	if len(ids) == 0 {
		return false
	}
	s.mu.Lock()
	s.ids = ids
	s.version = version
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
