package msgp

import (
	"encoding/binary"
	"errors"
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
		// TODO
	}

	return nil
}

func UnpackBool(r io.Reader, ptr interface{}) error {
	var val interface{}
	var err error
	if val, err := UnpackPrimitive(r); err != nil {
		return err
	}

	reflect.ValueOf(ptr).Elem().SetBool(val.(bool))
	return nil
}

func UnpackInt(r io.Reader, ptr interface{}) error {
	var val interface{}
	var err error
	if val, err := UnpackPrimitive(r); err != nil {
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
		reflect.ValueOf(ptr).Elem().Set(reflect.ValueOf(val)) // if not assignable, will panic.
	}
	return nil
}

func UnpackUint(r io.Reader, ptr interface{}) error {
	var val interface{}
	var err error
	if val, err := UnpackPrimitive(r); err != nil {
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
		reflect.ValueOf(ptr).Elem().Set(reflect.ValueOf(val)) // if not assignable, will panic.
	}
	return nil
}

func UnpackFloat(r io.Reader, ptr interface{}) error {
	var val interface{}
	var err error
	if val, err := UnpackPrimitive(r); err != nil {
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
		reflect.ValueOf(ptr).Elem().Set(reflect.ValueOf(val)) // if not assignable, will panic.
	}
	return nil
}

func UnpackString(r io.Reader, ptr interface{}) error {
	var val interface{}
	var err error
	if val, err := UnpackPrimitive(r); err != nil {
		return err
	}

	reflect.ValueOf(ptr).Elem().SetString(reflect.ValueOf(val).String())
	return nil
}

func UnpackArray(r io.Reader, ptr interface{}) error {
	var val interface{}
	var err error
	var head byte

	arrTyp := reflect.TypeOf(ptr).Elem()
	arrVal := reflect.ValueOf(ptr).Elem()
	arrLen := arrVal.Len()

	if err = binary.Read(r, binary.BigEndian, &head); err != nil {
		return err
	}

	var srcLen = 0
	if head == 0xc4 { // bin
		var temp uint8
		binary.Read(r, binary.BigEndian, &temp)
		srcLen = int(temp)
	} else if head == 0xc5 {
		var temp uint16
		binary.Read(r, binary.BigEndian, &temp)
		srcLen = int(temp)
	} else if head == 0xc6 {
		var temp uint32
		binary.Read(r, binary.BigEndian, &temp)
		srcLen = int(temp) // maybe overflow.
	} else if head&0xf0 == 0x90 { // array
		var temp uint8
		srcLen = int(head & 0x0f)
	} else if head == 0xdc {
		var temp uint16
		binary.Read(r, binary.BigEndian, &temp)
		srcLen = int(temp)
	} else if head == 0xdd {
		var temp uint32
		binary.Read(r, binary.BigEndian, &temp)
		srcLen = int(temp) // maybe overflow.
	} else {
		panic("!!!") // TODO
	}

	if arrLen < srcLen {
		panic("!!!") // TODO
	}

	arrVal.Set(reflect.Zero(reflect.ArrayOf(arrLen, arrTyp.Elem()))) // array 생성.
	for inx := 0; inx < srcLen; inx++ {
		UnpackValue(r, arrVal.Index(inx).Addr())
	}
	return nil
}

