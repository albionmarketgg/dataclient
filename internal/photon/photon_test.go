package photon

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
)

// ---- minimal Protocol18 encoder (mirror of the deserializer) for tests ----

func writeCount(b *bytes.Buffer, n int) {
	u := uint32(n)
	for {
		c := byte(u & 0x7f)
		u >>= 7
		if u != 0 {
			c |= 0x80
		}
		b.WriteByte(c)
		if u == 0 {
			return
		}
	}
}

func writeZigzag32(b *bytes.Buffer, v int32) {
	writeCount(b, 0) // placeholder unused
}

// encTyped writes [typeByte][payload].
func encTyped(b *bytes.Buffer, v any) {
	switch x := v.(type) {
	case int16:
		b.WriteByte(p18Short)
		_ = binary.Write(b, binary.LittleEndian, x)
	case int32:
		b.WriteByte(p18CompressedInt)
		u := uint32((x << 1) ^ (x >> 31))
		for {
			c := byte(u & 0x7f)
			u >>= 7
			if u != 0 {
				c |= 0x80
			}
			b.WriteByte(c)
			if u == 0 {
				break
			}
		}
	case int64:
		b.WriteByte(p18CompressedLong)
		u := uint64((x << 1) ^ (x >> 63))
		for {
			c := byte(u & 0x7f)
			u >>= 7
			if u != 0 {
				c |= 0x80
			}
			b.WriteByte(c)
			if u == 0 {
				break
			}
		}
	case byte:
		b.WriteByte(p18Byte)
		b.WriteByte(x)
	case bool:
		if x {
			b.WriteByte(p18BooleanTrue)
		} else {
			b.WriteByte(p18BooleanFalse)
		}
	case string:
		b.WriteByte(p18String)
		writeCount(b, len(x))
		b.WriteString(x)
	case nil:
		b.WriteByte(p18Null)
	case []string:
		b.WriteByte(p18StringArray)
		writeCount(b, len(x))
		for _, s := range x {
			writeCount(b, len(s))
			b.WriteString(s)
		}
	case []int32:
		b.WriteByte(p18CompressedIntArray)
		writeCount(b, len(x))
		for _, n := range x {
			u := uint32((n << 1) ^ (n >> 31))
			for {
				c := byte(u & 0x7f)
				u >>= 7
				if u != 0 {
					c |= 0x80
				}
				b.WriteByte(c)
				if u == 0 {
					break
				}
			}
		}
	default:
		panic("encTyped: unsupported type")
	}
}

func encParamTable(params map[byte]any) []byte {
	var b bytes.Buffer
	b.WriteByte(byte(len(params)))
	// deterministic order not required by parser
	for k, v := range params {
		b.WriteByte(k)
		encTyped(&b, v)
	}
	return b.Bytes()
}

func encEventBody(code byte, params map[byte]any) []byte {
	var b bytes.Buffer
	b.WriteByte(code)
	b.Write(encParamTable(params))
	return b.Bytes()
}

func encResponseBody(code byte, returnCode int16, params map[byte]any) []byte {
	var b bytes.Buffer
	b.WriteByte(code)
	_ = binary.Write(&b, binary.LittleEndian, returnCode)
	b.WriteByte(p18Null) // debug message slot
	b.Write(encParamTable(params))
	return b.Bytes()
}

// wrapReliable builds a full Photon packet with one reliable command.
func wrapReliable(messageType byte, message []byte) []byte {
	var body bytes.Buffer
	body.WriteByte(0) // signalByte
	body.WriteByte(messageType)
	body.Write(message)

	var cmd bytes.Buffer
	cmd.WriteByte(6) // SendReliable
	cmd.WriteByte(0) // channelId
	cmd.WriteByte(0) // commandFlags
	cmd.WriteByte(0) // reserved
	cl := int32(commandHeaderLength + body.Len())
	_ = binary.Write(&cmd, binary.BigEndian, cl)
	_ = binary.Write(&cmd, binary.BigEndian, int32(0)) // seq
	cmd.Write(body.Bytes())

	return wrapPacket(1, cmd.Bytes())
}

