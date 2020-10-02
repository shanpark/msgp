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

// UnpackValue reads an object from reader. And assigns to the value pointed by 'ptr'.
// Because UnpackValue depends on the type of ptr to extract values, 'ptr'
// must be an address of specific type.
// If possible, read value is converted to the type of ptr.
// If 'ptr' is a pointer of pointer, array, map type, a new value will be
// allocated. and the address of the new value will be assigned to 'ptr'
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

// UnpackBool reads a bool value from reader. And assigns to the value pointed by ptr.
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

// UnpackInt reads a integer value from reader. And assigns to the value pointed by ptr.
// Numeric types(int, uint, float) are compatible with each other.
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

// UnpackUint reads a unsigned integer value from reader. And assigns to the value pointed by ptr.
// Numeric types(int, uint, float) are compatible with each other.
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

// UnpackFloat reads a float value from reader. And assigns to the value pointed by ptr.
// Numeric types(int, uint, float) are compatible with each other.
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

// UnpackString reads a string value from reader. And assigns to the value pointed by ptr.
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

// UnpackArray reads a array value from reader. And assigns to the value pointed by ptr.
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

// UnpackSlice reads a array value from reader. And assigns to the value pointed by ptr.
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

// UnpackMap reads a map value from reader. And assigns to the value pointed by ptr.
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

// UnpackPtr reads a msgp object from reader. And assigns to the value pointed by ptr.
// ptr should be a pointer of pointer. And a new object filled with read object will be allocated.
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

// UnpackPtr reads a msgp object from reader. And assigns to the value pointed by ptr.
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

// UnpackPrimitive reads a primitive value from reader.
// Primitive values mean the values of nil, pool, int, uint, float, string, [] bytes.
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
	} else if head&0xf0 == 0x90 { // array
		return UnpackArray4(r, int(head&0x0f))
	} else if head == 0xdc {
		return UnpackArray16(r)
	} else if head == 0xdd {
		return UnpackArray32(r)
	} else if head&0xf0 == 0x80 { // map
		return UnpackMap4(r, int(head&0x0f))
	} else if head == 0xde {
		return UnpackMap16(r)
	} else if head == 0xdf {
		return UnpackMap32(r)
	}

	return nil, errors.New("msgp: UnpackPrimitive() reads unsupported(array, map) format family")
}

// UnpackInt8 reads a int8 object from reader.
func UnpackInt8(r io.Reader) (int8, error) {
	var val int8
	err := binary.Read(r, binary.BigEndian, &val)
	return val, err
}

// UnpackInt16 reads a int16 object from reader.
func UnpackInt16(r io.Reader) (int16, error) {
	var val int16
	err := binary.Read(r, binary.BigEndian, &val)
	return val, err
}

// UnpackInt32 reads a int32 object from reader.
func UnpackInt32(r io.Reader) (int32, error) {
	var val int32
	err := binary.Read(r, binary.BigEndian, &val)
	return val, err
}

// UnpackInt64 reads a int64 object from reader.
func UnpackInt64(r io.Reader) (int64, error) {
	var val int64
	err := binary.Read(r, binary.BigEndian, &val)
	return val, err
}

// UnpackUint8 reads a uint8 object from reader.
func UnpackUint8(r io.Reader) (uint8, error) {
	var val uint8
	err := binary.Read(r, binary.BigEndian, &val)
	return val, err
}

// UnpackUint16 reads a uint16 object from reader.
func UnpackUint16(r io.Reader) (uint16, error) {
	var val uint16
	err := binary.Read(r, binary.BigEndian, &val)
	return val, err
}

// UnpackUint32 reads a uint32 object from reader.
func UnpackUint32(r io.Reader) (uint32, error) {
	var val uint32
	err := binary.Read(r, binary.BigEndian, &val)
	return val, err
}

// UnpackUint64 reads a uint64 object from reader.
func UnpackUint64(r io.Reader) (uint64, error) {
	var val uint64
	err := binary.Read(r, binary.BigEndian, &val)
	return val, err
}

// UnpackFloat32 reads a float32 object from reader.
func UnpackFloat32(r io.Reader) (float32, error) {
	buf := make([]byte, 4)
	if _, err := r.Read(buf); err != nil {
		return 0, err
	}

	bits := binary.BigEndian.Uint32(buf)
	return math.Float32frombits(bits), nil
}

// UnpackFloat64 reads a float64 object from reader.
func UnpackFloat64(r io.Reader) (float64, error) {
	buf := make([]byte, 8)
	if _, err := r.Read(buf); err != nil {
		return 0, err
	}

	bits := binary.BigEndian.Uint64(buf)
	return math.Float64frombits(bits), nil
}

// UnpackString5 reads a string object with length of 5 bits from reader.
func UnpackString5(r io.Reader, len int) (string, error) {
	return unpackStringBody(r, len)
}

// UnpackString8 reads a string object with length of 8 bits from reader.
func UnpackString8(r io.Reader) (string, error) {
	var len uint8
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return "", err
	}
	return unpackStringBody(r, int(len))
}

// UnpackString16 reads a string object with length of 16 bits from reader.
func UnpackString16(r io.Reader) (string, error) {
	var len uint16
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return "", err
	}
	return unpackStringBody(r, int(len))
}

// UnpackString32 reads a string object with length of 32 bits from reader.
func UnpackString32(r io.Reader) (string, error) {
	var len uint32
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return "", err
	}
	return unpackStringBody(r, int(len))
}

// UnpackBin8 reads a binary object with length of 8 bits from reader.
func UnpackBin8(r io.Reader) ([]byte, error) {
	var len uint8
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return nil, err
	}
	return unpackBinBody(r, int(len))
}

// UnpackBin16 reads a binary object with length of 16 bits from reader.
func UnpackBin16(r io.Reader) ([]byte, error) {
	var len uint16
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return nil, err
	}
	return unpackBinBody(r, int(len))
}

// UnpackBin32 reads a binary object with length of 32 bits from reader.
func UnpackBin32(r io.Reader) ([]byte, error) {
	var len uint32
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return nil, err
	}
	return unpackBinBody(r, int(len))
}

// UnpackArray4 reads an array object with length of 4 bits from readder.
func UnpackArray4(r io.Reader, len int) (interface{}, error) {
	return unpackArrayBody(r, len)
}

// UnpackArray16 reads an array object with length of 16 bits from readder.
func UnpackArray16(r io.Reader) (interface{}, error) {
	var len uint16
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return nil, err
	}
	return unpackArrayBody(r, int(len))
}

// UnpackArray32 reads an array object with length of 32 bits from readder.
func UnpackArray32(r io.Reader) (interface{}, error) {
	var len uint32
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return nil, err
	}
	return unpackArrayBody(r, int(len))
}

// UnpackMap4 reads an map object with length of 4 bits from readder.
func UnpackMap4(r io.Reader, len int) (interface{}, error) {
	return unpackMapBody(r, len)
}

// UnpackMap16 reads an map object with length of 16 bits from readder.
func UnpackMap16(r io.Reader) (interface{}, error) {
	var len uint16
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return nil, err
	}
	return unpackMapBody(r, int(len))
}

// UnpackMap32 reads an map object with length of 32 bits from readder.
func UnpackMap32(r io.Reader) (interface{}, error) {
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
		return "", errors.New("msgp: broken string format family was found.")
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
		return nil, errors.New("msgp: broken binary format family was found.")
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
