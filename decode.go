package msgp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
)

// UnpackValue reads a value from the io.Reader. And assigns it to the value pointed by 'ptr'.
// Because UnpackValue depends on the type of 'ptr' to extract a value, 'ptr' should be
// an address of specific type variable.
// If possible, the read value will be converted to the type of variable pointed by 'ptr'.
// If 'ptr' is a pointer of pointer, a new value will be allocated. You don't have to
// allocate new one.
// It is recommended to use this function for all types.
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
	case reflect.Struct:
		err = UnpackStruct(r, ptr)
	case reflect.Ptr:
		err = UnpackPtr(r, ptr)
	case reflect.Interface:
		err = UnpackInterface(r, ptr)
	default:
		return fmt.Errorf("msgp: specified type[%v] is not supported", wantType.Kind())
	}

	return err
}

// UnpackBool reads a bool value from the io.Reader. And assigns it to the value pointed by 'ptr'.
func UnpackBool(r io.Reader, ptr interface{}) error {
	var err error
	var val interface{}

	if val, err = UnpackPrimitive(r); err != nil {
		return err
	}

	if val == nil {
		reflect.ValueOf(ptr).Elem().Set(reflect.Zero(reflect.TypeOf(ptr).Elem()))
	} else {
		if b, ok := val.(bool); ok {
			reflect.ValueOf(ptr).Elem().SetBool(b)
		} else {
			return fmt.Errorf("msgp: unpacked value[%v] is not assignable to bool type", val)
		}
	}
	return nil
}

// UnpackInt reads a integer value from the io.Reader. And assigns it to the value pointed by 'ptr'.
// Numeric types(int, uint, float) are compatible with each other.
// Even float value can be read by a int variable.
func UnpackInt(r io.Reader, ptr interface{}) error {
	var err error
	var val interface{}

	if val, err = UnpackPrimitive(r); err != nil {
		return err
	}

	if val == nil {
		reflect.ValueOf(ptr).Elem().Set(reflect.Zero(reflect.TypeOf(ptr).Elem()))
	} else {
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
	}
	return nil
}

// UnpackUint reads a unsigned integer value from the io.Reader. And assigns it to the value pointed by 'ptr'.
// Numeric types(int, uint, float) are compatible with each other.
// Even float value can be read by a uint variable.
func UnpackUint(r io.Reader, ptr interface{}) error {
	var err error
	var val interface{}

	if val, err = UnpackPrimitive(r); err != nil {
		return err
	}

	if val == nil {
		reflect.ValueOf(ptr).Elem().Set(reflect.Zero(reflect.TypeOf(ptr).Elem()))
	} else {
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
	}
	return nil
}

// UnpackFloat reads a float value from the io.Reader. And assigns it to the value pointed by 'ptr'.
// Numeric types(int, uint, float) are compatible with each other.
// Even int value can be read by a float32 variable.
func UnpackFloat(r io.Reader, ptr interface{}) error {
	var err error
	var val interface{}

	if val, err = UnpackPrimitive(r); err != nil {
		return err
	}

	if val == nil {
		reflect.ValueOf(ptr).Elem().Set(reflect.Zero(reflect.TypeOf(ptr).Elem()))
	} else {
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
	}
	return nil
}

// UnpackString reads a string value from the io.Reader. And assigns it to the value pointed by 'ptr'.
func UnpackString(r io.Reader, ptr interface{}) error {
	var err error
	var val interface{}

	if val, err = UnpackPrimitive(r); err != nil {
		return err
	}

	if val == nil {
		reflect.ValueOf(ptr).Elem().Set(reflect.Zero(reflect.TypeOf(ptr).Elem()))
	} else {
		valTyp := reflect.TypeOf(val)
		if valTyp.Kind() == reflect.String {
			reflect.ValueOf(ptr).Elem().SetString(reflect.ValueOf(val).String())
		} else if valTyp.Kind() == reflect.Slice {
			if valTyp.Elem().Kind() == reflect.Uint8 {
				reflect.ValueOf(ptr).Elem().SetString(string(val.([]byte)))
			} else {
				return fmt.Errorf("msgp: unpacked value[%v] is not assignable to string type", val)
			}
		} else {
			return fmt.Errorf("msgp: unpacked value[%v] is not assignable to string type", val)
		}
	}
	return nil
}

