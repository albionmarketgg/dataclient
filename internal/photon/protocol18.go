package photon

import "fmt"

// Protocol18 value type codes.
const (
	p18Unknown            byte = 0
	p18Boolean            byte = 2
	p18Byte               byte = 3
	p18Short              byte = 4
	p18Float              byte = 5
	p18Double             byte = 6
	p18String             byte = 7
	p18Null               byte = 8
	p18CompressedInt      byte = 9
	p18CompressedLong     byte = 10
	p18Int1               byte = 11
	p18Int1Negative       byte = 12
	p18Int2               byte = 13
	p18Int2Negative       byte = 14
	p18Long1              byte = 15
	p18Long1Negative      byte = 16
	p18Long2              byte = 17
	p18Long2Negative      byte = 18
	p18Custom             byte = 19
	p18Dictionary         byte = 20
	p18Hashtable          byte = 21
	p18ObjectArray        byte = 23
	p18OperationRequest   byte = 24
	p18OperationResponse  byte = 25
	p18EventData          byte = 26
	p18BooleanFalse       byte = 27
	p18BooleanTrue        byte = 28
	p18ShortZero          byte = 29
	p18IntZero            byte = 30
	p18LongZero           byte = 31
	p18FloatZero          byte = 32
	p18DoubleZero         byte = 33
	p18ByteZero           byte = 34
	p18Array              byte = 64
	p18BooleanArray       byte = 66
	p18ByteArray          byte = 67
	p18ShortArray         byte = 68
	p18FloatArray         byte = 69
	p18DoubleArray        byte = 70
	p18StringArray        byte = 71
	p18CompressedIntArray byte = 73
	p18CompressedLongArr  byte = 74
	p18CustomTypeArray    byte = 83
	p18DictionaryArray    byte = 84
	p18HashtableArray     byte = 85

	p18SlimCustomBase byte = 0x80
	p18SlimCustomMax  byte = 0xE4
)

// CustomData holds an undecoded Protocol18 custom type payload.
type CustomData struct {
	TypeCode byte
	Data     []byte
}

// deserialize reads a self-describing value (type byte then payload).
func (r *reader) deserialize() (any, error) {
	tc, err := r.readByte()
	if err != nil {
		return nil, err
	}
	return r.deserializeWithType(tc)
}