func wrapPacket(commandCount byte, commands []byte) []byte {
	var p bytes.Buffer
	_ = binary.Write(&p, binary.BigEndian, int16(0)) // peerId
	p.WriteByte(0)                                   // flags
	p.WriteByte(commandCount)
	_ = binary.Write(&p, binary.BigEndian, int32(0)) // timestamp
	_ = binary.Write(&p, binary.BigEndian, int32(0)) // challenge
	p.Write(commands)
	return p.Bytes()
}

// ---- test handler ----

type capture struct {
	events    []EventPacket
	responses []ResponsePacket
	requests  []RequestPacket
}

type EventPacket struct {
	Code   EventCode
	Params map[byte]any
}
type ResponsePacket struct {
	Code   OperationCode
	Params map[byte]any
}
type RequestPacket struct {
	Code   OperationCode
	Params map[byte]any
}

func (c *capture) HandleEvent(code EventCode, p map[byte]any) {
	c.events = append(c.events, EventPacket{code, p})
}
func (c *capture) HandleRequest(code OperationCode, p map[byte]any) {
	c.requests = append(c.requests, RequestPacket{code, p})
}
func (c *capture) HandleResponse(code OperationCode, rc int16, dbg string, p map[byte]any) {
	c.responses = append(c.responses, ResponsePacket{code, p})
}

// ---- value round-trip ----

func TestValueRoundTrip(t *testing.T) {
	cases := []any{
		int16(-1234),
		int32(123456),
		int32(-7),
		int64(9876543210),
		byte(200),
		true,
		false,
		"hello world",
		"",
		[]string{"a", "bb", "ccc"},
		[]int32{1, -2, 3, -4},
	}
	for i, want := range cases {
		var b bytes.Buffer
		encTyped(&b, want)
		r := newReader(b.Bytes())
		got, err := r.deserialize()
		if err != nil {
			t.Fatalf("case %d: %v", i, err)
		}
		// normalize []int32 vs []any
		if exp, ok := want.([]int32); ok {
			arr, ok := got.([]any)
			if !ok || len(arr) != len(exp) {
				t.Fatalf("case %d: int32 array mismatch: %#v", i, got)
			}
			for j := range exp {
				if arr[j].(int32) != exp[j] {
					t.Fatalf("case %d[%d]: got %v want %v", i, j, arr[j], exp[j])
				}
			}
			continue
		}
		if exp, ok := want.([]string); ok {
			arr, ok := got.([]any)
			if !ok || len(arr) != len(exp) {
				t.Fatalf("case %d: string array mismatch: %#v", i, got)
			}
			for j := range exp {
				if arr[j].(string) != exp[j] {
					t.Fatalf("case %d[%d]: got %v want %v", i, j, arr[j], exp[j])
				}
			}
			continue
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("case %d: got %#v want %#v", i, got, want)
		}
		if r.remaining() != 0 {
			t.Fatalf("case %d: %d bytes left over", i, r.remaining())
		}
	}
}

func TestDictionaryRoundTrip(t *testing.T) {
	// typed dict<byte,string>
	var b bytes.Buffer
	b.WriteByte(p18Dictionary)
	b.WriteByte(p18Byte)   // key type
	b.WriteByte(p18String) // value type
	writeCount(&b, 2)
	b.WriteByte(1)
	writeCount(&b, 3)
	b.WriteString("one")
	b.WriteByte(2)
	writeCount(&b, 3)
	b.WriteString("two")

	r := newReader(b.Bytes())
	got, err := r.deserialize()
	if err != nil {
		t.Fatal(err)
	}
	m, ok := got.(map[any]any)
	if !ok {
		t.Fatalf("not a map: %#v", got)
	}
	if m[byte(1)] != "one" || m[byte(2)] != "two" {
		t.Fatalf("dict mismatch: %#v", m)
	}
}

