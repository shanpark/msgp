package msgp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
)

// PackValue serialize a value.
func PackValue(w io.Writer, value interface{}) (err error) {
	if value == nil {
		err = PackNil(w)
	}

	switch reflect.ValueOf(value).Kind() {
	case reflect.Bool:
		err = PackBool(w, value.(bool))
	case reflect.Int:
		err = PackInt(w, int64(value.(int)))
	case reflect.Int8:
		err = PackInt(w, int64(value.(int8)))
	case reflect.Int16:
		err = PackInt(w, int64(value.(int16)))
	case reflect.Int32:
		err = PackInt(w, int64(value.(int32)))
	case reflect.Int64:
		err = PackInt(w, value.(int64))
	case reflect.Uint:
		err = PackUint(w, uint64(value.(uint)))
	case reflect.Uint8:
		err = PackUint(w, uint64(value.(uint8)))
	case reflect.Uint16:
		err = PackUint(w, uint64(value.(uint16)))
	case reflect.Uint32:
		err = PackUint(w, uint64(value.(uint32)))
	case reflect.Uint64:
		err = PackUint(w, value.(uint64))
	case reflect.Float32:
		err = PackFloat32(w, value.(float32))
	case reflect.Float64:
		err = PackFloat64(w, value.(float64))
	case reflect.String:
		err = PackString(w, value.(string))
	case reflect.Struct:
		err = PackStruct(w, value)
	case reflect.Map:
		err = PackMap(w, value)
	case reflect.Array, reflect.Slice:
		err = PackArray(w, value)
	case reflect.Ptr:
		err = PackPtr(w, value)
	}

	return err
}

// PackNil writes nil value.
func PackNil(w io.Writer) error {
	w.Write([]byte{0xc0})
	return nil
}

// PackBool writes a bool data to writer.
func PackBool(w io.Writer, value bool) error {
	var err error

	if value {
		_, err = w.Write([]byte{0xc3})
	} else {
		_, err = w.Write([]byte{0xc2})
	}

	return err
}

// PackInt writes an int type data to writer.
func PackInt(w io.Writer, value int64) error {
	var buf bytes.Buffer

	if (value >= -32) && (value <= 127) {
		buf.WriteByte(byte(value))
	} else if (value >= -0x80) && (value <= 0x7f) {
		buf.WriteByte(0xd0)
		buf.WriteByte(byte(value))
	} else if (value >= -0x8000) && (value <= 0x7fff) {
		buf.WriteByte(0xd1)
		binary.Write(&buf, binary.BigEndian, int16(value))
	} else if (value >= -0x80000000) && (value <= 0x7fffffff) {
		buf.WriteByte(0xd2)
		binary.Write(&buf, binary.BigEndian, int32(value))
	} else if (value >= -0x8000000000000000) && (value <= 0x7fffffffffffffff) {
		buf.WriteByte(0xd3)
		binary.Write(&buf, binary.BigEndian, int64(value))
	}

	_, err := w.Write(buf.Bytes())
	return err
}

// PackUint writes an uint type data to writer.
func PackUint(w io.Writer, value uint64) error {
	var buf bytes.Buffer

	if value <= 0xff {
		buf.WriteByte(0xcc)
		buf.WriteByte(uint8(value))
	} else if value <= 0xffff {
		buf.WriteByte(0xcd)
		binary.Write(&buf, binary.BigEndian, uint16(value))
	} else if value <= 0xffffffff {
		buf.WriteByte(0xce)
		binary.Write(&buf, binary.BigEndian, uint32(value))
	} else if value <= 0xffffffffffffffff {
		buf.WriteByte(0xcf)
		binary.Write(&buf, binary.BigEndian, uint64(value))
	}

	_, err := w.Write(buf.Bytes())
	return err
}

// PackFloat32 writes a float32 data to writer.
func PackFloat32(w io.Writer, value float32) error {
	var buf bytes.Buffer

	buf.WriteByte(0xca)
	binary.Write(&buf, binary.BigEndian, math.Float32bits(value))

	_, err := w.Write(buf.Bytes())
	return err
}

// PackFloat64 writes a float64 data to writer.
func PackFloat64(w io.Writer, value float64) error {
	var buf bytes.Buffer

	buf.WriteByte(0xcb)
	binary.Write(&buf, binary.BigEndian, math.Float64bits(value))

	_, err := w.Write(buf.Bytes())
	return err
}

// PackString writes a string data to writer.
func PackString(w io.Writer, value string) error {
	var buf bytes.Buffer

	len := len(value)
	if len <= 0x1f {
		buf.WriteByte(0xa0 | uint8(len))
		buf.WriteString(value)
	} else if len <= 0xff {
		buf.WriteByte(0xd9)
		buf.WriteByte(uint8(len))
		buf.WriteString(value)
	} else if len <= 0xffff {
		buf.WriteByte(0xda)
		binary.Write(&buf, binary.BigEndian, uint16(len))
		buf.WriteString(value)
	} else if len <= 0xffffffff {
		buf.WriteByte(0xdb)
		binary.Write(&buf, binary.BigEndian, uint32(len))
		buf.WriteString(value)
	}

	_, err := w.Write(buf.Bytes())
	return err
}

