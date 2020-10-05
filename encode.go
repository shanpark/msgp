package msgp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
)

// Pack writes a value to the io.Writer.
// It is recommended to use this function for all types.
func Pack(w io.Writer, value interface{}) error {
	var err error

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
	case reflect.Array, reflect.Slice:
		err = PackArray(w, value)
	case reflect.Map:
		err = PackMap(w, value)
	case reflect.Struct:
		err = PackStruct(w, value)
	case reflect.Ptr:
		err = PackPtr(w, value)
	default:
		err = errors.New("msgp: unsupported type value")
	}

	return err
}

// PackNil writes a nil value to the io.Writer.
func PackNil(w io.Writer) error {
	_, err := w.Write([]byte{0xc0})
	return err
}

// PackBool writes a bool value to the io.Writer.
func PackBool(w io.Writer, value bool) error {
	var err error

	if value {
		_, err = w.Write([]byte{0xc3})
	} else {
		_, err = w.Write([]byte{0xc2})
	}

	return err
}

// PackInt writes an integer value to the io.Writer.
func PackInt(w io.Writer, value int64) error {
	var err error
	var buf bytes.Buffer

	if value >= 0 {
		if value <= 0x7f { // int
			if err = buf.WriteByte(byte(value)); err != nil {
				return err
			}
		} else if value <= 0xff { // uint
			if err = buf.WriteByte(0xcc); err != nil {
				return err
			}
			if err = buf.WriteByte(uint8(value)); err != nil {
				return err
			}
		} else if value <= 0x7fff { // int
			if err = buf.WriteByte(0xd1); err != nil {
				return err
			}
			if err = binary.Write(&buf, binary.BigEndian, int16(value)); err != nil {
				return err
			}
		} else if value <= 0xffff { // uint
			if err = buf.WriteByte(0xcd); err != nil {
				return err
			}
			if err = binary.Write(&buf, binary.BigEndian, uint16(value)); err != nil {
				return err
			}
		} else if value <= 0x7fffffff { // int
			if err = buf.WriteByte(0xd2); err != nil {
				return err
			}
			if err = binary.Write(&buf, binary.BigEndian, int32(value)); err != nil {
				return err
			}
		} else if value <= 0xffffffff { // uint
			if err = buf.WriteByte(0xce); err != nil {
				return err
			}
			if err = binary.Write(&buf, binary.BigEndian, uint32(value)); err != nil {
				return err
			}
		} else { // int
			if err = buf.WriteByte(0xd3); err != nil {
				return err
			}
			if err = binary.Write(&buf, binary.BigEndian, int64(value)); err != nil {
				return err
			}
		}
	} else {
		if value >= -32 {
			if err = buf.WriteByte(byte(value)); err != nil {
				return err
			}
		} else if value >= -0x80 {
			if err = buf.WriteByte(0xd0); err != nil {
				return err
			}
			if err = buf.WriteByte(byte(value)); err != nil {
				return err
			}
		} else if value >= -0x8000 {
			if err = buf.WriteByte(0xd1); err != nil {
				return err
			}
			if err = binary.Write(&buf, binary.BigEndian, int16(value)); err != nil {
				return err
			}
		} else if value >= -0x80000000 {
			if err = buf.WriteByte(0xd2); err != nil {
				return err
			}
			if err = binary.Write(&buf, binary.BigEndian, int32(value)); err != nil {
				return err
			}
		} else {
			if err = buf.WriteByte(0xd3); err != nil {
				return err
			}
			if err = binary.Write(&buf, binary.BigEndian, int64(value)); err != nil {
				return err
			}
		}
	}

	_, err = w.Write(buf.Bytes())
	return err
}

// PackUint writes an unsigned integer value to the io.Writer.
func PackUint(w io.Writer, value uint64) error {
	var err error
	var buf bytes.Buffer

	if value <= 0xff {
		if err = buf.WriteByte(0xcc); err != nil {
			return err
		}
		if err = buf.WriteByte(uint8(value)); err != nil {
			return err
		}
	} else if value <= 0xffff {
		if err = buf.WriteByte(0xcd); err != nil {
			return err
		}
		if err = binary.Write(&buf, binary.BigEndian, uint16(value)); err != nil {
			return err
		}
	} else if value <= 0xffffffff {
		if err = buf.WriteByte(0xce); err != nil {
			return err
		}
		if err = binary.Write(&buf, binary.BigEndian, uint32(value)); err != nil {
			return err
		}
	} else {
		if err = buf.WriteByte(0xcf); err != nil {
			return err
		}
		if err = binary.Write(&buf, binary.BigEndian, value); err != nil {
			return err
		}
	}

	_, err = w.Write(buf.Bytes())
	return err
}