// UnpackArray reads an array from the io.Reader. And assigns it to the value pointed by 'ptr'.
// The length of array type pointed by 'ptr' should be large enough.
// Bin format family value is also read by this function.
// Bin format famaly should be read with a pointer of '[n]byte' or '[n]uint8'.
// (Numeric types are compatible, but bin format family cannot be read with [n]int)
func UnpackArray(r io.Reader, ptr interface{}) error {
	var err error
	var head byte

	arrTyp := reflect.TypeOf(ptr).Elem()
	arrVal := reflect.ValueOf(ptr).Elem()
	arrLen := arrVal.Len()

	if err = binary.Read(r, binary.BigEndian, &head); err != nil {
		return err
	}

	if head == 0xc0 { // nil
		reflect.ValueOf(ptr).Elem().Set(reflect.Zero(reflect.TypeOf(ptr).Elem()))
		return nil
	}

	// handle bin format family
	if head == 0xc4 || head == 0xc5 || head == 0xc6 {
		if arrTyp.Elem().Kind() == reflect.Uint8 {
			var byteSlice []byte

			switch head {
			case 0xc4:
				byteSlice, err = unpackBin8(r)
			case 0xc5:
				byteSlice, err = unpackBin16(r)
			case 0xc6:
				byteSlice, err = unpackBin32(r)
			}
			if err != nil {
				return err
			}

			arrVal.Set(reflect.Zero(reflect.ArrayOf(len(byteSlice), reflect.TypeOf(byteSlice).Elem())))
			reflect.Copy(arrVal, reflect.ValueOf(byteSlice))
			return nil
		}

		return fmt.Errorf("msgp: byte array can't be assigned to other type[%v] array", arrTyp.Elem().Kind())
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
		if err = UnpackValue(r, arrVal.Index(inx).Addr().Interface()); err != nil {
			return err
		}
	}
	return nil
}

// UnpackSlice reads an array from the io.Reader. And assigns it to the value pointed by 'ptr'.
// Bin format family value is also read by this function.
// Bin format famaly should be read with a pointer of '[]byte' or '[]uint8'.
// (Numeric types are compatible, but bin format family cannot be read with []int)
func UnpackSlice(r io.Reader, ptr interface{}) error {
	var err error
	var head byte

	sliceTyp := reflect.TypeOf(ptr).Elem()
	sliceVal := reflect.ValueOf(ptr).Elem()

	if err = binary.Read(r, binary.BigEndian, &head); err != nil {
		return err
	}

	if head == 0xc0 { // nil
		reflect.ValueOf(ptr).Elem().Set(reflect.Zero(reflect.TypeOf(ptr).Elem()))
		return nil
	}

	// handle bin format family
	if head == 0xc4 || head == 0xc5 || head == 0xc6 {
		if sliceTyp.Elem().Kind() == reflect.Uint8 {
			var byteSlice []byte

			switch head {
			case 0xc4:
				byteSlice, err = unpackBin8(r)
			case 0xc5:
				byteSlice, err = unpackBin16(r)
			case 0xc6:
				byteSlice, err = unpackBin32(r)
			}
			if err != nil {
				return err
			}

			sliceVal.Set(reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(byteSlice).Elem()), len(byteSlice), len(byteSlice)))
			reflect.Copy(sliceVal, reflect.ValueOf(byteSlice))
			return nil
		}
		return fmt.Errorf("msgp: byte array can't be assigned to other type[%v] slice", sliceTyp.Elem().Kind())
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
		if err = UnpackValue(r, sliceVal.Index(inx).Addr().Interface()); err != nil {
			return err
		}
	}
	return nil
}

