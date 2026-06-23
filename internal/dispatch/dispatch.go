// Package dispatch routes parsed Photon messages to registered handlers by code.
package dispatch

import (
	"sync"

	"github.com/niick1231/albionmarket_dataclient/internal/photon"
)

// EventFunc handles an event's parameter table.
type EventFunc func(params map[byte]any)

// RequestFunc handles a request's parameter table.
type RequestFunc func(params map[byte]any)

// ResponseFunc handles a response's parameter table.
type ResponseFunc func(returnCode int16, debugMessage string, params map[byte]any)

// Dispatcher implements photon.Handler and fans out to registered handlers.
type Dispatcher struct {
	mu        sync.RWMutex
	events    map[photon.EventCode][]EventFunc
	requests  map[photon.OperationCode][]RequestFunc
	responses map[photon.OperationCode][]ResponseFunc
	onAny     func()
}

// New creates an empty dispatcher.
func New() *Dispatcher {
	return &Dispatcher{
		events:    map[photon.EventCode][]EventFunc{},
		requests:  map[photon.OperationCode][]RequestFunc{},
		responses: map[photon.OperationCode][]ResponseFunc{},
	}
}

// OnAny registers a callback fired for every dispatched message (telemetry).
func (d *Dispatcher) OnAny(fn func()) { d.onAny = fn }

// OnEvent registers an event handler.
func (d *Dispatcher) OnEvent(code photon.EventCode, fn EventFunc) {
	d.mu.Lock()
	d.events[code] = append(d.events[code], fn)
	d.mu.Unlock()
}

// OnRequest registers a request handler.
func (d *Dispatcher) OnRequest(code photon.OperationCode, fn RequestFunc) {
	d.mu.Lock()
	d.requests[code] = append(d.requests[code], fn)
	d.mu.Unlock()
}

// OnResponse registers a response handler.
func (d *Dispatcher) OnResponse(code photon.OperationCode, fn ResponseFunc) {
	d.mu.Lock()
	d.responses[code] = append(d.responses[code], fn)
	d.mu.Unlock()
}

// HandleEvent implements photon.Handler.
func (d *Dispatcher) HandleEvent(code photon.EventCode, params map[byte]any) {
	if d.onAny != nil {
		d.onAny()
	}
	d.mu.RLock()
	fns := d.events[code]
	d.mu.RUnlock()
	for _, fn := range fns {
		fn(params)
	}
}

// HandleRequest implements photon.Handler.
func (d *Dispatcher) HandleRequest(code photon.OperationCode, params map[byte]any) {
	if d.onAny != nil {
		d.onAny()
	}
	d.mu.RLock()
	fns := d.requests[code]
	d.mu.RUnlock()
	for _, fn := range fns {
		fn(params)
	}
}

// HandleResponse implements photon.Handler.
func (d *Dispatcher) HandleResponse(code photon.OperationCode, returnCode int16, debugMessage string, params map[byte]any) {
	if d.onAny != nil {
		d.onAny()
	}
	d.mu.RLock()
	fns := d.responses[code]
	d.mu.RUnlock()
	for _, fn := range fns {
		fn(returnCode, debugMessage, params)
	}
}