func TestEventDispatch(t *testing.T) {
	c := &capture{}
	p := NewParser(c)
	// EstimatedMarketValueUpdate = 464, param 252 carries the event code
	params := map[byte]any{
		paramEventCode: int16(EvEstimatedMarketValueUpdate),
		0:              []int32{100, 200},
		1:              []int32{5000000, 9000000},
	}
	pkt := wrapReliable(4, encEventBody(7, params))
	if st := p.ReceivePacket(pkt); st != StatusSuccess {
		t.Fatalf("status %v", st)
	}
	if len(c.events) != 1 {
		t.Fatalf("got %d events", len(c.events))
	}
	if c.events[0].Code != EvEstimatedMarketValueUpdate {
		t.Fatalf("code %v", c.events[0].Code)
	}
}

func TestMarketResponseDispatch(t *testing.T) {
	c := &capture{}
	p := NewParser(c)
	order := `{"Id":1,"ItemTypeId":"T4_BAG","LocationId":"3005","QualityLevel":1,"UnitPriceSilver":12340000,"Amount":3,"AuctionType":"offer","Expires":"2026-07-01T00:00:00"}`
	params := map[byte]any{
		paramOperationCode: int16(OpAuctionGetOffers),
		0:                  []string{order},
	}
	pkt := wrapReliable(3, encResponseBody(byte(OpAuctionGetOffers), 0, params))
	if st := p.ReceivePacket(pkt); st != StatusSuccess {
		t.Fatalf("status %v", st)
	}
	if len(c.responses) != 1 {
		t.Fatalf("got %d responses", len(c.responses))
	}
	if c.responses[0].Code != OpAuctionGetOffers {
		t.Fatalf("code %v", c.responses[0].Code)
	}
	arr, ok := c.responses[0].Params[0].([]any)
	if !ok || len(arr) != 1 || arr[0].(string) != order {
		t.Fatalf("order param mismatch: %#v", c.responses[0].Params[0])
	}
}

func TestEncryptedMessage(t *testing.T) {
	c := &capture{}
	p := NewParser(c)
	pkt := wrapReliable(131, []byte{0x01, 0x02, 0x03})
	if st := p.ReceivePacket(pkt); st != StatusEncrypted {
		t.Fatalf("expected Encrypted, got %v", st)
	}
}

func TestFragmentReassembly(t *testing.T) {
	c := &capture{}
	p := NewParser(c)
	// build a full reliable body, then split into 2 fragments
	params := map[byte]any{
		paramEventCode: int16(EvEstimatedMarketValueUpdate),
		0:              []int32{1},
		1:              []int32{2},
	}
	var body bytes.Buffer
	body.WriteByte(0) // signal
	body.WriteByte(4) // event
	body.Write(encEventBody(7, params))
	full := body.Bytes()
	mid := len(full) / 2

	makeFrag := func(startSeq int32, fragOffset int, chunk []byte) []byte {
		var f bytes.Buffer
		_ = binary.Write(&f, binary.BigEndian, startSeq)
		_ = binary.Write(&f, binary.BigEndian, int32(2))            // fragmentCount
		_ = binary.Write(&f, binary.BigEndian, int32(0))            // fragmentNumber
		_ = binary.Write(&f, binary.BigEndian, int32(len(full)))    // totalLength
		_ = binary.Write(&f, binary.BigEndian, int32(fragOffset))   // fragmentOffset
		f.Write(chunk)
		var cmd bytes.Buffer
		cmd.WriteByte(8) // SendFragment
		cmd.WriteByte(0)
		cmd.WriteByte(0)
		cmd.WriteByte(0)
		_ = binary.Write(&cmd, binary.BigEndian, int32(commandHeaderLength+f.Len()))
		_ = binary.Write(&cmd, binary.BigEndian, int32(0))
		cmd.Write(f.Bytes())
		return cmd.Bytes()
	}

	p.ReceivePacket(wrapPacket(1, makeFrag(10, 0, full[:mid])))
	if len(c.events) != 0 {
		t.Fatalf("event dispatched before reassembly complete")
	}
	p.ReceivePacket(wrapPacket(1, makeFrag(10, mid, full[mid:])))
	if len(c.events) != 1 || c.events[0].Code != EvEstimatedMarketValueUpdate {
		t.Fatalf("fragment reassembly failed: %#v", c.events)
	}
}

var _ = writeZigzag32
