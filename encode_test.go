package msgp

import (
	"bytes"
	"fmt"
)

func ExamplePackValue() {
	var buf bytes.Buffer

	// nil
	PackValue(&buf, nil)
	fmt.Printf("% x\n", buf.Bytes())

	// bool
	buf.Reset()

	PackValue(&buf, true)
	PackValue(&buf, false)
	fmt.Printf("% x\n", buf.Bytes())

	// integers
	buf.Reset()

	PackValue(&buf, 0x7f)
	PackValue(&buf, 0x80)
	PackValue(&buf, 0x7fff)
	PackValue(&buf, 0x8000)
	PackValue(&buf, 0x7fffffff)
	PackValue(&buf, 0x80000000)
	PackValue(&buf, 0x7fffffffffffffff)
	fmt.Printf("% x\n", buf.Bytes())

	buf.Reset()

	PackValue(&buf, -0x80)
	PackValue(&buf, -0x8000)
	PackValue(&buf, -0x80000000)
	PackValue(&buf, -0x8000000000000000)
	fmt.Printf("% x\n", buf.Bytes())

	buf.Reset()

	PackValue(&buf, 0xff)
	PackValue(&buf, 0xffff)
	PackValue(&buf, 0xffffffff)
	PackValue(&buf, uint64(0xffffffffffffffff))
	fmt.Printf("% x\n", buf.Bytes())

	// floats
	buf.Reset()

	PackValue(&buf, float32(3.14))
	fmt.Printf("% x\n", buf.Bytes())

	buf.Reset()

	PackValue(&buf, float64(3.14))
	fmt.Printf("% x\n", buf.Bytes())

	// string
	buf.Reset()

	PackValue(&buf, "1234567890123456789012345678901")
	fmt.Printf("% x\n", buf.Bytes())

	// array
	buf.Reset()

	a := []string{"aaa", "bbb", "ccc"}
	PackValue(&buf, a)
	fmt.Printf("% x\n", buf.Bytes())

	// map
	buf.Reset()

	m := map[string]int{"aaa": 1}
	PackValue(&buf, m)
	fmt.Printf("% x\n", buf.Bytes())

	// struct
	buf.Reset()

	type myStruct struct {
		AAA string
		BBB int
		CCC string `msgp:"ccc"`
		DDD int    `msgp:"-"`
		FFF int32  `msgp:"-,omitempty"`
		GGG int32  `msgp:",omitempty"`
		HHH int32  `msgp:",string"`
		III int32  `msgp:",omitempty"`
	}
	st := myStruct{"1234567890", 0xff, "12345", 0x11, 0x22, 0x33, 100, 0}
	PackValue(&buf, st)
	fmt.Printf("% x\n", buf.Bytes())

	// pointer
	buf.Reset()

	PackValue(&buf, &a)
	fmt.Printf("% x\n", buf.Bytes())

	// Output:
	// c0
	// c3 c2
	// 7f cc 80 d1 7f ff cd 80 00 d2 7f ff ff ff ce 80 00 00 00 d3 7f ff ff ff ff ff ff ff
	// d0 80 d1 80 00 d2 80 00 00 00 d3 80 00 00 00 00 00 00 00
	// cc ff cd ff ff ce ff ff ff ff cf ff ff ff ff ff ff ff ff
	// ca 40 48 f5 c3
	// cb 40 09 1e b8 51 eb 85 1f
	// bf 31 32 33 34 35 36 37 38 39 30 31 32 33 34 35 36 37 38 39 30 31 32 33 34 35 36 37 38 39 30 31
	// 93 a3 61 61 61 a3 62 62 62 a3 63 63 63
	// 81 a3 61 61 61 01
	// 86 a3 41 41 41 aa 31 32 33 34 35 36 37 38 39 30 a3 42 42 42 cc ff a3 63 63 63 a5 31 32 33 34 35 a1 5f 22 a3 47 47 47 33 a3 48 48 48 a3 31 30 30
	// 93 a3 61 61 61 a3 62 62 62 a3 63 63 63
}

func ExamplePackValue_nil() {
	var buf bytes.Buffer

	PackNil(&buf)
	fmt.Printf("% x\n", buf.Bytes())

	// Output:
	// c0
}

func ExamplePackValue_bool() {
	var buf bytes.Buffer

	PackBool(&buf, true)
	PackBool(&buf, false)
	fmt.Printf("% x\n", buf.Bytes())

	// Output:
	// c3 c2
}

func ExamplePackValue_int() {
	var buf bytes.Buffer

	PackInt(&buf, 0x7f)
	PackInt(&buf, 0x80)
	PackInt(&buf, 0x7fff)
	PackInt(&buf, 0x8000)
	PackInt(&buf, 0x7fffffff)
	PackInt(&buf, 0x80000000)
	PackInt(&buf, 0x7fffffffffffffff)
	fmt.Printf("% x\n", buf.Bytes())

	buf.Reset()

	PackInt(&buf, -0x80)
	PackInt(&buf, -0x8000)
	PackInt(&buf, -0x80000000)
	PackInt(&buf, -0x8000000000000000)
	fmt.Printf("% x\n", buf.Bytes())

	// Output:
	// 7f cc 80 d1 7f ff cd 80 00 d2 7f ff ff ff ce 80 00 00 00 d3 7f ff ff ff ff ff ff ff
	// d0 80 d1 80 00 d2 80 00 00 00 d3 80 00 00 00 00 00 00 00
}