func UnpackSlice(r io.Reader, ptr interface{}) error {
	var val interface{}
	var err error
	var head byte

	arrTyp := reflect.TypeOf(ptr).Elem()
	arrVal := reflect.ValueOf(ptr).Elem()

	if err = binary.Read(r, binary.BigEndian, &head); err != nil {
		return err
	}

	var srcLen = 0
	if head == 0xc4 { // bin
		var temp uint8
		err = binary.Read(r, binary.BigEndian, &temp) // TODO handle err
		srcLen = int(temp)
	} else if head == 0xc5 {
		var temp uint16
		err = binary.Read(r, binary.BigEndian, &temp) // TODO handle err
		srcLen = int(temp)
	} else if head == 0xc6 {
		var temp uint32
		err = binary.Read(r, binary.BigEndian, &temp) // TODO handle err
		srcLen = int(temp) // maybe overflow.
	} else if head&0xf0 == 0x90 { // array
		var temp uint8
		srcLen = int(head & 0x0f)
	} else if head == 0xdc {
		var temp uint16
		err = binary.Read(r, binary.BigEndian, &temp) // TODO handle err
		srcLen = int(temp)
	} else if head == 0xdd {
		var temp uint32
		err = binary.Read(r, binary.BigEndian, &temp) // TODO handle err
		srcLen = int(temp) // maybe overflow.
	} else {
		panic("!!!") // TODO
	}

	arrVal.Set(reflect.MakeSlice(reflect.SliceOf(arrTyp.Elem()), srcLen, srcLen)) // slice 생성.
	for inx := 0; inx < srcLen; inx++ {
		UnpackValue(r, arrVal.Index(inx).Addr())
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
		var temp uint8
		srcLen = int(head & 0x0f)
	} else if head == 0xde {
		var temp uint16
		err = binary.Read(r, binary.BigEndian, &temp) // TODO handle err
		srcLen = int(temp)
	} else if head == 0xdf {
		var temp uint32
		err = binary.Read(r, binary.BigEndian, &temp) // TODO handle err
		srcLen = int(temp) // maybe overflow.
	} else {
		panic("!!!") // TODO
	}

	mapVal.Set(reflect.MakeMap(reflect.MapOf(mapTyp.Key(), mapTyp.Elem()))) // map 생성.
	for inx := 0; inx < srcLen; inx++ {
		keyPtr := reflect.New(mapTyp.Key())
		err = UnpackValue(r, keyPtr.Interface())

		valPtr := reflect.New(mapTyp.Elem())
		err = UnpackValue(r, valPtr.Interface())

		mapVal.SetMapIndex(keyPtr.Elem(), valPtr.Elem())
	}
	return nil
}

