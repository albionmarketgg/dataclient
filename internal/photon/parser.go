package photon

// PacketStatus is the result of parsing one Photon packet.
type PacketStatus int

const (
	StatusUndefined PacketStatus = iota
	StatusSuccess
	StatusEncrypted
	StatusInvalidCrc
	StatusInvalidHeader
	StatusDisconnectCommand
)

const (
	photonHeaderLength  = 12
	commandHeaderLength = 12
)

// Handler receives routed Photon messages. Implementations are dispatched on
// the caller's goroutine; keep them fast or hand off internally.
type Handler interface {
	HandleEvent(code EventCode, params map[byte]any)
	HandleRequest(code OperationCode, params map[byte]any)
	HandleResponse(code OperationCode, returnCode int16, debugMessage string, params map[byte]any)
}

type segment struct {
	totalLength  int
	payload      []byte
	bytesWritten int
}

// Parser parses Photon UDP payloads and dispatches messages to a Handler.
type Parser struct {
	handler Handler
	pending map[int32]*segment
	verbose bool

	// inspect, if set, is called for every parsed message including unknown
	// codes that routing would drop. For packet-inspection tools.
	inspect func(typ string, leadingCode byte, returnCode int16, params map[byte]any)
}

// NewParser builds a parser dispatching to handler.
func NewParser(handler Handler) *Parser {
	return &Parser{handler: handler, pending: make(map[int32]*segment)}
}

// SetInspector installs an observer fired for every parsed message before code
// routing, so unknown codes are seen too.
func (p *Parser) SetInspector(fn func(typ string, leadingCode byte, returnCode int16, params map[byte]any)) {
	p.inspect = fn
}

// ReceivePacket parses one UDP payload (everything after the UDP header).
func (p *Parser) ReceivePacket(payload []byte) PacketStatus {
	if len(payload) < photonHeaderLength {
		return StatusInvalidHeader
	}
	r := newReader(payload)
	r.pos = 2 // skip peerId
	flags, _ := r.readByte()
	commandCount, _ := r.readByte()
	r.pos += 8 // skip timestamp(4) + challenge(4)

	if flags == 1 {
		return StatusEncrypted
	}
	// flags == 0xCC indicates CRC, unused by live Albion traffic, so we skip
	// validation to avoid dropping valid packets.

	status := StatusUndefined
	for i := 0; i < int(commandCount); i++ {
		if !r.has(commandHeaderLength) {
			return StatusInvalidHeader
		}
		s := p.handleCommand(r)
		if s == StatusInvalidHeader {
			return StatusInvalidHeader
		}
		status = s
	}
	return status
}

func (p *Parser) handleCommand(r *reader) PacketStatus {
	commandType, _ := r.readByte()
	r.pos++ // channelId
	r.pos++ // commandFlags
	r.pos++ // reserved
	commandLength, err := r.readInt32BE()
	if err != nil {
		return StatusInvalidHeader
	}
	_, err = r.readInt32BE() // sequenceNumber
	if err != nil {
		return StatusInvalidHeader
	}
	cl := int(commandLength) - commandHeaderLength
	if cl < 0 || !r.has(cl) {
		return StatusInvalidHeader
	}

	switch commandType {
	case 4: // Disconnect
		return StatusDisconnectCommand
	case 6: // SendReliable
		return p.handleReliable(r, cl)
	case 7: // SendUnreliable: 4-byte prefix then reliable body
		if cl < 4 {
			return StatusInvalidHeader
		}
		r.pos += 4
		return p.handleReliable(r, cl-4)
	case 8: // SendFragment
		return p.handleFragment(r, cl)
	default:
		r.pos += cl
		return StatusUndefined
	}
}

func (p *Parser) handleReliable(r *reader, commandLength int) PacketStatus {
	if commandLength < 2 {
		return StatusInvalidHeader
	}
	r.pos++ // signalByte
	messageType, _ := r.readByte()
	commandLength -= 2
	body, err := r.readBytes(commandLength)
	if err != nil {
		return StatusInvalidHeader
	}
	if messageType == 131 {
		return StatusEncrypted
	}
	mr := newReader(body)
	switch messageType {
	case 2: // OperationRequest
		req, err := mr.readOperationRequestBody()
		if err != nil {
			return StatusSuccess
		}
		if p.inspect != nil {
			p.inspect("request", byte(req.code), 0, req.params)
		}
		p.onRequest(req)
	case 3: // OperationResponse
		resp, err := mr.readOperationResponseBody()
		if err != nil {
			return StatusSuccess
		}
		if p.inspect != nil {
			p.inspect("response", byte(resp.code), resp.returnCode, resp.params)
		}
		p.onResponse(resp)
	case 4: // Event
		ev, err := mr.readEventDataBody()
		if err != nil {
			return StatusSuccess
		}
		if p.inspect != nil {
			p.inspect("event", ev.code, 0, ev.params)
		}
		p.onEvent(ev)
	}
	return StatusSuccess
}

