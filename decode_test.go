package msgp

import (
	"bytes"
	"fmt"
)

func ExampleUnpackValue_nil() {
	var err error
	var buf bytes.Buffer

	var s string
	var ba [1]byte
	var bs []byte
	var m map[string]int
	var pb *bool
	var unknown interface{}

	if IsLittleEndian() {
		fmt.Println("Little")
	} else {
		fmt.Println("Big")
	}

	PackNil(&buf)
	err = UnpackValue(&buf, &s)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", s)
	}

	PackNil(&buf)
	err = UnpackValue(&buf, &ba)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", ba)
	}

	PackNil(&buf)
	err = UnpackValue(&buf, &bs)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", bs)
	}

	PackNil(&buf)
	err = UnpackValue(&buf, &m)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", m)
	}

	PackNil(&buf)
	err = UnpackValue(&buf, &pb)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", *pb)
	}

	PackNil(&buf)
	err = UnpackValue(&buf, &unknown)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", unknown)
	}

	// Output:
	//
	// [0]
	// []
	// map[]
	// false
	// msgp: specified type[interface] is not supported
}

func ExampleUnpackValue_bool() {
	var err error
	var buf bytes.Buffer
	var val bool

	PackBool(&buf, true)
	err = UnpackValue(&buf, &val)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", val)
	}
	PackBool(&buf, false)
	err = UnpackValue(&buf, &val)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", val)
	}
	PackNil(&buf)
	err = UnpackValue(&buf, &val)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", val)
	}

	// Output:
	// true
	// false
	// false
}

func ExampleUnpackValue_int() {
	var err error
	var buf bytes.Buffer
	var i int
	var i8 int8
	var i16 int16
	var i32 int32
	var i64 int64
	var ui uint
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
	err = UnpackValue(&buf, &ui)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%x\n", ui)
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
	// 0
	// 0
}

func ExampleUnpackValue_float() {
	var err error
	var buf bytes.Buffer
	var f32 float32
	var f64 float64

	PackFloat32(&buf, 3.14)
	err = UnpackValue(&buf, &f32)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", f32)
	}

	PackFloat64(&buf, 3.14)
	err = UnpackValue(&buf, &f64)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", f64)
	}

	PackNil(&buf)
	err = UnpackValue(&buf, &f64)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", f64)
	}

	// Output:
	// 3.14
	// 3.14
	// 0
}