func (r *reader) deserializeWithType(tc byte) (any, error) {
	if tc >= p18SlimCustomBase && tc <= p18SlimCustomMax {
		return r.readCustomPayload(tc - p18SlimCustomBase)
	}
	switch tc {
	case p18Unknown, p18Null:
		return nil, nil
	case p18Boolean:
		b, err := r.readByte()
		return b != 0, err
	case p18Byte:
		return r.readByte()
	case p18Short:
		return r.readInt16LE()
	case p18Float:
		return r.readFloat32()
	case p18Double:
		return r.readFloat64()
	case p18String:
		return r.readStringValue()
	case p18CompressedInt:
		return r.readCompressedInt32()
	case p18CompressedLong:
		return r.readCompressedInt64()
	case p18Int1:
		b, err := r.readByte()
		return int32(b), err
	case p18Int1Negative:
		b, err := r.readByte()
		return -int32(b), err
	case p18Int2:
		v, err := r.readUint16LE()
		return int32(v), err
	case p18Int2Negative:
		v, err := r.readUint16LE()
		return -int32(v), err
	case p18Long1:
		b, err := r.readByte()
		return int64(b), err
	case p18Long1Negative:
		b, err := r.readByte()
		return -int64(b), err
	case p18Long2:
		v, err := r.readUint16LE()
		return int64(v), err
	case p18Long2Negative:
		v, err := r.readUint16LE()
		return -int64(v), err
	case p18Custom:
		ct, err := r.readByte()
		if err != nil {
			return nil, err
		}
		return r.readCustomPayload(ct)
	case p18Dictionary:
		return r.readDictionary()
	case p18Hashtable:
		return r.readHashtable()
	case p18ObjectArray:
		return r.readObjectArray()
	case p18OperationRequest:
		return r.readOperationRequest()
	case p18OperationResponse:
		return r.readOperationResponse()
	case p18EventData:
		return r.readEventData()
	case p18BooleanFalse:
		return false, nil
	case p18BooleanTrue:
		return true, nil
	case p18ShortZero:
		return int16(0), nil
	case p18IntZero:
		return int32(0), nil
	case p18LongZero:
		return int64(0), nil
	case p18FloatZero:
		return float32(0), nil
	case p18DoubleZero:
		return float64(0), nil
	case p18ByteZero:
		return byte(0), nil
	case p18Array:
		return r.readArrayInArray()
	case p18BooleanArray:
		return r.readBooleanArray()
	case p18ByteArray:
		n, err := r.readCount()
		if err != nil {
			return nil, err
		}
		b, err := r.readBytes(n)
		if err != nil {
			return nil, err
		}
		out := make([]byte, n)
		copy(out, b)
		return out, nil
	case p18ShortArray:
		return readTypedArray(r, func(r *reader) (any, error) { return r.readInt16LE() })
	case p18FloatArray:
		return readTypedArray(r, func(r *reader) (any, error) { return r.readFloat32() })
	case p18DoubleArray:
		return readTypedArray(r, func(r *reader) (any, error) { return r.readFloat64() })
	case p18StringArray:
		return readTypedArray(r, func(r *reader) (any, error) { return r.readStringValue() })
	case p18CompressedIntArray:
		return readTypedArray(r, func(r *reader) (any, error) { return r.readCompressedInt32() })
	case p18CompressedLongArr:
		return readTypedArray(r, func(r *reader) (any, error) { return r.readCompressedInt64() })
	case p18CustomTypeArray:
		return r.readCustomTypeArray()
	case p18DictionaryArray:
		return r.readDictionaryArray()
	case p18HashtableArray:
		return readTypedArray(r, func(r *reader) (any, error) { return r.readHashtable() })
	default:
		return nil, fmt.Errorf("photon: type code %d not implemented", tc)
	}
}

func readTypedArray(r *reader, read func(*reader) (any, error)) ([]any, error) {
	n, err := r.readCount()
	if err != nil {
		return nil, err
	}
	out := make([]any, n)
	for i := 0; i < n; i++ {
		v, err := read(r)
		if err != nil {
			return nil, err
		}
		out[i] = v
	}
	return out, nil
}

func (r *reader) readCustomPayload(typeCode byte) (CustomData, error) {
	n, err := r.readCount()
	if err != nil {
		return CustomData{}, err
	}
	b, err := r.readBytes(n)
	if err != nil {
		return CustomData{}, err
	}
	data := make([]byte, n)
	copy(data, b)
	return CustomData{TypeCode: typeCode, Data: data}, nil
}

func (r *reader) readCustomTypeArray() ([]CustomData, error) {
	n, err := r.readCount()
	if err != nil {
		return nil, err
	}
	tc, err := r.readByte()
	if err != nil {
		return nil, err
	}
	out := make([]CustomData, n)
	for i := 0; i < n; i++ {
		cd, err := r.readCustomPayload(tc)
		if err != nil {
			return nil, err
		}
		out[i] = cd
	}
	return out, nil
}

func (r *reader) readBooleanArray() ([]bool, error) {
	n, err := r.readCount()
	if err != nil {
		return nil, err
	}
	out := make([]bool, n)
	full := n / 8
	idx := 0
	for i := 0; i < full; i++ {
		b, err := r.readByte()
		if err != nil {
			return nil, err
		}
		for bit := 0; bit < 8; bit++ {
			out[idx] = b&(1<<uint(bit)) != 0
			idx++
		}
	}
	if n%8 != 0 {
		b, err := r.readByte()
		if err != nil {
			return nil, err
		}
		for bit := 0; idx < n; bit++ {
			out[idx] = b&(1<<uint(bit)) != 0
			idx++
		}
	}
	return out, nil
}

