package msgp

import (
	"bytes"
	"fmt"
)

type structType struct {
	AAA string
	BBB int
	CCC string `msgp:"ccc"`
	DDD int    `msgp:"-"`
	FFF int32  `msgp:"-,omitempty"`
	GGG int32  `msgp:",omitempty"`
	HHH int32  `msgp:",string"`
	III int32  `msgp:",omitempty"`
}

func ExamplePackNil() {
	var buf bytes.Buffer

	PackNil(&buf)

	fmt.Printf("% x\n", buf.Bytes())
	// Output:
	// c0
}

func ExamplePackBool() {
	var buf bytes.Buffer

	PackBool(&buf, true)
	PackBool(&buf, false)

	fmt.Printf("% x\n", buf.Bytes())
	// Output:
	// c3 c2
}

func ExamplePackInt() {
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

func ExamplePackUint() {
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

func ExamplePackFloat32() {
	var buf bytes.Buffer

	PackFloat32(&buf, float32(3.14))

	fmt.Printf("% x\n", buf.Bytes())

	// Output:
	// ca 40 48 f5 c3
}

func ExamplePackFloat64() {
	var buf bytes.Buffer

	PackFloat64(&buf, float64(3.14))

	fmt.Printf("% x\n", buf.Bytes())

	// Output:
	// cb 40 09 1e b8 51 eb 85 1f
}

func ExamplePackString() {
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

func ExamplePackStruct() {
	var buf bytes.Buffer
	var str = structType{"1234567890", 0xff, "12345", 0x11, 0x22, 0x33, 100, 0}

	PackStruct(&buf, str)
	fmt.Printf("% x\n", buf.Bytes())

	// Output:
	// 86 a3 41 41 41 aa 31 32 33 34 35 36 37 38 39 30 a3 42 42 42 d1 00 ff a3 63 63 63 a5 31 32 33 34 35 a1 5f 22 a3 47 47 47 33 a3 48 48 48 a3 31 30 30
}

func ExamplePackMap() {
	var buf bytes.Buffer
	m := map[string]int{"aaa": 1}

	PackMap(&buf, m)
	fmt.Printf("% x\n", buf.Bytes())

	// Output:
	// 81 a3 61 61 61 01
}

func ExamplePackArray() {
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

func ExamplePackPtr() {
	var buf bytes.Buffer
	a := []string{"aaa", "bbb", "ccc"}
	b := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	m := map[string]int{"aaa": 1}
	var str = structType{"1234567890", 0xff, "12345", 0x11, 0x22, 0x33, 100, 0}

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
	// 86 a3 41 41 41 aa 31 32 33 34 35 36 37 38 39 30 a3 42 42 42 d1 00 ff a3 63 63 63 a5 31 32 33 34 35 a1 5f 22 a3 47 47 47 33 a3 48 48 48 a3 31 30 30
}
