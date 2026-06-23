package photon

import (
	"encoding/binary"
	"errors"
	"math"
	"strconv"
)

func itoa(i int) string { return strconv.Itoa(i) }

var errEOF = errors.New("photon: unexpected end of buffer")

// reader is a byte-slice cursor supporting both big-endian Photon framing and
// little-endian / varint Protocol18 value primitives.
type reader struct {
	buf []byte
	pos int
}

func newReader(b []byte) *reader { return &reader{buf: b} }

func (r *reader) remaining() int { return len(r.buf) - r.pos }
func (r *reader) has(n int) bool { return n >= 0 && r.remaining() >= n }

func (r *reader) readByte() (byte, error) {
	if r.remaining() < 1 {
		return 0, errEOF
	}
	b := r.buf[r.pos]
	r.pos++
	return b, nil
}

func (r *reader) readBytes(n int) ([]byte, error) {
	if n < 0 || r.remaining() < n {
		return nil, errEOF
	}
	b := r.buf[r.pos : r.pos+n]
	r.pos += n
	return b, nil
}

// ---- big-endian (Photon framing) ----

func (r *reader) readInt16BE() (int16, error) {
	b, err := r.readBytes(2)
	if err != nil {
		return 0, err
	}
	return int16(binary.BigEndian.Uint16(b)), nil
}

func (r *reader) readInt32BE() (int32, error) {
	b, err := r.readBytes(4)
	if err != nil {
		return 0, err
	}
	return int32(binary.BigEndian.Uint32(b)), nil
}

// ---- little-endian (Protocol18 values) ----

func (r *reader) readInt16LE() (int16, error) {
	b, err := r.readBytes(2)
	if err != nil {
		return 0, err
	}
	return int16(binary.LittleEndian.Uint16(b)), nil
}

func (r *reader) readUint16LE() (uint16, error) {
	b, err := r.readBytes(2)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(b), nil
}

func (r *reader) readInt32LE() (int32, error) {
	b, err := r.readBytes(4)
	if err != nil {
		return 0, err
	}
	return int32(binary.LittleEndian.Uint32(b)), nil
}

func (r *reader) readFloat32() (float32, error) {
	b, err := r.readBytes(4)
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(binary.LittleEndian.Uint32(b)), nil
}

func (r *reader) readFloat64() (float64, error) {
	b, err := r.readBytes(8)
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(binary.LittleEndian.Uint64(b)), nil
}

// ---- varints (LEB128, low byte first) ----

func (r *reader) readCompressedUint32() (uint32, error) {
	var value uint32
	var shift uint
	for {
		if shift >= 35 {
			return 0, errors.New("photon: compressed uint32 too long")
		}
		b, err := r.readByte()
		if err != nil {
			return 0, err
		}
		value |= uint32(b&0x7f) << shift
		if b&0x80 == 0 {
			return value, nil
		}
		shift += 7
	}
}

func (r *reader) readCompressedUint64() (uint64, error) {
	var value uint64
	var shift uint
	for {
		if shift >= 70 {
			return 0, errors.New("photon: compressed uint64 too long")
		}
		b, err := r.readByte()
		if err != nil {
			return 0, err
		}
		value |= uint64(b&0x7f) << shift
		if b&0x80 == 0 {
			return value, nil
		}
		shift += 7
	}
}

func (r *reader) readCompressedInt32() (int32, error) {
	u, err := r.readCompressedUint32()
	if err != nil {
		return 0, err
	}
	return int32(u>>1) ^ -int32(u&1), nil
}

func (r *reader) readCompressedInt64() (int64, error) {
	u, err := r.readCompressedUint64()
	if err != nil {
		return 0, err
	}
	return int64(u>>1) ^ -int64(u&1), nil
}

// readCount is the varint length/size prefix used for all Protocol18 collections.
func (r *reader) readCount() (int, error) {
	u, err := r.readCompressedUint32()
	if err != nil {
		return 0, err
	}
	return int(u), nil
}

func (r *reader) readStringValue() (string, error) {
	n, err := r.readCount()
	if err != nil {
		return "", err
	}
	if n == 0 {
		return "", nil
	}
	b, err := r.readBytes(n)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
