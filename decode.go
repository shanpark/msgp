package msgp

import (
	"io"
	"reflect"
)

// UnpackValue read a msgpack object from reader to ptr.
func UnpackValue(r io.Reader, ptr interface{}) error {
	wantTyp := reflect.TypeOf(ptr).Elem()

	switch wantTyp.Kind() {
	case reflect.Ptr:
		newVal := reflect.New(wantTyp.Elem())
		UnpackValue(r, newVal.Interface())
		// newVal.Elem().SetInt(newVal.Int()) // 위의 Unpack 대신이다.
		reflect.ValueOf(ptr).Elem().Set(newVal)
	}
	return nil
}

// // UnpackValue reads a value from reader.
// func UnpackPrimitive(r io.Reader) (interface{}, error) {
// 	var head byte

// 	if err := binary.Read(r, binary.BigEndian, &head); err != nil {
// 		return nil, err
// 	}

// 	if head == 0xc0 { // nil
// 		return nil, nil
// 	} else if head == 0xc2 { // bool
// 		return false, nil
// 	} else if head == 0xc3 {
// 		return true, nil
// 	} else if head&0x80 == 0 { // int8
// 		return int8(head), nil
// 	} else if head&0xe0 == 0xe0 {
// 		return int8(head), nil
// 	} else if head == 0xd0 {
// 		return UnpackInt8(r)
// 	} else if head == 0xd1 {
// 		return UnpackInt16(r)
// 	} else if head == 0xd2 {
// 		return UnpackInt32(r)
// 	} else if head == 0xd3 {
// 		return UnpackInt64(r)
// 	} else if head == 0xcc {
// 		return UnpackUint8(r)
// 	} else if head == 0xcd {
// 		return UnpackUint16(r)
// 	} else if head == 0xce {
// 		return UnpackUint32(r)
// 	} else if head == 0xcf {
// 		return UnpackUint64(r)
// 	} else if head == 0xca {
// 		return UnpackFloat32(r)
// 	} else if head == 0xcb {
// 		return UnpackFloat64(r)
// 	} else if head&0xa0 == 0xa0 {
// 		return UnpackString(r, head)
// 	} else if head == 0xd9 {
// 		return UnpackString(r, head)
// 	} else if head == 0xda {
// 		return UnpackString(r, head)
// 	} else if head == 0xdb {
// 		return UnpackString(r, head)
// 	}
// }

// // UnpackInt8 reads a int8 value from reader.
// func UnpackInt8(r io.Reader) (int8, error) {
// 	var data int8
// 	err := binary.Read(r, binary.BigEndian, &data)
// 	return data, err
// }

// // UnpackInt16 reads a int8 value from reader.
// func UnpackInt16(r io.Reader) (int16, error) {
// 	var data int16
// 	err := binary.Read(r, binary.BigEndian, &data)
// 	return data, err
// }

// // UnpackInt32 reads a int8 value from reader.
// func UnpackInt32(r io.Reader) (int32, error) {
// 	var data int32
// 	err := binary.Read(r, binary.BigEndian, &data)
// 	return data, err
// }

// // UnpackInt64 reads a int8 value from reader.
// func UnpackInt64(r io.Reader) (int64, error) {
// 	var data int64
// 	err := binary.Read(r, binary.BigEndian, &data)
// 	return data, err
// }

// // UnpackUint8 reads a uint8 value from reader.
// func UnpackUint8(r io.Reader) (uint8, error) {
// 	var data uint8
// 	err := binary.Read(r, binary.BigEndian, &data)
// 	return data, err
// }

// // UnpackUint16 reads a uint16 value from reader.
// func UnpackUint16(r io.Reader) (uint16, error) {
// 	var data uint16
// 	err := binary.Read(r, binary.BigEndian, &data)
// 	return data, err
// }

// // UnpackUint32 reads a uint32 value from reader.
// func UnpackUint32(r io.Reader) (uint32, error) {
// 	var data uint32
// 	err := binary.Read(r, binary.BigEndian, &data)
// 	return data, err
// }

// // UnpackUint64 reads a uint64 value from reader.
// func UnpackUint64(r io.Reader) (uint64, error) {
// 	var data uint64
// 	err := binary.Read(r, binary.BigEndian, &data)
// 	return data, err
// }

// // UnpackFloat32 reads a float32 value from reader.
// func UnpackFloat32(r io.Reader) (float32, error) {
// 	buf := make([]byte, 4)
// 	if _, err := r.Read(buf); err != nil {
// 		return 0, err
// 	}

// 	bits := binary.LittleEndian.Uint32(buf)
// 	return math.Float32frombits(bits), nil
// }

// // UnpackFloat64 reads a float64 value from reader.
// func UnpackFloat64(r io.Reader) (float64, error) {
// 	buf := make([]byte, 8)
// 	if _, err := r.Read(buf); err != nil {
// 		return 0, err
// 	}

// 	bits := binary.LittleEndian.Uint64(buf)
// 	return math.Float64frombits(bits), nil
// }

// func UnpackString(r io.Reader, head uint8) {
// 	var len uint64
// 	if head&0xa0 == 0xa0 {
// 		len = uint64(head & 0x1f)
// 	} else if head == 0xd9 {
// 		var temp uint8
// 		binary.Read(r, binary.BigEndian, &temp)
// 		len = uint64(temp)
// 	} else if head == 0xda {
// 		var temp uint16
// 		binary.Read(r, binary.BigEndian, &temp)
// 		len = uint64(temp)
// 	} else if head == 0xdb {
// 		var temp uint32
// 		binary.Read(r, binary.BigEndian, &temp)
// 		len = uint64(temp)
// 	}
// }
