package msgp

import (
	"bytes"
	"fmt"
)

func ExampleUnpackValue_nil() {
	var buf bytes.Buffer
	var val interface{}

	PackNil(&buf)
	UnpackValue(&buf, &val)
	fmt.Printf("%v\n", val)

	// Output:
	// <nil>
}

func ExampleUnpackValue_int() {
	var err error
	var buf bytes.Buffer
	var i int
	var i8 int8
	var i16 int16
	var i32 int32
	var i64 int64
	var ui8 uint8
	var ui16 uint16
	var ui32 uint32
	var ui64 uint64

	PackInt(&buf, 0x7f)
	UnpackValue(&buf, &i8)
	fmt.Printf("%x\n", i8)
	PackInt(&buf, 0x7fff)
	UnpackValue(&buf, &i16)
	fmt.Printf("%x\n", i16)
	PackInt(&buf, 0x7fffffff)
	UnpackValue(&buf, &i32)
	fmt.Printf("%x\n", i32)
	PackInt(&buf, 0x7fffffffffffffff)
	UnpackValue(&buf, &i64)
	fmt.Printf("%x\n", i64)

	PackUint(&buf, 0xff)
	UnpackValue(&buf, &ui8)
	fmt.Printf("%x\n", ui8)
	PackUint(&buf, 0xffff)
	UnpackValue(&buf, &ui16)
	fmt.Printf("%x\n", ui16)
	PackUint(&buf, 0xffffffff)
	UnpackValue(&buf, &ui32)
	fmt.Printf("%x\n", ui32)
	PackUint(&buf, 0xffffffffffffffff)
	UnpackValue(&buf, &ui64)
	fmt.Printf("%x\n", ui64)

	PackNil(&buf)
	err = UnpackValue(&buf, &i)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%x\n", i)
	}

	PackNil(&buf)
	err = UnpackValue(&buf, &i8)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%x\n", i)
	}

	// Output:
	// 7f
	// 7fff
	// 7fffffff
	// 7fffffffffffffff
	// ff
	// ffff
	// ffffffff
	// ffffffffffffffff
	// msgp: unpacked value[<nil>] is not assignable to integer type
	// msgp: unpacked value[<nil>] is not assignable to integer type
}