// UnpackValueWithPtr reads a value from reader to pointer.
func UnpackValueWithPtr(r io.Reader, vp interface{}) error {
	var val interface{}
	var err error

	if val, err = UnpackValue(r); err != nil {
		return err
	}

	wantType := reflect.TypeOf(vp).Elem()
	switch wantType.Kind() {
	case reflect.Bool:
		reflect.ValueOf(vp).Elem().SetBool(val.(bool))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch reflect.ValueOf(val).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			reflect.ValueOf(vp).Elem().SetInt(reflect.ValueOf(val).Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			reflect.ValueOf(vp).Elem().SetInt(int64(reflect.ValueOf(val).Uint()))
		case reflect.Float32, reflect.Float64:
			reflect.ValueOf(vp).Elem().SetInt(int64(reflect.ValueOf(val).Float()))
		default:
			reflect.ValueOf(vp).Elem().Set(reflect.ValueOf(val)) // if not assignable, will panic.
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch reflect.ValueOf(val).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			reflect.ValueOf(vp).Elem().SetUint(uint64(reflect.ValueOf(val).Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			reflect.ValueOf(vp).Elem().SetUint(reflect.ValueOf(val).Uint())
		case reflect.Float32, reflect.Float64:
			reflect.ValueOf(vp).Elem().SetUint(uint64(reflect.ValueOf(val).Float()))
		default:
			reflect.ValueOf(vp).Elem().Set(reflect.ValueOf(val)) // if not assignable, will panic.
		}

	case reflect.Float32, reflect.Float64:
		switch reflect.ValueOf(val).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			reflect.ValueOf(vp).Elem().SetFloat(float64(reflect.ValueOf(val).Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			reflect.ValueOf(vp).Elem().SetFloat(float64(reflect.ValueOf(val).Uint()))
		case reflect.Float32, reflect.Float64:
			reflect.ValueOf(vp).Elem().SetFloat(reflect.ValueOf(val).Float())
		default:
			reflect.ValueOf(vp).Elem().Set(reflect.ValueOf(val)) // if not assignable, will panic.
		}

	case reflect.String:
		reflect.ValueOf(vp).Elem().SetString(reflect.ValueOf(val).String())

	case reflect.Array:
		srcVal := reflect.ValueOf(val)
		srcLen := srcVal.Len() // if src is not an array, will panic
		tgtVal := reflect.ValueOf(vp).Elem()
		tgtVal.Set(reflect.Zero(reflect.ArrayOf(srcLen, wantType.Elem())))
		for inx := 0; inx < srcLen; inx++ {
			switch wantType.Elem().Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				switch srcVal.Index(inx).Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					tgtVal.Index(inx).SetInt(srcVal.Index(inx).Int())
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					tgtVal.Index(inx).SetInt(int64(srcVal.Index(inx).Uint()))
				case reflect.Float32, reflect.Float64:
					tgtVal.Index(inx).SetInt(int64(srcVal.Index(inx).Float()))
				default:
					tgtVal.Index(inx).Set(srcVal.Index(inx)) // if not assignable, will panic.
				}

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				switch srcVal.Index(inx).Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					tgtVal.Index(inx).SetUint(uint64(srcVal.Index(inx).Int()))
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					tgtVal.Index(inx).SetUint(srcVal.Index(inx).Uint())
				case reflect.Float32, reflect.Float64:
					tgtVal.Index(inx).SetUint(uint64(srcVal.Index(inx).Float()))
				default:
					tgtVal.Index(inx).Set(srcVal.Index(inx)) // if not assignable, will panic.
				}

			case reflect.Float32, reflect.Float64:
				switch srcVal.Index(inx).Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					tgtVal.Index(inx).SetFloat(float64(srcVal.Index(inx).Int()))
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					tgtVal.Index(inx).SetFloat(float64(srcVal.Index(inx).Uint()))
				case reflect.Float32, reflect.Float64:
					tgtVal.Index(inx).SetFloat(srcVal.Index(inx).Float())
				default:
					tgtVal.Index(inx).Set(srcVal.Index(inx)) // if not assignable, will panic.
				}
			default:
				tgtVal.Index(inx).Set(srcVal.Index(inx)) // non-numeric values are just set.
			}
		}

	case reflect.Slice:
		srcVal := reflect.ValueOf(val)
		srcLen := srcVal.Len() // if src is not an array, will panic
		tgtVal := reflect.ValueOf(vp).Elem()
		tgtVal.Set(reflect.MakeSlice(wantType.Elem(), srcLen, srcLen))
		for inx := 0; inx < srcLen; inx++ {
			switch wantType.Elem().Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				switch srcVal.Index(inx).Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					tgtVal.Index(inx).SetInt(srcVal.Index(inx).Int())
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					tgtVal.Index(inx).SetInt(int64(srcVal.Index(inx).Uint()))
				case reflect.Float32, reflect.Float64:
					tgtVal.Index(inx).SetInt(int64(srcVal.Index(inx).Float()))
				default:
					tgtVal.Index(inx).Set(srcVal.Index(inx)) // if not assignable, will panic.
				}

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				switch srcVal.Index(inx).Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					tgtVal.Index(inx).SetUint(uint64(srcVal.Index(inx).Int()))
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					tgtVal.Index(inx).SetUint(srcVal.Index(inx).Uint())
				case reflect.Float32, reflect.Float64:
					tgtVal.Index(inx).SetUint(uint64(srcVal.Index(inx).Float()))
				default:
					tgtVal.Index(inx).Set(srcVal.Index(inx)) // if not assignable, will panic.
				}

			case reflect.Float32, reflect.Float64:
				switch srcVal.Index(inx).Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					tgtVal.Index(inx).SetFloat(float64(srcVal.Index(inx).Int()))
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					tgtVal.Index(inx).SetFloat(float64(srcVal.Index(inx).Uint()))
				case reflect.Float32, reflect.Float64:
					tgtVal.Index(inx).SetFloat(srcVal.Index(inx).Float())
				default:
					tgtVal.Index(inx).Set(srcVal.Index(inx)) // if not assignable, will panic.
				}
			default:
				tgtVal.Index(inx).Set(srcVal.Index(inx)) // if not assignable, will panic.
			}
		}

	case reflect.Map:
		srcVal := reflect.ValueOf(val)
		tgtVal := reflect.ValueOf(vp).Elem()
		tgtVal.Set(reflect.MakeMap(reflect.MapOf(wantType.Key(), wantType.Elem())))
		for _, keyVal := range srcVal.MapKeys() {
			var key interface{}
			var data interface{}
			switch wantType.Key().Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				// var key int64
				switch keyVal.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					key = keyVal.Int()
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					key = int64(keyVal.Uint())
				case reflect.Float32, reflect.Float64:
					key = int64(keyVal.Float())
				default:
					key = keyVal.Interface().(int64) // if not assignable, will panic.
				}

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				// var key uint64
				switch keyVal.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					key = uint64(keyVal.Int())
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					key = keyVal.Uint()
				case reflect.Float32, reflect.Float64:
					key = uint64(keyVal.Float())
				default:
					key = keyVal.Interface().(uint64) // if not assignable, will panic.
				}

			case reflect.Float32, reflect.Float64:
				// var key float64
				switch keyVal.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					key = float64(keyVal.Int())
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					key = float64(keyVal.Uint())
				case reflect.Float32, reflect.Float64:
					key = keyVal.Float()
				default:
					key = keyVal.Interface().(float64) // if not assignable, will panic.
				}

			default:
				key = keyVal.Interface() // if not assignable, will panic.
			}

			dataVal := srcVal.MapIndex(keyVal)
			switch wantType.Elem().Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				// var data int64
				switch dataVal.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					data = dataVal.Int()
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					data = int64(dataVal.Uint())
				case reflect.Float32, reflect.Float64:
					data = int64(dataVal.Float())
				default:
					data = dataVal.Interface().(int64) // if not assignable, will panic.
				}

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				// var data uint64
				switch dataVal.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					data = uint64(dataVal.Int())
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					data = dataVal.Uint()
				case reflect.Float32, reflect.Float64:
					data = uint64(dataVal.Float())
				default:
					data = dataVal.Interface().(uint64) // if not assignable, will panic.
				}

			case reflect.Float32, reflect.Float64:
				// var data float64
				switch dataVal.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					data = float64(dataVal.Int())
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					data = float64(dataVal.Uint())
				case reflect.Float32, reflect.Float64:
					data = dataVal.Float()
				default:
					data = dataVal.Interface().(float64) // if not assignable, will panic.
				}

			default:
				data = dataVal.Interface() // if not assignable, will panic.
			}

			tgtVal.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(data))
		}
	}

	return errors.New("msgp: Specified type is not supported.")
}