// UnpackMap reads a map from the io.Reader. And assigns it to the value pointed by 'ptr'.
func UnpackMap(r io.Reader, ptr interface{}) error {
	var err error
	var head byte

	mapTyp := reflect.TypeOf(ptr).Elem()
	mapVal := reflect.ValueOf(ptr).Elem()

	if err = binary.Read(r, binary.BigEndian, &head); err != nil {
		return err
	}

	if head == 0xc0 { // nil
		reflect.ValueOf(ptr).Elem().Set(reflect.Zero(reflect.TypeOf(ptr).Elem()))
		return nil
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

// UnpackStruct reads a struct value from the io.Reader. And assigns it to the value pointed by 'ptr'.
// The struct value is deserialized from a map value.
// If the fields of struct are not compatible with the value read, an error is returned.
func UnpackStruct(r io.Reader, ptr interface{}) error {
	var err error
	var head byte

	if err = binary.Read(r, binary.BigEndian, &head); err != nil {
		return err
	}
	if head == 0xc0 { // nil
		reflect.ValueOf(ptr).Elem().Set(reflect.Zero(reflect.TypeOf(ptr).Elem()))
		return nil
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

	type StructField struct {
		Props FieldProps
		Val   reflect.Value
	}
	fieldMap := make(map[string]StructField)

	structTyp := reflect.TypeOf(ptr).Elem()
	structVal := reflect.ValueOf(ptr).Elem()

	structVal.Set(reflect.Zero(structTyp)) // init with zero value

	structNumField := structTyp.NumField()
	for inx := 0; inx < structNumField; inx++ {
		var fp FieldProps

		fieldTyp := structTyp.Field(inx)
		fieldVal := structVal.Field(inx)
		fp.parseTag(fieldTyp)
		if fp.Skip {
			continue
		}
		fieldMap[fp.Name] = StructField{fp, fieldVal}
	}

	for inx := 0; inx < srcLen; inx++ {
		var key string
		if err = UnpackValue(r, &key); err != nil {
			return err
		}

		structField, ok := fieldMap[key]
		if ok {
			if structField.Props.Skip {
				continue
			}

			if structField.Props.String {
				var str string
				if err = UnpackValue(r, &str); err != nil {
					return err
				}
				if err = assignValueFromString(structField.Val, str); err != nil {
					return err
				}
			} else {
				if err = UnpackValue(r, structField.Val.Addr().Interface()); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// UnpackPtr reads a value from the io.Reader. And assigns it to the value pointed by 'ptr'.
// 'ptr' should be a pointer of pointer.
func UnpackPtr(r io.Reader, ptr interface{}) error {
	var err error
	var peek byte

	br := NewPeekableReader(r)
	if peek, err = br.Peek(); err != nil {
		return err
	}
	if peek == 0xc0 { // nil value unpacked.
		reflect.ValueOf(ptr).Elem().Set(reflect.Zero(reflect.TypeOf(ptr).Elem()))
		return nil
	}

	newVal := reflect.New(reflect.TypeOf(ptr).Elem().Elem())
	if err = UnpackValue(br, newVal.Interface()); err != nil { // peeked byte will be consumed in UnpackValue()
		return err
	}

	reflect.ValueOf(ptr).Elem().Set(newVal)
	return nil
}

// UnpackInterface reads a value from the io.Reader. And assigns it to the value pointed by 'ptr'.
// If you don't know the type, you can use this function. but you will have to use reflection to discover the type of the value read.
func UnpackInterface(r io.Reader, ptr interface{}) error {
	var err error

	pi, ok := ptr.(*interface{})
	if !ok {
		return fmt.Errorf("msgp: specified type[%v] is not supported", reflect.TypeOf(ptr).Elem())
	}

	if *pi, err = UnpackPrimitive(r); err != nil {
		return err
	}

	return nil
}

// UnpackPrimitive reads a value from the io.Reader.
// but no type casting takes place.
// It is generally recommended to use UnpackValue().
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
		return unpackInt8(r)
	} else if head == 0xd1 {
		return unpackInt16(r)
	} else if head == 0xd2 {
		return unpackInt32(r)
	} else if head == 0xd3 {
		return unpackInt64(r)
	} else if head == 0xcc {
		return unpackUint8(r)
	} else if head == 0xcd {
		return unpackUint16(r)
	} else if head == 0xce {
		return unpackUint32(r)
	} else if head == 0xcf {
		return unpackUint64(r)
	} else if head == 0xca {
		return unpackFloat32(r)
	} else if head == 0xcb {
		return unpackFloat64(r)
	} else if head&0xe0 == 0xa0 {
		return unpackString5(r, int(head&0x1f))
	} else if head == 0xd9 {
		return unpackString8(r)
	} else if head == 0xda {
		return unpackString16(r)
	} else if head == 0xdb {
		return unpackString32(r)
	} else if head == 0xc4 { // bin
		return unpackBin8(r)
	} else if head == 0xc5 {
		return unpackBin16(r)
	} else if head == 0xc6 {
		return unpackBin32(r)
	} else if head&0xf0 == 0x90 { // array
		return unpackArray4(r, int(head&0x0f))
	} else if head == 0xdc {
		return unpackArray16(r)
	} else if head == 0xdd {
		return unpackArray32(r)
	} else if head&0xf0 == 0x80 { // map
		return unpackMap4(r, int(head&0x0f))
	} else if head == 0xde {
		return unpackMap16(r)
	} else if head == 0xdf {
		return unpackMap32(r)
	}

	return nil, errors.New("msgp: UnpackPrimitive() reads unsupported(array, map) format family")
}

func unpackInt8(r io.Reader) (int8, error) {
	var val int8
	err := binary.Read(r, binary.BigEndian, &val)
	return val, err
}

func unpackInt16(r io.Reader) (int16, error) {
	var val int16
	err := binary.Read(r, binary.BigEndian, &val)
	return val, err
}

func unpackInt32(r io.Reader) (int32, error) {
	var val int32
	err := binary.Read(r, binary.BigEndian, &val)
	return val, err
}

func unpackInt64(r io.Reader) (int64, error) {
	var val int64
	err := binary.Read(r, binary.BigEndian, &val)
	return val, err
}

func unpackUint8(r io.Reader) (uint8, error) {
	var val uint8
	err := binary.Read(r, binary.BigEndian, &val)
	return val, err
}

func unpackUint16(r io.Reader) (uint16, error) {
	var val uint16
	err := binary.Read(r, binary.BigEndian, &val)
	return val, err
}

func unpackUint32(r io.Reader) (uint32, error) {
	var val uint32
	err := binary.Read(r, binary.BigEndian, &val)
	return val, err
}

func unpackUint64(r io.Reader) (uint64, error) {
	var val uint64
	err := binary.Read(r, binary.BigEndian, &val)
	return val, err
}

func unpackFloat32(r io.Reader) (float32, error) {
	buf := make([]byte, 4)
	if _, err := r.Read(buf); err != nil {
		return 0, err
	}

	bits := binary.BigEndian.Uint32(buf)
	return math.Float32frombits(bits), nil
}

func unpackFloat64(r io.Reader) (float64, error) {
	buf := make([]byte, 8)
	if _, err := r.Read(buf); err != nil {
		return 0, err
	}

	bits := binary.BigEndian.Uint64(buf)
	return math.Float64frombits(bits), nil
}

func unpackString5(r io.Reader, len int) (string, error) {
	return unpackStringBody(r, len)
}

func unpackString8(r io.Reader) (string, error) {
	var len uint8
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return "", err
	}
	return unpackStringBody(r, int(len))
}

func unpackString16(r io.Reader) (string, error) {
	var len uint16
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return "", err
	}
	return unpackStringBody(r, int(len))
}

func unpackString32(r io.Reader) (string, error) {
	var len uint32
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return "", err
	}
	return unpackStringBody(r, int(len))
}

func unpackBin8(r io.Reader) ([]byte, error) {
	var len uint8
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return nil, err
	}
	return unpackBinBody(r, int(len))
}

func unpackBin16(r io.Reader) ([]byte, error) {
	var len uint16
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return nil, err
	}
	return unpackBinBody(r, int(len))
}

