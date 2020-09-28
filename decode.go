package msgp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
)

func UnpackValue(r io.Reader, ptr interface{}) error {
	var err error

	wantType := reflect.TypeOf(ptr).Elem()
	switch wantType.Kind() {
	case reflect.Bool:
		err = UnpackBool(r, ptr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		err = UnpackInt(r, ptr)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		err = UnpackUint(r, ptr)
	case reflect.Float32, reflect.Float64:
		err = UnpackFloat(r, ptr)
	case reflect.String:
		err = UnpackString(r, ptr)
	case reflect.Array:
		err = UnpackArray(r, ptr)
	case reflect.Slice:
		err = UnpackSlice(r, ptr)
	case reflect.Map:
		err = UnpackMap(r, ptr)
	case reflect.Ptr:
		err = UnpackPtr(r, ptr)
	default:
		return fmt.Errorf("msgp: specified type[%v] is not supported", wantType.Kind())
	}

	return err
}

func UnpackBool(r io.Reader, ptr interface{}) error {
	var err error
	var val interface{}

	if val, err = UnpackPrimitive(r); err != nil {
		return err
	}

	reflect.ValueOf(ptr).Elem().SetBool(val.(bool))
	return nil
}

func UnpackInt(r io.Reader, ptr interface{}) error {
	var err error
	var val interface{}

	if val, err = UnpackPrimitive(r); err != nil {
		return err
	}

	switch reflect.ValueOf(val).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		reflect.ValueOf(ptr).Elem().SetInt(reflect.ValueOf(val).Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		reflect.ValueOf(ptr).Elem().SetInt(int64(reflect.ValueOf(val).Uint()))
	case reflect.Float32, reflect.Float64:
		reflect.ValueOf(ptr).Elem().SetInt(int64(reflect.ValueOf(val).Float()))
	default:
		return fmt.Errorf("msgp: unpacked value[%v] is not assignable to integer type", val)
	}
	return nil
}

func UnpackUint(r io.Reader, ptr interface{}) error {
	var err error
	var val interface{}

	if val, err = UnpackPrimitive(r); err != nil {
		return err
	}

	switch reflect.ValueOf(val).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		reflect.ValueOf(ptr).Elem().SetUint(uint64(reflect.ValueOf(val).Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		reflect.ValueOf(ptr).Elem().SetUint(reflect.ValueOf(val).Uint())
	case reflect.Float32, reflect.Float64:
		reflect.ValueOf(ptr).Elem().SetUint(uint64(reflect.ValueOf(val).Float()))
	default:
		return fmt.Errorf("msgp: unpacked value[%v] is not assignable to unsigned integer type", val)
	}
	return nil
}

func UnpackFloat(r io.Reader, ptr interface{}) error {
	var err error
	var val interface{}

	if val, err = UnpackPrimitive(r); err != nil {
		return err
	}

	switch reflect.ValueOf(val).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		reflect.ValueOf(ptr).Elem().SetFloat(float64(reflect.ValueOf(val).Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		reflect.ValueOf(ptr).Elem().SetFloat(float64(reflect.ValueOf(val).Uint()))
	case reflect.Float32, reflect.Float64:
		reflect.ValueOf(ptr).Elem().SetFloat(reflect.ValueOf(val).Float())
	default:
		return fmt.Errorf("msgp: unpacked value[%v] is not assignable to float type", val)
	}
	return nil
}

func UnpackString(r io.Reader, ptr interface{}) error {
	var err error
	var val interface{}

	if val, err = UnpackPrimitive(r); err != nil {
		return err
	}

	if reflect.ValueOf(val).Kind() == reflect.String {
		reflect.ValueOf(ptr).Elem().SetString(reflect.ValueOf(val).String())
	} else {
		return fmt.Errorf("msgp: unpacked value[%v] is not assignable to string type", val)
	}
	return nil
}

func UnpackArray(r io.Reader, ptr interface{}) error {
	var err error
	var head byte

	arrTyp := reflect.TypeOf(ptr).Elem()
	arrVal := reflect.ValueOf(ptr).Elem()
	arrLen := arrVal.Len()

	if err = binary.Read(r, binary.BigEndian, &head); err != nil {
		return err
	}

	// handle bin format family
	if head == 0xc4 || head == 0xc5 || head == 0xc6 {
		if arrTyp.Elem().Kind() == reflect.Uint8 {
			var byteSlice []byte

			switch head {
			case 0xc4:
				byteSlice, err = UnpackBin8(r)
			case 0xc5:
				byteSlice, err = UnpackBin16(r)
			case 0xc6:
				byteSlice, err = UnpackBin32(r)
			}
			if err != nil {
				return err
			}

			arrVal.Set(reflect.Zero(reflect.ArrayOf(len(byteSlice), reflect.TypeOf(byteSlice).Elem())))
			reflect.Copy(arrVal, reflect.ValueOf(byteSlice))
			return nil
		} else {
			return fmt.Errorf("msgp: byte array can't be assigned to other type[%v] array", arrTyp.Elem().Kind())
		}
	}

	// handle array format family
	var srcLen = 0
	if head&0xf0 == 0x90 { // array
		srcLen = int(head & 0x0f)
	} else if head == 0xdc {
		var temp uint16
		if err = binary.Read(r, binary.BigEndian, &temp); err != nil {
			return err
		}
		srcLen = int(temp)
	} else if head == 0xdd {
		var temp uint32
		if err = binary.Read(r, binary.BigEndian, &temp); err != nil {
			return err
		}
		srcLen = int(temp) // maybe overflow.
	} else {
		return fmt.Errorf("msgp: unpacked value is not an array")
	}

	if arrLen < srcLen {
		return fmt.Errorf("msgp: array size is too small")
	}

	arrVal.Set(reflect.Zero(reflect.ArrayOf(arrLen, arrTyp.Elem()))) // array 생성.
	for inx := 0; inx < srcLen; inx++ {
		if err = UnpackValue(r, arrVal.Index(inx).Addr()); err != nil {
			return err
		}
	}
	return nil
}

func UnpackSlice(r io.Reader, ptr interface{}) error {
	var err error
	var head byte

	sliceTyp := reflect.TypeOf(ptr).Elem()
	sliceVal := reflect.ValueOf(ptr).Elem()

	if err = binary.Read(r, binary.BigEndian, &head); err != nil {
		return err
	}

	// handle bin format family
	if head == 0xc4 || head == 0xc5 || head == 0xc6 {
		if sliceTyp.Elem().Kind() == reflect.Uint8 {
			var byteSlice []byte

			switch head {
			case 0xc4:
				byteSlice, err = UnpackBin8(r)
			case 0xc5:
				byteSlice, err = UnpackBin16(r)
			case 0xc6:
				byteSlice, err = UnpackBin32(r)
			}
			if err != nil {
				return err
			}

			sliceVal.Set(reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(byteSlice).Elem()), len(byteSlice), len(byteSlice)))
			reflect.Copy(sliceVal, reflect.ValueOf(byteSlice))
			return nil
		} else {
			return fmt.Errorf("msgp: byte array can't be assigned to other type[%v] slice", sliceTyp.Elem().Kind())
		}
	}

	// handle array format family
	var srcLen = 0
	if head&0xf0 == 0x90 { // array
		srcLen = int(head & 0x0f)
	} else if head == 0xdc {
		var temp uint16
		if err = binary.Read(r, binary.BigEndian, &temp); err != nil {
			return err
		}
		srcLen = int(temp)
	} else if head == 0xdd {
		var temp uint32
		if err = binary.Read(r, binary.BigEndian, &temp); err != nil {
			return err
		}
		srcLen = int(temp) // maybe overflow.
	} else {
		return fmt.Errorf("msgp: unpacked value is not an array")
	}

	sliceVal.Set(reflect.MakeSlice(reflect.SliceOf(sliceTyp.Elem()), srcLen, srcLen)) // slice 생성.
	for inx := 0; inx < srcLen; inx++ {
		if err = UnpackValue(r, sliceVal.Index(inx).Addr()); err != nil {
			return err
		}
	}
	return nil
}

func UnpackMap(r io.Reader, ptr interface{}) error {
	var err error
	var head byte

	mapTyp := reflect.TypeOf(ptr).Elem()
	mapVal := reflect.ValueOf(ptr).Elem()

	if err = binary.Read(r, binary.BigEndian, &head); err != nil {
		return err
	}

	var srcLen = 0
	if head&0xf0 == 0x80 { // map
		srcLen = int(head & 0x0f)
	} else if head == 0xde {
		var temp uint16
		if err = binary.Read(r, binary.BigEndian, &temp); err != nil {
			return err
		}
		srcLen = int(temp)
	} else if head == 0xdf {
		var temp uint32
		if err = binary.Read(r, binary.BigEndian, &temp); err != nil {
			return err
		}
		srcLen = int(temp)
	} else {
		return fmt.Errorf("msgp: unpacked value is not a map")
	}

	mapVal.Set(reflect.MakeMap(reflect.MapOf(mapTyp.Key(), mapTyp.Elem()))) // map 생성.
	for inx := 0; inx < srcLen; inx++ {
		keyPtr := reflect.New(mapTyp.Key())
		if err = UnpackValue(r, keyPtr.Interface()); err != nil {
			return err
		}

		valPtr := reflect.New(mapTyp.Elem())
		if err = UnpackValue(r, valPtr.Interface()); err != nil {
			return err
		}
		mapVal.SetMapIndex(keyPtr.Elem(), valPtr.Elem())
	}
	return nil
}

func UnpackPtr(r io.Reader, ptr interface{}) error {
	var err error

	newVal := reflect.New(reflect.TypeOf(ptr).Elem().Elem())
	if err = UnpackValue(r, newVal.Interface()); err != nil {
		return err
	}

	reflect.ValueOf(ptr).Elem().Set(newVal)
	return nil
}

// UnpackValue reads a value from reader.
func UnpackPrimitive(r io.Reader) (interface{}, error) {
	var err error
	var head byte

	if err = binary.Read(r, binary.BigEndian, &head); err != nil {
		return nil, err
	}

	if head == 0xc0 { // nil
		return nil, nil
	} else if head == 0xc2 { // bool
		return false, nil
	} else if head == 0xc3 {
		return true, nil
	} else if head&0x80 == 0 { // int8
		return int8(head), nil
	} else if head&0xe0 == 0xe0 {
		return int8(head), nil
	} else if head == 0xd0 {
		return UnpackInt8(r)
	} else if head == 0xd1 {
		return UnpackInt16(r)
	} else if head == 0xd2 {
		return UnpackInt32(r)
	} else if head == 0xd3 {
		return UnpackInt64(r)
	} else if head == 0xcc {
		return UnpackUint8(r)
	} else if head == 0xcd {
		return UnpackUint16(r)
	} else if head == 0xce {
		return UnpackUint32(r)
	} else if head == 0xcf {
		return UnpackUint64(r)
	} else if head == 0xca {
		return UnpackFloat32(r)
	} else if head == 0xcb {
		return UnpackFloat64(r)
	} else if head&0xe0 == 0xa0 {
		return UnpackString5(r, int(head&0x1f))
	} else if head == 0xd9 {
		return UnpackString8(r)
	} else if head == 0xda {
		return UnpackString16(r)
	} else if head == 0xdb {
		return UnpackString32(r)
	} else if head == 0xc4 { // bin
		return UnpackBin8(r)
	} else if head == 0xc5 {
		return UnpackBin16(r)
	} else if head == 0xc6 {
		return UnpackBin32(r)
	}

	return nil, errors.New("msgp: Unknown or unsupported format family found.")
}

// UnpackInt8 reads a int8 value from reader.
func UnpackInt8(r io.Reader) (int8, error) {
	var data int8
	err := binary.Read(r, binary.BigEndian, &data)
	return data, err
}

// UnpackInt16 reads a int8 value from reader.
func UnpackInt16(r io.Reader) (int16, error) {
	var data int16
	err := binary.Read(r, binary.BigEndian, &data)
	return data, err
}

// UnpackInt32 reads a int8 value from reader.
func UnpackInt32(r io.Reader) (int32, error) {
	var data int32
	err := binary.Read(r, binary.BigEndian, &data)
	return data, err
}

// UnpackInt64 reads a int8 value from reader.
func UnpackInt64(r io.Reader) (int64, error) {
	var data int64
	err := binary.Read(r, binary.BigEndian, &data)
	return data, err
}

// UnpackUint8 reads a uint8 value from reader.
func UnpackUint8(r io.Reader) (uint8, error) {
	var data uint8
	err := binary.Read(r, binary.BigEndian, &data)
	return data, err
}

// UnpackUint16 reads a uint16 value from reader.
func UnpackUint16(r io.Reader) (uint16, error) {
	var data uint16
	err := binary.Read(r, binary.BigEndian, &data)
	return data, err
}

// UnpackUint32 reads a uint32 value from reader.
func UnpackUint32(r io.Reader) (uint32, error) {
	var data uint32
	err := binary.Read(r, binary.BigEndian, &data)
	return data, err
}

// UnpackUint64 reads a uint64 value from reader.
func UnpackUint64(r io.Reader) (uint64, error) {
	var data uint64
	err := binary.Read(r, binary.BigEndian, &data)
	return data, err
}

// UnpackFloat32 reads a float32 value from reader.
func UnpackFloat32(r io.Reader) (float32, error) {
	buf := make([]byte, 4)
	if _, err := r.Read(buf); err != nil {
		return 0, err
	}

	bits := binary.LittleEndian.Uint32(buf)
	return math.Float32frombits(bits), nil
}

// UnpackFloat64 reads a float64 value from reader.
func UnpackFloat64(r io.Reader) (float64, error) {
	buf := make([]byte, 8)
	if _, err := r.Read(buf); err != nil {
		return 0, err
	}

	bits := binary.LittleEndian.Uint64(buf)
	return math.Float64frombits(bits), nil
}

func UnpackString5(r io.Reader, len int) (string, error) {
	strBuf := make([]byte, len)
	if _, err := r.Read(strBuf); err != nil {
		return "", err
	}

	return string(strBuf), nil
}

func UnpackString8(r io.Reader) (string, error) {
	var len uint8
	binary.Read(r, binary.BigEndian, &len)

	strBuf := make([]byte, len)
	if _, err := r.Read(strBuf); err != nil {
		return "", err
	}

	return string(strBuf), nil
}

func UnpackString16(r io.Reader) (string, error) {
	var len uint16
	binary.Read(r, binary.BigEndian, &len)

	strBuf := make([]byte, len)
	if _, err := r.Read(strBuf); err != nil {
		return "", err
	}

	return string(strBuf), nil
}

func UnpackString32(r io.Reader) (string, error) {
	var len uint32
	binary.Read(r, binary.BigEndian, &len)

	strBuf := make([]byte, len)
	if _, err := r.Read(strBuf); err != nil {
		return "", err
	}

	return string(strBuf), nil
}

func UnpackBin8(r io.Reader) ([]byte, error) {
	var len uint8
	binary.Read(r, binary.BigEndian, &len)

	bin := make([]byte, len)
	if _, err := r.Read(bin); err != nil {
		return nil, err
	}

	return bin, nil
}

func UnpackBin16(r io.Reader) ([]byte, error) {
	var len uint16
	binary.Read(r, binary.BigEndian, &len)

	bin := make([]byte, len)
	if _, err := r.Read(bin); err != nil {
		return nil, err
	}

	return bin, nil
}

func UnpackBin32(r io.Reader) ([]byte, error) {
	var len uint32
	binary.Read(r, binary.BigEndian, &len)

	bin := make([]byte, len)
	if _, err := r.Read(bin); err != nil {
		return nil, err
	}

	return bin, nil
}