// UnpackValue reads a value from reader.
func UnpackPrimitive(r io.Reader) (interface{}, error) {
	var head byte

	if err := binary.Read(r, binary.BigEndian, &head); err != nil {
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
	// } else if head == 0xc4 { // bin
	// 	return UnpackBin8(r)
	// } else if head == 0xc5 {
	// 	return UnpackBin16(r)
	// } else if head == 0xc6 {
	// 	return UnpackBin32(r)
	// } else if head&0xf0 == 0x90 { // array
	// 	return UnpackArray4(r, int(head&0x0f))
	// } else if head == 0xdc {
	// 	return UnpackArray16(r)
	// } else if head == 0xdd {
	// 	return UnpackArray32(r)
	// } else if head&0xf0 == 0x80 { // map
	// 	return UnpackMap4(r, int(head&0x0f))
	// } else if head == 0xde {
	// 	return UnpackMap16(r)
	// } else if head == 0xdf {
	// 	return UnpackMap32(r)
	// }

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

// func UnpackBin8(r io.Reader) ([]byte, error) {
// 	var len uint8
// 	binary.Read(r, binary.BigEndian, &len)

// 	bin := make([]byte, len)
// 	if _, err := r.Read(bin); err != nil {
// 		return nil, err
// 	}

// 	return bin, nil
// }

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

func UnpackArray4(r io.Reader, len int) (interface{}, error) {
	if len == 0 {
		return nil, nil // nil as an empty slice
	}

	var val interface{}
	var err error
	if val, err = UnpackValue(r); err != nil {
		return nil, err
	}

	sliceVal := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(val)), 0, len)
	sliceVal = reflect.Append(sliceVal, reflect.ValueOf(val))
	for inx := 1; inx < len; inx++ {
		if val, err = UnpackValue(r); err != nil {
			return nil, err
		}
		sliceVal = reflect.Append(sliceVal, reflect.ValueOf(val))
	}

	return sliceVal.Interface(), nil
}

