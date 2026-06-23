// Package upload ships captured market data to our ingest endpoint
// (config.IngestBaseURL).
package upload

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/niick1231/albionmarket_dataclient/internal/config"
	"github.com/niick1231/albionmarket_dataclient/internal/market"
	"github.com/niick1231/albionmarket_dataclient/internal/pow"
)

// Stats is an upload counter snapshot for the UI.
type Stats struct {
	Queued        int64 `json:"queued"`
	MarketOrders  int64 `json:"marketOrders"`
	MarketHistory int64 `json:"marketHistory"`
	GoldPrices    int64 `json:"goldPrices"`
	EMV           int64 `json:"emv"`
	Failed        int64 `json:"failed"`
}

type job struct {
	topic   string
	payload any
}

// Uploader serializes uploads, performs the PoW handshake, and POSTs to ingest.
type Uploader struct {
	cfg    config.Config
	client *http.Client
	logf   func(string)

	queue  chan job
	wg     sync.WaitGroup
	cancel context.CancelFunc

	stats     Stats
	serverID  atomic.Int64
	identifier func() string

	// tokenFn, when set, returns a bearer token + user id to attribute uploads
	// to a logged-in user. ok=false means anonymous.
	tokenFn func() (token, userID string, ok bool)
}

// SetTokenProvider wires an auth token source for attributing uploads.
func (u *Uploader) SetTokenProvider(fn func() (token, userID string, ok bool)) {
	u.tokenFn = fn
}

// New builds an Uploader.
func New(cfg config.Config, logf func(string)) *Uploader {
	if logf == nil {
		logf = func(string) {}
	}
	return &Uploader{
		cfg:        cfg,
		client:     &http.Client{Timeout: 15 * time.Second},
		logf:       logf,
		queue:      make(chan job, 256),
		identifier: defaultIdentifier,
	}
}

// SetServerID records the current server id used in uploads.
func (u *Uploader) SetServerID(id int) { u.serverID.Store(int64(id)) }

// Start launches the worker.
func (u *Uploader) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	u.cancel = cancel
	u.wg.Add(1)
	go u.worker(ctx)
}

// Stop drains and stops the worker.
func (u *Uploader) Stop() {
	if u.cancel != nil {
		u.cancel()
	}
	u.wg.Wait()
}

// Stats returns a snapshot of upload counters.
func (u *Uploader) Stats() Stats {
	return Stats{
		Queued:        atomic.LoadInt64(&u.stats.Queued),
		MarketOrders:  atomic.LoadInt64(&u.stats.MarketOrders),
		MarketHistory: atomic.LoadInt64(&u.stats.MarketHistory),
		GoldPrices:    atomic.LoadInt64(&u.stats.GoldPrices),
		EMV:           atomic.LoadInt64(&u.stats.EMV),
		Failed:        atomic.LoadInt64(&u.stats.Failed),
	}
}

// EnqueueMarket queues a market-orders upload.
func (u *Uploader) EnqueueMarket(up market.Upload) {
	if len(up.Orders) == 0 {
		return
	}
	u.enqueue(u.cfg.MarketOrdersTopic, up)
}

// EnqueueHistories queues a market-history upload.
func (u *Uploader) EnqueueHistories(up market.HistoriesUpload) {
	u.enqueue(u.cfg.MarketHistoriesTopic, up)
}

// EnqueueGold queues a gold-price upload.
func (u *Uploader) EnqueueGold(up market.GoldPriceUpload) {
	if len(up.Prices) == 0 {
		return
	}
	u.enqueue(u.cfg.GoldPricesTopic, up)
}

// UploadEMV ships an estimated-market-value batch (camelCase JSON, no PoW).
func (u *Uploader) UploadEMV(up market.EstimatedValueUpload) {
	if len(up.Items) == 0 {
		return
	}
	go func() {
		body, err := json.Marshal(up)
		if err != nil {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := u.post(ctx, u.cfg.IngestBaseURL+"/itemEstimatedMarketValues", "application/json", body); err != nil {
			atomic.AddInt64(&u.stats.Failed, 1)
			u.logf("EMV upload failed: " + err.Error())
			return
		}
		atomic.AddInt64(&u.stats.EMV, int64(len(up.Items)))
	}()
}

func (u *Uploader) enqueue(topic string, payload any) {
	atomic.AddInt64(&u.stats.Queued, 1)
	select {
	case u.queue <- job{topic: topic, payload: payload}:
	default:
		u.logf("Upload queue full; dropping " + topic)
		atomic.AddInt64(&u.stats.Queued, -1)
		atomic.AddInt64(&u.stats.Failed, 1)
	}
}

func (u *Uploader) worker(ctx context.Context) {
	defer u.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case j := <-u.queue:
			atomic.AddInt64(&u.stats.Queued, -1)
			if err := u.send(ctx, j); err != nil {
				atomic.AddInt64(&u.stats.Failed, 1)
				u.logf("Upload failed (" + j.topic + "): " + err.Error())
				continue
			}
			u.countSuccess(j.topic)
		}
	}
}

func (u *Uploader) countSuccess(topic string) {
	switch topic {
	case u.cfg.MarketOrdersTopic:
		atomic.AddInt64(&u.stats.MarketOrders, 1)
	case u.cfg.MarketHistoriesTopic:
		atomic.AddInt64(&u.stats.MarketHistory, 1)
	case u.cfg.GoldPricesTopic:
		atomic.AddInt64(&u.stats.GoldPrices, 1)
	}
}

func (u *Uploader) send(ctx context.Context, j job) error {
	body, err := json.Marshal(j.payload)
	if err != nil {
		return err
	}
	if !u.cfg.RequirePoW {
		return u.post(ctx, u.cfg.IngestBaseURL+"/"+j.topic, "application/json", body)
	}
	req, err := u.getPow(ctx)
	if err != nil {
		return fmt.Errorf("pow challenge: %w", err)
	}
	solution := pow.Solve(req)
	form := url.Values{}
	form.Set("key", req.Key)
	form.Set("solution", solution)
	form.Set("serverid", strconv.FormatInt(u.serverID.Load(), 10))
	form.Set("natsmsg", string(body))
	// attribute to the logged-in user when available, else an anonymous id.
	identifier := u.identifier()
	if u.tokenFn != nil {
		if _, userID, ok := u.tokenFn(); ok && userID != "" {
			identifier = userID
		}
	}
	form.Set("identifier", identifier)
	return u.post(ctx, u.cfg.IngestBaseURL+"/pow/"+j.topic, "application/x-www-form-urlencoded", []byte(form.Encode()))
}

func (u *Uploader) getPow(ctx context.Context) (pow.Request, error) {
	var pr pow.Request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.cfg.IngestBaseURL+"/pow", nil)
	if err != nil {
		return pr, err
	}
	resp, err := u.client.Do(req)
	if err != nil {
		return pr, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return pr, fmt.Errorf("status %d", resp.StatusCode)
	}
	b, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(b, &pr); err != nil {
		return pr, err
	}
	return pr, nil
}

func (u *Uploader) post(ctx context.Context, url, contentType string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", "AlbionMarketDataClient")
	if u.tokenFn != nil {
		if token, userID, ok := u.tokenFn(); ok {
			req.Header.Set("Authorization", "Bearer "+token)
			if userID != "" {
				req.Header.Set("X-User-Id", userID)
			}
		}
	}
	resp, err := u.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	return nil
}

func defaultIdentifier() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
