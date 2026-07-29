// Package protomap consumes the remotely-served protocol map (protocol-map.json)
// and turns it into wire→compiled translations applied at the parser boundary.
// The remote map can move packet CODES and param POSITIONS (data); it can never
// define new decoding logic — handlers keep reading compiled positions, and a
// translation rewrites incoming packets into that compiled layout first.
// Contract: PROTOCOL_MAP_ENDPOINT.md.
package protomap

import (
	"encoding/json"

	"github.com/albionmarketgg/dataclient/internal/photon"
)

// baselineVersion is the version of the compiled-in defaults. Remote payloads
// with a lower version are ignored (stale), equal versions are accepted (they
// should be a 1:1 mirror = identity translation).
const baselineVersion = 1

// payloadPacket is one packet entry of the remote JSON. Param values are int
// positions in v1; v2 reserves objects (typed extractors) which are skipped
// ("keep compiled default") for forward compatibility.
type payloadPacket struct {
	Code           *int16                     `json:"code"`
	Params         map[string]json.RawMessage `json:"params"`
	RequestParams  map[string]json.RawMessage `json:"requestParams"`
	ResponseParams map[string]json.RawMessage `json:"responseParams"`
}

type payload struct {
	Version    int                      `json:"version"`
	Events     map[string]payloadPacket `json:"events"`
	Operations map[string]payloadPacket `json:"operations"`
}

// move is one param translation: the value at wire position `from` becomes the
// value at compiled position `to`.
type move struct{ to, from byte }

// Remap holds the computed wire→compiled translations. Empty maps = identity.
type Remap struct {
	eventCodes map[photon.EventCode]photon.EventCode // wire code -> compiled code
	opCodes    map[photon.OperationCode]photon.OperationCode
	eventMoves map[photon.EventCode][]move // keyed by COMPILED code
	reqMoves   map[photon.OperationCode][]move
	respMoves  map[photon.OperationCode][]move
}

// parseRemap validates a payload and computes the translation set against the
// compiled defaults. Returns nil when the payload must be rejected wholesale
// (bad JSON, stale version, out-of-range positions) — the caller then keeps the
// compiled defaults, so a bad push can never be worse than no push.
func parseRemap(body []byte) *Remap {
	var p payload
	if err := json.Unmarshal(body, &p); err != nil || p.Version < baselineVersion {
		return nil
	}
	r := &Remap{
		eventCodes: map[photon.EventCode]photon.EventCode{},
		opCodes:    map[photon.OperationCode]photon.OperationCode{},
		eventMoves: map[photon.EventCode][]move{},
		reqMoves:   map[photon.OperationCode][]move{},
		respMoves:  map[photon.OperationCode][]move{},
	}
	for name, pp := range p.Events {
		def, ok := defaultEvents[name]
		if !ok {
			continue // unknown packet name: ignore (forward compat)
		}
		compiled := photon.EventCode(def.code)
		if pp.Code != nil && *pp.Code != def.code {
			if *pp.Code < 0 {
				return nil
			}
			r.eventCodes[photon.EventCode(*pp.Code)] = compiled
		}
		moves, valid := paramMoves(def.params, pp.Params)
		if !valid {
			return nil
		}
		if len(moves) > 0 {
			r.eventMoves[compiled] = moves
		}
	}
	for name, pp := range p.Operations {
		reqDef, hasReq := defaultOpRequests[name]
		respDef, hasResp := defaultOpResponses[name]
		if !hasReq && !hasResp {
			continue
		}
		code := reqDef.code
		if !hasReq {
			code = respDef.code
		}
		compiled := photon.OperationCode(code)
		if pp.Code != nil && *pp.Code != code {
			if *pp.Code < 0 {
				return nil
			}
			r.opCodes[photon.OperationCode(*pp.Code)] = compiled
		}
		// "params" applies to both directions; requestParams/responseParams override.
		if hasReq {
			moves, valid := paramMoves(reqDef.params, merged(pp.Params, pp.RequestParams))
			if !valid {
				return nil
			}
			if len(moves) > 0 {
				r.reqMoves[compiled] = moves
			}
		}
		if hasResp {
			moves, valid := paramMoves(respDef.params, merged(pp.Params, pp.ResponseParams))
			if !valid {
				return nil
			}
			if len(moves) > 0 {
				r.respMoves[compiled] = moves
			}
		}
	}
	return r
}

// merged overlays specific direction params over the shared ones.
func merged(shared, specific map[string]json.RawMessage) map[string]json.RawMessage {
	if len(specific) == 0 {
		return shared
	}
	if len(shared) == 0 {
		return specific
	}
	out := make(map[string]json.RawMessage, len(shared)+len(specific))
	for k, v := range shared {
		out[k] = v
	}
	for k, v := range specific {
		out[k] = v
	}
	return out
}

// paramMoves diffs remote param positions against the compiled defaults.
// Unknown field names and v2 object values are skipped; an int position outside
// 0-255 invalidates the whole payload.
func paramMoves(def map[string]byte, remote map[string]json.RawMessage) ([]move, bool) {
	var out []move
	for name, raw := range remote {
		to, known := def[name]
		if !known {
			continue
		}
		var pos int
		if err := json.Unmarshal(raw, &pos); err != nil {
			continue // v2 typed-extractor object (or junk): keep compiled default
		}
		if pos < 0 || pos > 255 {
			return nil, false
		}
		if byte(pos) != to {
			out = append(out, move{to: to, from: byte(pos)})
		}
	}
	return out, true
}

// translate applies moves copy-on-write: the value at each wire position lands
// at its compiled position; untouched params carry over. When the source param
// is absent, the compiled position is cleared — whatever the new layout put
// there is some OTHER field and must not be misread.
func translate(params map[byte]any, moves []move) map[byte]any {
	if len(moves) == 0 {
		return params
	}
	out := make(map[byte]any, len(params))
	for k, v := range params {
		out[k] = v
	}
	for _, mv := range moves {
		if v, ok := params[mv.from]; ok {
			out[mv.to] = v
		} else {
			delete(out, mv.to)
		}
	}
	return out
}

// Event implements the wire→compiled translation for events (photon.Remapper).
func (r *Remap) Event(code photon.EventCode, params map[byte]any) (photon.EventCode, map[byte]any) {
	if r == nil {
		return code, params
	}
	if c, ok := r.eventCodes[code]; ok {
		code = c
	}
	return code, translate(params, r.eventMoves[code])
}

// Request implements photon.Remapper for operation requests.
func (r *Remap) Request(code photon.OperationCode, params map[byte]any) (photon.OperationCode, map[byte]any) {
	if r == nil {
		return code, params
	}
	if c, ok := r.opCodes[code]; ok {
		code = c
	}
	return code, translate(params, r.reqMoves[code])
}

// Response implements photon.Remapper for operation responses.
func (r *Remap) Response(code photon.OperationCode, params map[byte]any) (photon.OperationCode, map[byte]any) {
	if r == nil {
		return code, params
	}
	if c, ok := r.opCodes[code]; ok {
		code = c
	}
	return code, translate(params, r.respMoves[code])
}

// identity reports whether the remap changes nothing (pure 1:1 mirror payload).
func (r *Remap) identity() bool {
	return r == nil || (len(r.eventCodes) == 0 && len(r.opCodes) == 0 &&
		len(r.eventMoves) == 0 && len(r.reqMoves) == 0 && len(r.respMoves) == 0)
}