// PackStruct writes a struct data to writer.
func PackStruct(w io.Writer, value interface{}) error {
	var headBuf bytes.Buffer
	var dataBuf bytes.Buffer
	var err error

	structType := reflect.TypeOf(value)
	structValue := reflect.ValueOf(value)
	structNumField := structType.NumField()

	numField := uint32(0)
	for inx := 0; inx < structNumField; inx++ {
		var fp FieldProp

		field := structType.Field(inx)
		fp.parseTag(field)
		if fp.Skip {
			continue
		}

		fieldValue := structValue.Field(inx)
		if fp.OmitEmpty {
			if fieldValue.Interface() == reflect.Zero(fieldValue.Type()).Interface() {
				continue
			}
		}

		if err = PackString(&dataBuf, fp.Name); err != nil {
			return err
		}

		if fp.String {
			switch fieldValue.Kind() {
			case reflect.Bool,
				reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
				reflect.Float32, reflect.Float64:
				err = PackString(&dataBuf, fmt.Sprintf("%v", fieldValue.Interface()))
			default:
				err = PackValue(&dataBuf, fieldValue.Interface())
			}
		} else {
			err = PackValue(&dataBuf, fieldValue.Interface())
		}
		if err != nil {
			return err
		}

		numField++
	}

	if numField <= 0x0f {
		headBuf.WriteByte(0x80 | uint8(numField))
	} else if numField <= 0xffff {
		headBuf.WriteByte(0xde)
		binary.Write(&headBuf, binary.BigEndian, uint16(numField))
	} else if numField <= 0xffffffff {
		headBuf.WriteByte(0xdf)
		binary.Write(&headBuf, binary.BigEndian, uint32(numField))
	}

	if _, err = w.Write(headBuf.Bytes()); err != nil {
		return err
	}

	_, err = w.Write(dataBuf.Bytes())
	return err
}

// PackMap writes a map data to writer.
func PackMap(w io.Writer, value interface{}) error {
	var buf bytes.Buffer

	m := reflect.ValueOf(value)
	mapSize := m.Len()
	if mapSize <= 0x0f {
		buf.WriteByte(0x80 | uint8(mapSize))
	} else if mapSize <= 0xffff {
		buf.WriteByte(0xde)
		binary.Write(&buf, binary.BigEndian, uint16(mapSize))
	} else if mapSize <= 0xffffffff {
		buf.WriteByte(0xdf)
		binary.Write(&buf, binary.BigEndian, uint32(mapSize))
	}

	for _, key := range m.MapKeys() {
		err := PackValue(&buf, key.Interface())
		if err != nil {
			return err
		}
		err = PackValue(&buf, m.MapIndex(key).Interface())
		if err != nil {
			return err
		}
	}

	_, err := w.Write(buf.Bytes())
	return err
}

// PackArray writes a array data to writer.
func PackArray(w io.Writer, value interface{}) error {
	var buf bytes.Buffer

	if reflect.TypeOf(value).Elem().Kind() == reflect.Uint8 { // for []byte
		b := reflect.ValueOf(value)
		arraySize := b.Len()
		if arraySize <= 0xff {
			buf.WriteByte(0xc4)
			binary.Write(&buf, binary.BigEndian, uint8(arraySize))
		} else if arraySize <= 0xffff {
			buf.WriteByte(0xc5)
			binary.Write(&buf, binary.BigEndian, uint16(arraySize))
		} else if arraySize <= 0xffffffff {
			buf.WriteByte(0xc6)
			binary.Write(&buf, binary.BigEndian, uint32(arraySize))
		}

		_, err := w.Write(buf.Bytes())
		if err != nil {
			return err
		}

		_, err = w.Write(value.([]byte))
		return err
	}

	a := reflect.ValueOf(value)
	arraySize := a.Len()
	if arraySize <= 0x0f {
		buf.WriteByte(0x90 | uint8(arraySize))
	} else if arraySize <= 0xffff {
		buf.WriteByte(0xdc)
		binary.Write(&buf, binary.BigEndian, uint16(arraySize))
	} else if arraySize <= 0xffffffff {
		buf.WriteByte(0xdd)
		binary.Write(&buf, binary.BigEndian, uint32(arraySize))
	}

	for inx := 0; inx < a.Len(); inx++ {
		PackValue(&buf, a.Index(inx).Interface())
	}

	_, err := w.Write(buf.Bytes())
	return err
}

// PackPtr writes the data pointed by ptr to writer.
func PackPtr(w io.Writer, value interface{}) error {
	return PackValue(w, reflect.ValueOf(value).Elem().Interface())
}