func unpackBin32(r io.Reader) ([]byte, error) {
	var len uint32
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return nil, err
	}
	return unpackBinBody(r, int(len))
}

func unpackArray4(r io.Reader, len int) (interface{}, error) {
	return unpackArrayBody(r, len)
}

func unpackArray16(r io.Reader) (interface{}, error) {
	var len uint16
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return nil, err
	}
	return unpackArrayBody(r, int(len))
}

func unpackArray32(r io.Reader) (interface{}, error) {
	var len uint32
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return nil, err
	}
	return unpackArrayBody(r, int(len))
}

func unpackMap4(r io.Reader, len int) (interface{}, error) {
	return unpackMapBody(r, len)
}

func unpackMap16(r io.Reader) (interface{}, error) {
	var len uint16
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return nil, err
	}
	return unpackMapBody(r, int(len))
}

func unpackMap32(r io.Reader) (interface{}, error) {
	var len uint32
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return nil, err
	}
	return unpackMapBody(r, int(len))
}

func unpackStringBody(r io.Reader, len int) (string, error) {
	if len == 0 {
		return "", nil
	}

	var n int
	var err error
	str := make([]byte, len)
	if n, err = r.Read(str); err != nil {
		return "", err
	}
	if n != len {
		return "", errors.New("msgp: broken string format family was found")
	}

	return string(str), nil
}