// PackFloat32 writes a float32 value to the io.Writer.
func PackFloat32(w io.Writer, value float32) error {
	var err error
	var buf bytes.Buffer

	if err = buf.WriteByte(0xca); err != nil {
		return err
	}
	if err = binary.Write(&buf, binary.BigEndian, math.Float32bits(value)); err != nil {
		return err
	}

	_, err = w.Write(buf.Bytes())
	return err
}

// PackFloat64 writes a float64 value to the io.Writer.
func PackFloat64(w io.Writer, value float64) error {
	var err error
	var buf bytes.Buffer

	if err = buf.WriteByte(0xcb); err != nil {
		return err
	}
	if err = binary.Write(&buf, binary.BigEndian, math.Float64bits(value)); err != nil {
		return err
	}

	_, err = w.Write(buf.Bytes())
	return err
}

// PackString writes a string value to the io.Writer.
func PackString(w io.Writer, value string) error {
	var err error
	var buf bytes.Buffer

	len := len(value)
	if len <= 0x1f {
		if err = buf.WriteByte(0xa0 | uint8(len)); err != nil {
			return err
		}
		if _, err = buf.WriteString(value); err != nil {
			return err
		}
	} else if len <= 0xff {
		if err = buf.WriteByte(0xd9); err != nil {
			return err
		}
		if err = buf.WriteByte(uint8(len)); err != nil {
			return err
		}
		if _, err = buf.WriteString(value); err != nil {
			return err
		}
	} else if len <= 0xffff {
		if err = buf.WriteByte(0xda); err != nil {
			return err
		}
		if err = binary.Write(&buf, binary.BigEndian, uint16(len)); err != nil {
			return err
		}
		if _, err = buf.WriteString(value); err != nil {
			return err
		}
	} else if len <= 0xffffffff {
		if err = buf.WriteByte(0xdb); err != nil {
			return err
		}
		if err = binary.Write(&buf, binary.BigEndian, uint32(len)); err != nil {
			return err
		}
		if _, err = buf.WriteString(value); err != nil {
			return err
		}
	} else {
		return errors.New("msgp: try to pack too long string")
	}

	_, err = w.Write(buf.Bytes())
	return err
}

// PackArray writes an array to the io.Writer.
func PackArray(w io.Writer, value interface{}) error {
	var err error
	var buf bytes.Buffer

	if reflect.TypeOf(value).Elem().Kind() == reflect.Uint8 { // for []byte
		b := reflect.ValueOf(value)
		arraySize := b.Len()
		if arraySize <= 0xff {
			if err = buf.WriteByte(0xc4); err != nil {
				return err
			}
			if err = binary.Write(&buf, binary.BigEndian, uint8(arraySize)); err != nil {
				return err
			}
		} else if arraySize <= 0xffff {
			if err = buf.WriteByte(0xc5); err != nil {
				return err
			}
			if err = binary.Write(&buf, binary.BigEndian, uint16(arraySize)); err != nil {
				return err
			}
		} else if arraySize <= 0xffffffff {
			if err = buf.WriteByte(0xc6); err != nil {
				return err
			}
			if err = binary.Write(&buf, binary.BigEndian, uint32(arraySize)); err != nil {
				return err
			}
		}

		if _, err = w.Write(buf.Bytes()); err != nil {
			return err
		}

		_, err = w.Write(value.([]byte))
		return err
	}

	a := reflect.ValueOf(value)
	arraySize := a.Len()
	if arraySize <= 0x0f {
		if err = buf.WriteByte(0x90 | uint8(arraySize)); err != nil {
			return err
		}
	} else if arraySize <= 0xffff {
		if err = buf.WriteByte(0xdc); err != nil {
			return err
		}
		if err = binary.Write(&buf, binary.BigEndian, uint16(arraySize)); err != nil {
			return err
		}
	} else if arraySize <= 0xffffffff {
		if err = buf.WriteByte(0xdd); err != nil {
			return err
		}
		if err = binary.Write(&buf, binary.BigEndian, uint32(arraySize)); err != nil {
			return err
		}
	}

	for inx := 0; inx < a.Len(); inx++ {
		if err = Pack(&buf, a.Index(inx).Interface()); err != nil {
			return err
		}
	}

	_, err = w.Write(buf.Bytes())
	return err
}

