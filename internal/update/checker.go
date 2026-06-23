// Package update asks the backend whether a newer client build is available.
package update

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Result mirrors GET /client/version?current=<v>.
type Result struct {
	Latest         string `json:"latest"`
	MinSupported   string `json:"minSupported"`
	DownloadURL    string `json:"downloadUrl"`
	Available      bool   `json:"available"`
	Notes          string `json:"notes"`
	Current        string `json:"current"`
	UpToDate       bool   `json:"upToDate"`
	UpdateRequired bool   `json:"updateRequired"`
}

// ShouldNotify reports whether the user should be told about an update: either a
// forced one (below the support floor) or an available newer build.
func (r Result) ShouldNotify() bool {
	return r.UpdateRequired || (r.Available && !r.UpToDate)
}

// Checker queries the version endpoint.
type Checker struct {
	baseURL string
	version string
	client  *http.Client
}

// New builds a Checker. baseURL is the site origin; version is this build.
func New(baseURL, version string) *Checker {
	return &Checker{
		baseURL: strings.TrimRight(baseURL, "/"),
		version: version,
		client:  &http.Client{Timeout: 15 * time.Second},
	}
}

// Check asks the backend whether a newer build exists (it computes upToDate /
// updateRequired from the current version we send).
func (c *Checker) Check(ctx context.Context) (Result, error) {
	var r Result
	u := c.baseURL + "/client/version?current=" + url.QueryEscape(c.version)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return r, err
	}
	req.Header.Set("Cache-Control", "no-store")
	resp, err := c.client.Do(req)
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return r, fmt.Errorf("status %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return r, err
	}
	if r.Current == "" {
		r.Current = c.version
	}
	return r, nil
}