func UnpackArray16(r io.Reader) (interface{}, error) {
	var len uint16
	binary.Read(r, binary.BigEndian, &len)

	if len == 0 {
		return nil, nil // nil as an empty slice
	}

	var val interface{}
	var err error
	if val, err = UnpackValue(r); err != nil {
		return nil, err
	}

	sliceVal := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(val)), 0, int(len))
	sliceVal = reflect.Append(sliceVal, reflect.ValueOf(val))
	for inx := uint16(1); inx < len; inx++ {
		if val, err = UnpackValue(r); err != nil {
			return nil, err
		}
		sliceVal = reflect.Append(sliceVal, reflect.ValueOf(val))
	}

	return sliceVal.Interface(), nil
}

func UnpackArray32(r io.Reader) (interface{}, error) {
	var len uint32
	binary.Read(r, binary.BigEndian, &len)

	if len == 0 {
		return nil, nil // nil as an empty slice
	}

	var val interface{}
	var err error
	if val, err = UnpackValue(r); err != nil {
		return nil, err
	}

	sliceVal := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(val)), 0, int(len))
	sliceVal = reflect.Append(sliceVal, reflect.ValueOf(val))
	for inx := uint32(1); inx < len; inx++ {
		if val, err = UnpackValue(r); err != nil {
			return nil, err
		}
		sliceVal = reflect.Append(sliceVal, reflect.ValueOf(val))
	}

	return sliceVal.Interface(), nil
}

func UnpackMap4(r io.Reader, len int) (interface{}, error) {
	if len == 0 {
		return nil, nil // nil as an empty map
	}

	var key, val interface{}
	var err error
	if key, err = UnpackValue(r); err != nil {
		return nil, err
	}
	if val, err = UnpackValue(r); err != nil {
		return nil, err
	}

	mapVal := reflect.MakeMap(reflect.MapOf(reflect.TypeOf(key), reflect.TypeOf(val)))
	mapVal.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
	for inx := 1; inx < len; inx++ {
		if key, err = UnpackValue(r); err != nil {
			return nil, err
		}
		if val, err = UnpackValue(r); err != nil {
			return nil, err
		}
		mapVal.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
	}

	return mapVal.Interface(), nil
}

func UnpackMap16(r io.Reader) (interface{}, error) {
	var len uint16
	binary.Read(r, binary.BigEndian, &len)

	if len == 0 {
		return nil, nil // nil as an empty map
	}

	var key, val interface{}
	var err error
	if key, err = UnpackValue(r); err != nil {
		return nil, err
	}
	if val, err = UnpackValue(r); err != nil {
		return nil, err
	}

	mapVal := reflect.MakeMap(reflect.MapOf(reflect.TypeOf(key), reflect.TypeOf(val)))
	mapVal.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
	for inx := uint16(1); inx < len; inx++ {
		if key, err = UnpackValue(r); err != nil {
			return nil, err
		}
		if val, err = UnpackValue(r); err != nil {
			return nil, err
		}
		mapVal.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
	}

	return mapVal.Interface(), nil
}

func UnpackMap32(r io.Reader) (interface{}, error) {
	var len uint32
	binary.Read(r, binary.BigEndian, &len)

	if len == 0 {
		return nil, nil // nil as an empty map
	}

	var key, val interface{}
	var err error
	if key, err = UnpackValue(r); err != nil {
		return nil, err
	}
	if val, err = UnpackValue(r); err != nil {
		return nil, err
	}

	mapVal := reflect.MakeMap(reflect.MapOf(reflect.TypeOf(key), reflect.TypeOf(val)))
	mapVal.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
	for inx := uint32(1); inx < len; inx++ {
		if key, err = UnpackValue(r); err != nil {
			return nil, err
		}
		if val, err = UnpackValue(r); err != nil {
			return nil, err
		}
		mapVal.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
	}

	return mapVal.Interface(), nil
}