func (r *reader) readObjectArray() ([]any, error) {
	n, err := r.readCount()
	if err != nil {
		return nil, err
	}
	out := make([]any, n)
	for i := 0; i < n; i++ {
		v, err := r.deserialize()
		if err != nil {
			return nil, err
		}
		out[i] = v
	}
	return out, nil
}

// readArrayInArray (type 64) consumes ReadCount self-describing values.
func (r *reader) readArrayInArray() ([]any, error) {
	n, err := r.readCount()
	if err != nil {
		return nil, err
	}
	out := make([]any, n)
	for i := 0; i < n; i++ {
		v, err := r.deserialize()
		if err != nil {
			return nil, err
		}
		out[i] = v
	}
	return out, nil
}

func (r *reader) readHashtable() (map[any]any, error) {
	n, err := r.readCount()
	if err != nil {
		return nil, err
	}
	out := make(map[any]any, n)
	for i := 0; i < n; i++ {
		k, err := r.deserialize()
		if err != nil {
			return nil, err
		}
		v, err := r.deserialize()
		if err != nil {
			return nil, err
		}
		if k != nil {
			out[k] = v
		}
	}
	return out, nil
}

// readDictionary reads a typed dictionary (type 20).
func (r *reader) readDictionary() (map[any]any, error) {
	keyType, valType, err := r.readDictionaryHeader()
	if err != nil {
		return nil, err
	}
	return r.readDictionaryElements(keyType, valType)
}

func (r *reader) readDictionaryHeader() (keyType, valType byte, err error) {
	keyType, err = r.readByte()
	if err != nil {
		return
	}
	valType, err = r.readByte()
	if err != nil {
		return
	}
	// Resolve composite value types so element reading consumes correctly.
	switch valType {
	case p18Array:
		// consume the array element-type chain; elements become self-describing
		t, e := r.readByte()
		if e != nil {
			return 0, 0, e
		}
		for t == p18Array {
			t, e = r.readByte()
			if e != nil {
				return 0, 0, e
			}
		}
		valType = p18Unknown
	case p18Dictionary:
		// nested dictionary header consumed lazily during element read
	}
	return keyType, valType, nil
}

func (r *reader) readDictionaryElements(keyType, valType byte) (map[any]any, error) {
	n, err := r.readCount()
	if err != nil {
		return nil, err
	}
	out := make(map[any]any, n)
	for i := 0; i < n; i++ {
		var k any
		if keyType == p18Unknown {
			k, err = r.deserialize()
		} else {
			k, err = r.deserializeWithType(keyType)
		}
		if err != nil {
			return nil, err
		}
		var v any
		if valType == p18Unknown {
			v, err = r.deserialize()
		} else {
			v, err = r.deserializeWithType(valType)
		}
		if err != nil {
			return nil, err
		}
		if k != nil {
			out[k] = v
		}
	}
	return out, nil
}

func (r *reader) readDictionaryArray() ([]map[any]any, error) {
	keyType, valType, err := r.readDictionaryHeader()
	if err != nil {
		return nil, err
	}
	n, err := r.readCount()
	if err != nil {
		return nil, err
	}
	out := make([]map[any]any, n)
	for i := 0; i < n; i++ {
		d, err := r.readDictionaryElements(keyType, valType)
		if err != nil {
			return nil, err
		}
		out[i] = d
	}
	return out, nil
}

// parameterTable reads the message parameter table: 1-byte count, then count
// (key byte, value-type byte, value) triples.
func (r *reader) parameterTable() (map[byte]any, error) {
	count, err := r.readByte()
	if err != nil {
		return nil, err
	}
	out := make(map[byte]any, count)
	for i := 0; i < int(count); i++ {
		key, err := r.readByte()
		if err != nil {
			return nil, err
		}
		tc, err := r.readByte()
		if err != nil {
			return nil, err
		}
		v, err := r.deserializeWithType(tc)
		if err != nil {
			return nil, err
		}
		out[key] = v
	}
	return out, nil
}