func unpackBinBody(r io.Reader, len int) ([]byte, error) {
	if len == 0 {
		return nil, nil // nil as an empty slice
	}

	var n int
	var err error
	bin := make([]byte, len)
	if n, err = r.Read(bin); err != nil {
		return nil, err
	}
	if n != len {
		return nil, errors.New("msgp: broken binary format family was found")
	}

	return bin, nil
}

func unpackArrayBody(r io.Reader, len int) (interface{}, error) {
	if len == 0 {
		return nil, nil // nil as an empty slice
	}

	var err error
	var val interface{}

	slice := make([]interface{}, len, len)
	for inx := 0; inx < len; inx++ {
		if val, err = UnpackPrimitive(r); err != nil {
			return nil, err
		}
		slice[inx] = val
	}
	return slice, nil
}

func unpackMapBody(r io.Reader, len int) (interface{}, error) {
	if len == 0 {
		return nil, nil // nil as an empty map
	}

	var err error
	var key, val interface{}

	mapVal := make(map[interface{}]interface{})
	for inx := 0; inx < len; inx++ {
		if key, err = UnpackPrimitive(r); err != nil {
			return nil, err
		}
		if val, err = UnpackPrimitive(r); err != nil {
			return nil, err
		}
		mapVal[key] = val
	}
	return mapVal, nil
}

func assignValueFromString(dest reflect.Value, str string) error {
	var err error

	if str == "null" || str == "nil" {
		dest.Set(reflect.Zero(dest.Type()))
	} else {
		switch dest.Type().Kind() {
		case reflect.Bool:
			var b bool
			if b, err = strconv.ParseBool(str); err != nil {
				return err
			}
			dest.SetBool(b)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var i int64
			if i, err = strconv.ParseInt(str, 10, 64); err != nil {
				return err
			}
			dest.SetInt(i)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var u uint64
			if u, err = strconv.ParseUint(str, 10, 64); err != nil {
				return err
			}
			dest.SetUint(u)
		case reflect.Float32, reflect.Float64:
			var f float64
			if f, err = strconv.ParseFloat(str, 64); err != nil {
				return err
			}
			dest.SetFloat(f)
		case reflect.String:
			dest.SetString(str)
		case reflect.Ptr:
			dest.Set(reflect.New(dest.Elem().Type()))
			return assignValueFromString(dest.Elem(), str)
		}
	}

	return nil
}
