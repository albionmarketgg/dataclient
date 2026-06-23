// Package phototest builds synthetic Photon packets for tests and replay tools.
package phototest

import (
	"bytes"
	"encoding/binary"
)

// type codes (subset, mirrors internal/photon).
const (
	tShort        byte = 4
	tString       byte = 7
	tNull         byte = 8
	tCompInt      byte = 9
	tCompLong     byte = 10
	tByte         byte = 3
	tBoolFalse    byte = 27
	tBoolTrue     byte = 28
	tStringArray  byte = 71
	tCompIntArray byte = 73
	tCompLongArr  byte = 74
)

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

func zig32(b *bytes.Buffer, v int32) {
	u := uint32((v << 1) ^ (v >> 31))
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

func zig64(b *bytes.Buffer, v int64) {
	u := uint64((v << 1) ^ (v >> 63))
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

// EncTyped writes a typed value [typeByte][payload].
func EncTyped(b *bytes.Buffer, v any) {
	switch x := v.(type) {
	case int16:
		b.WriteByte(tShort)
		binary.Write(b, binary.LittleEndian, x)
	case int32:
		b.WriteByte(tCompInt)
		zig32(b, x)
	case int64:
		b.WriteByte(tCompLong)
		zig64(b, x)
	case byte:
		b.WriteByte(tByte)
		b.WriteByte(x)
	case bool:
		if x {
			b.WriteByte(tBoolTrue)
		} else {
			b.WriteByte(tBoolFalse)
		}
	case string:
		b.WriteByte(tString)
		writeCount(b, len(x))
		b.WriteString(x)
	case nil:
		b.WriteByte(tNull)
	case []string:
		b.WriteByte(tStringArray)
		writeCount(b, len(x))
		for _, s := range x {
			writeCount(b, len(s))
			b.WriteString(s)
		}
	case []int32:
		b.WriteByte(tCompIntArray)
		writeCount(b, len(x))
		for _, n := range x {
			zig32(b, n)
		}
	case []int64:
		b.WriteByte(tCompLongArr)
		writeCount(b, len(x))
		for _, n := range x {
			zig64(b, n)
		}
	default:
		panic("phototest.EncTyped: unsupported type")
	}
}

func paramTable(params map[byte]any) []byte {
	var b bytes.Buffer
	b.WriteByte(byte(len(params)))
	for k, v := range params {
		b.WriteByte(k)
		EncTyped(&b, v)
	}
	return b.Bytes()
}

func wrapPacket(commands []byte) []byte {
	var p bytes.Buffer
	binary.Write(&p, binary.BigEndian, int16(0)) // peerId
	p.WriteByte(0)                               // flags
	p.WriteByte(1)                               // commandCount
	binary.Write(&p, binary.BigEndian, int32(0)) // timestamp
	binary.Write(&p, binary.BigEndian, int32(0)) // challenge
	p.Write(commands)
	return p.Bytes()
}

func reliableCommand(messageType byte, message []byte) []byte {
	var body bytes.Buffer
	body.WriteByte(0) // signal
	body.WriteByte(messageType)
	body.Write(message)
	var cmd bytes.Buffer
	cmd.WriteByte(6) // reliable
	cmd.WriteByte(0)
	cmd.WriteByte(0)
	cmd.WriteByte(0)
	binary.Write(&cmd, binary.BigEndian, int32(12+body.Len()))
	binary.Write(&cmd, binary.BigEndian, int32(0))
	cmd.Write(body.Bytes())
	return cmd.Bytes()
}

// ResponsePacket builds a Photon packet carrying one OperationResponse. The
// params must include key 253 with the operation code (as int16) for routing.
func ResponsePacket(respCodeByte byte, returnCode int16, params map[byte]any) []byte {
	var msg bytes.Buffer
	msg.WriteByte(respCodeByte)
	binary.Write(&msg, binary.LittleEndian, returnCode)
	msg.WriteByte(tNull) // debug slot
	msg.Write(paramTable(params))
	return wrapPacket(reliableCommand(3, msg.Bytes()))
}

// EventPacket builds a Photon packet carrying one Event. The params must include
// key 252 with the event code (as int16) for routing.
func EventPacket(eventCodeByte byte, params map[byte]any) []byte {
	var msg bytes.Buffer
	msg.WriteByte(eventCodeByte)
	msg.Write(paramTable(params))
	return wrapPacket(reliableCommand(4, msg.Bytes()))
}

// RequestPacket builds a Photon packet carrying one OperationRequest. The params
// must include key 253 with the operation code (as int16) for routing.
func RequestPacket(reqCodeByte byte, params map[byte]any) []byte {
	var msg bytes.Buffer
	msg.WriteByte(reqCodeByte)
	msg.Write(paramTable(params))
	return wrapPacket(reliableCommand(2, msg.Bytes()))
}
