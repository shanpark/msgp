package msgp

import (
	"bytes"
	"fmt"
)

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
	PackInt(&buf, 0x7fff)
	PackInt(&buf, 0x7fffffff)
	PackInt(&buf, 0x7fffffffffffffff)

	fmt.Printf("% x\n", buf.Bytes())

	buf.Reset()

	PackInt(&buf, -0x80)
	PackInt(&buf, -0x8000)
	PackInt(&buf, -0x80000000)
	PackInt(&buf, -0x8000000000000000)

	fmt.Printf("% x\n", buf.Bytes())

	// Output:
	// 7f d1 7f ff d2 7f ff ff ff d3 7f ff ff ff ff ff ff ff
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

type structType struct {
	AAA string
	BBB int32
}

func ExamplePackStruct() {
	var buf bytes.Buffer
	var str = structType{"1234567890", 0xff}

	PackStruct(&buf, str)
	fmt.Printf("% x\n", buf.Bytes())

	// Output:
	// 82 a3 41 41 41 aa 31 32 33 34 35 36 37 38 39 30 a3 42 42 42 d1 00 ff
}

func ExamplePackMap() {
	var buf bytes.Buffer
	m := map[string]int{"aaa": 1, "bbb": 2}

	PackMap(&buf, m)
	fmt.Printf("% x\n", buf.Bytes())

	// Output:
	// 82 a3 61 61 61 01 a3 62 62 62 02
}