// PackMap writes a map to the io.Writer.
func PackMap(w io.Writer, value interface{}) error {
	var err error
	var buf bytes.Buffer

	m := reflect.ValueOf(value)
	mapSize := m.Len()
	if mapSize <= 0x0f {
		if err = buf.WriteByte(0x80 | uint8(mapSize)); err != nil {
			return err
		}
	} else if mapSize <= 0xffff {
		if err = buf.WriteByte(0xde); err != nil {
			return err
		}
		if err = binary.Write(&buf, binary.BigEndian, uint16(mapSize)); err != nil {
			return err
		}
	} else if mapSize <= 0xffffffff {
		if err = buf.WriteByte(0xdf); err != nil {
			return err
		}
		if err = binary.Write(&buf, binary.BigEndian, uint32(mapSize)); err != nil {
			return err
		}
	}

	for _, key := range m.MapKeys() {
		if err = Pack(&buf, key.Interface()); err != nil {
			return err
		}
		if err = Pack(&buf, m.MapIndex(key).Interface()); err != nil {
			return err
		}
	}

	_, err = w.Write(buf.Bytes())
	return err
}

// PackStruct writes a struct value to the io.Writer.
// The struct value is serialized as a map[string]interface{}.
func PackStruct(w io.Writer, value interface{}) error {
	var err error
	var headBuf bytes.Buffer
	var dataBuf bytes.Buffer

	structTyp := reflect.TypeOf(value)
	structVal := reflect.ValueOf(value)
	structNumField := structTyp.NumField()

	numField := uint32(0)
	for inx := 0; inx < structNumField; inx++ {
		var fp FieldProps

		field := structTyp.Field(inx)
		fp.parseTag(field)
		if fp.Skip {
			continue
		}

		fieldValue := structVal.Field(inx)
		if fp.OmitEmpty {
			if fieldValue.Interface() == reflect.Zero(fieldValue.Type()).Interface() {
				continue
			}
		}

		if err = PackString(&dataBuf, fp.Name); err != nil {
			return err
		}

		if fp.String {
			if fieldValue.Interface() == nil {
				err = PackString(&dataBuf, "nil")
			} else {
			Loop:
				for {
					switch fieldValue.Kind() {
					case reflect.Bool,
						reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
						reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
						reflect.Float32, reflect.Float64:
						err = PackString(&dataBuf, fmt.Sprintf("%v", fieldValue.Interface()))
						break Loop
					case reflect.Ptr:
						fieldValue = fieldValue.Elem()
					default:
						err = fmt.Errorf("msgp: cannot pack Go struct field %v.%s of type %v into string", structTyp, field.Name, field.Type)
						break Loop
					}
				}
			}
		} else {
			err = Pack(&dataBuf, fieldValue.Interface())
		}
		if err != nil {
			return err
		}

		numField++
	}

	if numField <= 0x0f {
		if err = headBuf.WriteByte(0x80 | uint8(numField)); err != nil {
			return err
		}
	} else if numField <= 0xffff {
		if err = headBuf.WriteByte(0xde); err != nil {
			return err
		}
		if err = binary.Write(&headBuf, binary.BigEndian, uint16(numField)); err != nil {
			return err
		}
	} else if numField <= 0xffffffff {
		if err = headBuf.WriteByte(0xdf); err != nil {
			return err
		}
		if err = binary.Write(&headBuf, binary.BigEndian, uint32(numField)); err != nil {
			return err
		}
	}

	if _, err = w.Write(headBuf.Bytes()); err != nil {
		return err
	}

	_, err = w.Write(dataBuf.Bytes())
	return err
}

// PackPtr writes a value pointed by ptr to the io.Writer.
func PackPtr(w io.Writer, ptr interface{}) error {
	return Pack(w, reflect.ValueOf(ptr).Elem().Interface())
}