func (p *Parser) handleFragment(r *reader, commandLength int) PacketStatus {
	if commandLength < 20 {
		return StatusInvalidHeader
	}
	startSeq, _ := r.readInt32BE()
	_, _ = r.readInt32BE() // fragmentCount
	_, _ = r.readInt32BE() // fragmentNumber
	totalLength, _ := r.readInt32BE()
	fragmentOffset, _ := r.readInt32BE()
	fragLen := commandLength - 20
	frag, err := r.readBytes(fragLen)
	if err != nil {
		return StatusInvalidHeader
	}

	seg := p.pending[startSeq]
	if seg == nil {
		if totalLength < 0 {
			return StatusInvalidHeader
		}
		seg = &segment{totalLength: int(totalLength), payload: make([]byte, totalLength)}
		p.pending[startSeq] = seg
	}
	fo := int(fragmentOffset)
	if fo < 0 || fo+fragLen > len(seg.payload) {
		delete(p.pending, startSeq)
		return StatusUndefined
	}
	copy(seg.payload[fo:], frag)
	seg.bytesWritten += fragLen
	if seg.bytesWritten >= seg.totalLength {
		delete(p.pending, startSeq)
		fr := newReader(seg.payload)
		return p.handleReliable(fr, len(seg.payload))
	}
	return StatusSuccess
}

// ---- message bodies ----

type operationRequest struct {
	code   OperationCode
	params map[byte]any
}

type operationResponse struct {
	code         OperationCode
	returnCode   int16
	debugMessage string
	params       map[byte]any
}

type eventData struct {
	code   byte
	params map[byte]any
}

func (r *reader) readOperationRequestBody() (operationRequest, error) {
	code, err := r.readByte()
	if err != nil {
		return operationRequest{}, err
	}
	params, err := r.parameterTable()
	if err != nil {
		return operationRequest{}, err
	}
	return operationRequest{code: OperationCode(code), params: params}, nil
}

func (r *reader) readOperationResponseBody() (operationResponse, error) {
	code, err := r.readByte()
	if err != nil {
		return operationResponse{}, err
	}
	returnCode, err := r.readInt16LE()
	if err != nil {
		return operationResponse{}, err
	}
	debug := ""
	if r.remaining() > 0 {
		dt, err := r.readByte()
		if err != nil {
			return operationResponse{}, err
		}
		dv, err := r.deserializeWithType(dt)
		if err != nil {
			return operationResponse{}, err
		}
		if s, ok := dv.(string); ok {
			debug = s
		}
	}
	params, err := r.parameterTable()
	if err != nil {
		return operationResponse{}, err
	}
	return operationResponse{code: OperationCode(code), returnCode: returnCode, debugMessage: debug, params: params}, nil
}

func (r *reader) readEventDataBody() (eventData, error) {
	code, err := r.readByte()
	if err != nil {
		return eventData{}, err
	}
	params, err := r.parameterTable()
	if err != nil {
		return eventData{}, err
	}
	return eventData{code: code, params: params}, nil
}

// nested-value type codes 24/25/26
func (r *reader) readOperationRequest() (operationRequest, error)   { return r.readOperationRequestBody() }
func (r *reader) readOperationResponse() (operationResponse, error) { return r.readOperationResponseBody() }
func (r *reader) readEventData() (eventData, error)                 { return r.readEventDataBody() }

// ---- routing ----

const (
	paramEventCode     byte = 252 // 0xFC
	paramOperationCode byte = 253 // 0xFD
)

func (p *Parser) onEvent(ev eventData) {
	if ev.code == 3 {
		if _, ok := ev.params[paramEventCode]; !ok {
			ev.params[paramEventCode] = int16(EvMove)
		}
	}
	code, ok := parseEventCode(ev.params)
	if !ok {
		return
	}
	if p.handler != nil {
		p.handler.HandleEvent(code, ev.params)
	}
}

func (p *Parser) onRequest(req operationRequest) {
	code, ok := parseOperationCode(req.params)
	if !ok {
		return
	}
	if p.handler != nil {
		p.handler.HandleRequest(code, req.params)
	}
}

func (p *Parser) onResponse(resp operationResponse) {
	code, ok := parseOperationCode(resp.params)
	if !ok {
		return
	}
	if p.handler != nil {
		p.handler.HandleResponse(code, resp.returnCode, resp.debugMessage, resp.params)
	}
}

func parseEventCode(params map[byte]any) (EventCode, bool) {
	raw, ok := toInt16(params[paramEventCode])
	if !ok {
		return 0, false
	}
	code := EventCode(raw)
	if IsKnownEventCode(code) {
		return code, true
	}
	// packed-nibble fallback: real code shifted left 4 with low nibble 1
	u := uint16(raw)
	if u&0x0f == 0x01 {
		shifted := EventCode(u >> 4)
		if IsKnownEventCode(shifted) {
			return shifted, true
		}
	}
	return 0, false
}

func parseOperationCode(params map[byte]any) (OperationCode, bool) {
	raw, ok := toInt16(params[paramOperationCode])
	if !ok {
		return 0, false
	}
	code := OperationCode(raw)
	if IsKnownOperationCode(code) {
		return code, true
	}
	return 0, false
}

// toInt16 coerces byte/short/int parameter representations to int16.
func toInt16(v any) (int16, bool) {
	switch n := v.(type) {
	case int16:
		return n, true
	case int32:
		return int16(n), true
	case int64:
		return int16(n), true
	case byte:
		return int16(n), true
	case uint16:
		return int16(n), true
	default:
		return 0, false
	}
}