func ExamplePackValue_uint() {
	var buf bytes.Buffer

	PackUint(&buf, 0xff)
	PackUint(&buf, 0xffff)
	PackUint(&buf, 0xffffffff)
	PackUint(&buf, 0xffffffffffffffff)
	PackUint(&buf, 0x1ff)
	fmt.Printf("% x\n", buf.Bytes())

	// Output:
	// cc ff cd ff ff ce ff ff ff ff cf ff ff ff ff ff ff ff ff cd 01 ff
}

func ExamplePackValue_float32() {
	var buf bytes.Buffer

	PackFloat32(&buf, float32(3.14))
	fmt.Printf("% x\n", buf.Bytes())

	// Output:
	// ca 40 48 f5 c3
}

func ExamplePackValue_float64() {
	var buf bytes.Buffer

	PackFloat64(&buf, float64(3.14))
	fmt.Printf("% x\n", buf.Bytes())

	// Output:
	// cb 40 09 1e b8 51 eb 85 1f
}

func ExamplePackValue_string() {
	var buf bytes.Buffer

	PackString(&buf, "1234567890123456789012345678901")
	fmt.Printf("% x\n", buf.Bytes())

	buf.Reset()

	PackString(&buf, "12345678901234567890123456789012")
	fmt.Printf("% x\n", buf.Bytes())

	// Output:
	// bf 31 32 33 34 35 36 37 38 39 30 31 32 33 34 35 36 37 38 39 30 31 32 33 34 35 36 37 38 39 30 31
	// d9 20 31 32 33 34 35 36 37 38 39 30 31 32 33 34 35 36 37 38 39 30 31 32 33 34 35 36 37 38 39 30 31 32
}

func ExamplePackValue_map() {
	var buf bytes.Buffer
	m := map[string]int{"aaa": 1}

	PackMap(&buf, m)
	fmt.Printf("% x\n", buf.Bytes())

	// Output:
	// 81 a3 61 61 61 01
}

func ExamplePackValue_struct() {
	type myStruct struct {
		AAA string
		BBB int
		CCC string `msgp:"ccc"`
		DDD int    `msgp:"-"`
		FFF int32  `msgp:"-,omitempty"`
		GGG int32  `msgp:",omitempty"`
		HHH int32  `msgp:",string"`
		III int32  `msgp:",omitempty"`
	}

	var buf bytes.Buffer
	var str = myStruct{"1234567890", 0xff, "12345", 0x11, 0x22, 0x33, 100, 0}

	PackStruct(&buf, str)
	fmt.Printf("% x\n", buf.Bytes())

	// Output:
	// 86 a3 41 41 41 aa 31 32 33 34 35 36 37 38 39 30 a3 42 42 42 cc ff a3 63 63 63 a5 31 32 33 34 35 a1 5f 22 a3 47 47 47 33 a3 48 48 48 a3 31 30 30
}

func ExamplePackValue_array() {
	var buf bytes.Buffer
	a := []string{"aaa", "bbb", "ccc"}
	b := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	c := []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

	PackArray(&buf, a)
	fmt.Printf("% x\n", buf.Bytes())

	buf.Reset()

	PackArray(&buf, b)
	fmt.Printf("% x\n", buf.Bytes())

	buf.Reset()

	PackArray(&buf, c)
	fmt.Printf("% x\n", buf.Bytes())

	// Output:
	// 93 a3 61 61 61 a3 62 62 62 a3 63 63 63
	// c4 20 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f 10 11 12 13 14 15 16 17 18 19 1a 1b 1c 1d 1e 1f 20
	// c4 20 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f 10 11 12 13 14 15 16 17 18 19 1a 1b 1c 1d 1e 1f 20
}

func ExamplePackValue_pointer() {
	type myStruct struct {
		AAA string
		BBB int
		CCC string `msgp:"ccc"`
		DDD int    `msgp:"-"`
		FFF int32  `msgp:"-,omitempty"`
		GGG int32  `msgp:",omitempty"`
		HHH int32  `msgp:",string"`
		III int32  `msgp:",omitempty"`
	}

	var buf bytes.Buffer
	a := []string{"aaa", "bbb", "ccc"}
	b := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	m := map[string]int{"aaa": 1}
	var str = myStruct{"1234567890", 0xff, "12345", 0x11, 0x22, 0x33, 100, 0}

	PackPtr(&buf, &a)
	fmt.Printf("% x\n", buf.Bytes())

	buf.Reset()

	PackPtr(&buf, &b)
	fmt.Printf("% x\n", buf.Bytes())

	buf.Reset()

	PackPtr(&buf, &m)
	fmt.Printf("% x\n", buf.Bytes())

	buf.Reset()

	PackPtr(&buf, &str)
	fmt.Printf("% x\n", buf.Bytes())

	// Output:
	// 93 a3 61 61 61 a3 62 62 62 a3 63 63 63
	// c4 20 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f 10 11 12 13 14 15 16 17 18 19 1a 1b 1c 1d 1e 1f 20
	// 81 a3 61 61 61 01
	// 86 a3 41 41 41 aa 31 32 33 34 35 36 37 38 39 30 a3 42 42 42 cc ff a3 63 63 63 a5 31 32 33 34 35 a1 5f 22 a3 47 47 47 33 a3 48 48 48 a3 31 30 30
}
