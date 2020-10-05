package msgp

import (
	"bytes"
	"fmt"
)

func ExampleUnpack() {
	var err error
	var buf bytes.Buffer
	var b bool
	var i8 int8
	var ui16 uint16
	var f32 float32
	var str string
	var ia []int
	var msi map[string]int
	var ptr *string
	var unknown interface{}

	Pack(&buf, true)
	err = Unpack(&buf, &b)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", b)
	}

	Pack(&buf, 0x7f)
	err = Unpack(&buf, &i8)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%x\n", i8)
	}

	Pack(&buf, 0x11)
	err = Unpack(&buf, &ui16)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%x\n", ui16)
	}

	Pack(&buf, 3.14)
	err = Unpack(&buf, &f32)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", f32)
	}

	Pack(&buf, "test string")
	err = Unpack(&buf, &str)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", str)
	}

	Pack(&buf, []int8{1, 2, 3, 4, 5, 6, 7, 8, 9})
	err = Unpack(&buf, &ia)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", ia)
	}

	Pack(&buf, map[string]int16{"a": 1, "b": 2})
	err = Unpack(&buf, &msi)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", msi)
	}

	Pack(&buf, "some text")
	err = Unpack(&buf, &ptr)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", *ptr)
	}

	Pack(&buf, "string in interface{}")
	err = Unpack(&buf, &unknown)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", unknown)
	}

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
	var st myStruct

	Pack(&buf, myStruct{"1234567890", 255, "12345", 17, 34, 51, 100, 0})
	err = Unpack(&buf, &st)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", st)
	}

	// Output:
	// true
	// 7f
	// 11
	// 3.14
	// test string
	// [1 2 3 4 5 6 7 8 9]
	// map[a:1 b:2]
	// some text
	// string in interface{}
	// {1234567890 255 12345 0 34 51 100 0}
}

func ExampleUnpack_nil() {
	var err error
	var buf bytes.Buffer

	var ba [1]byte
	var bs []byte
	var m map[string]int
	var pb *bool
	var unknown interface{}

	PackNil(&buf)
	err = Unpack(&buf, &ba)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", ba)
	}

	PackNil(&buf)
	err = Unpack(&buf, &bs)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", bs)
	}

	PackNil(&buf)
	err = Unpack(&buf, &m)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", m)
	}

	PackNil(&buf)
	err = Unpack(&buf, &pb)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		if pb == nil {
			fmt.Printf("%v\n", pb)
		} else {
			fmt.Printf("%v\n", *pb)
		}
	}

	PackNil(&buf)
	err = Unpack(&buf, &unknown)
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
	// <nil>
	// <nil>
}

func ExampleUnpack_bool() {
	var err error
	var buf bytes.Buffer
	var val bool
	var unknown interface{}

	Pack(&buf, true)
	Pack(&buf, false)
	Pack(&buf, nil)
	Pack(&buf, true)

	err = UnpackBool(&buf, &val)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", val)
	}

	err = UnpackBool(&buf, &val)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", val)
	}

	err = UnpackBool(&buf, &val)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", val)
	}

	err = Unpack(&buf, &unknown)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", unknown)
	}

	// Output:
	// true
	// false
	// false
	// true
}

func ExampleUnpack_int() {
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
	var unknown interface{}

	PackInt(&buf, 0x7f)
	UnpackInt(&buf, &i8)
	fmt.Printf("%x\n", i8)
	PackInt(&buf, 0x7fff)
	UnpackInt(&buf, &i16)
	fmt.Printf("%x\n", i16)
	PackInt(&buf, 0x7fffffff)
	UnpackInt(&buf, &i32)
	fmt.Printf("%x\n", i32)
	PackInt(&buf, 0x7fffffffffffffff)
	UnpackInt(&buf, &i64)
	fmt.Printf("%x\n", i64)

	PackUint(&buf, 0xff)
	UnpackUint(&buf, &ui8)
	fmt.Printf("%x\n", ui8)
	PackUint(&buf, 0xffff)
	UnpackUint(&buf, &ui16)
	fmt.Printf("%x\n", ui16)
	PackUint(&buf, 0xffffffff)
	UnpackUint(&buf, &ui32)
	fmt.Printf("%x\n", ui32)
	PackUint(&buf, 0xffffffffffffffff)
	UnpackUint(&buf, &ui64)
	fmt.Printf("%x\n", ui64)

	PackNil(&buf)
	err = UnpackInt(&buf, &i)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%x\n", i)
	}

	PackNil(&buf)
	err = UnpackUint(&buf, &ui)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%x\n", ui)
	}

	Pack(&buf, 1234)
	err = Unpack(&buf, &unknown)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", unknown)
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
	// 1234
}

func ExampleUnpack_float() {
	var err error
	var buf bytes.Buffer
	var f32 float32
	var f64 float64
	var unknown interface{}

	PackFloat32(&buf, 3.14)
	PackFloat64(&buf, 3.14)
	Pack(&buf, nil)
	Pack(&buf, 3.14)

	err = UnpackFloat(&buf, &f32)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", f32)
	}

	err = UnpackFloat(&buf, &f64)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", f64)
	}

	err = UnpackFloat(&buf, &f64)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", f64)
	}

	err = Unpack(&buf, &unknown)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", unknown)
	}

	// Output:
	// 3.14
	// 3.14
	// 0
	// 3.14
}

func ExampleUnpack_string() {
	var err error
	var buf bytes.Buffer
	var str string
	var unknown interface{}

	PackString(&buf, "test string")
	err = UnpackString(&buf, &str)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", str)
	}

	PackString(&buf, "012345678901234567890123456789012")
	err = UnpackString(&buf, &str)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", str)
	}

	PackArray(&buf, []byte("0123456789"))
	err = UnpackString(&buf, &str)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", str)
	}

	PackNil(&buf)
	err = UnpackString(&buf, &str)
	if err != nil {
		fmt.Println(err)
	} else {
		if len(str) == 0 {
			fmt.Printf("empty\n")
		} else {
			fmt.Printf("%v\n", str)
		}
	}

	Pack(&buf, "string in interface{}")
	err = Unpack(&buf, &unknown)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", unknown)
	}

	// Output:
	// test string
	// 012345678901234567890123456789012
	// 0123456789
	// empty
	// string in interface{}
}

func ExampleUnpack_array() {
	var err error
	var buf bytes.Buffer
	var is []int
	var i16s []int16
	var i32a [10]int32
	var strs []string
	var bin []byte
	var unknown interface{}

	Pack(&buf, []int8{1, 2, 3, 4, 5, 6, 7, 8, 9})
	err = UnpackSlice(&buf, &is)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", is)
	}

	Pack(&buf, []int64{1, 2, 3, 4, 5, 6, 7, 8, 9})
	err = UnpackSlice(&buf, &is)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", is)
	}

	Pack(&buf, []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9})
	err = UnpackSlice(&buf, &i16s)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", i16s)
	}

	Pack(&buf, []int16{1, 2, 3, 4, 5})
	err = UnpackArray(&buf, &i32a)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", i32a)
	}

	Pack(&buf, nil)
	err = UnpackSlice(&buf, &is)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", is)
	}

	Pack(&buf, nil)
	err = UnpackArray(&buf, &i32a)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", i32a)
	}

	Pack(&buf, []string{"aaa", "bbb", "ccc"})
	err = UnpackSlice(&buf, &strs)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", strs)
	}

	Pack(&buf, []byte{1, 2, 3, 4, 5, 6})
	err = UnpackSlice(&buf, &bin)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", bin)
	}

	var arr2d [][]int
	src := [][]int{{1, 2}, {3, 4}, {5, 6}}
	Pack(&buf, src)
	err = UnpackSlice(&buf, &arr2d)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", arr2d)
	}

	Pack(&buf, src)
	err = Unpack(&buf, &unknown)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", unknown)
	}

	// Output:
	// [1 2 3 4 5 6 7 8 9]
	// [1 2 3 4 5 6 7 8 9]
	// [1 2 3 4 5 6 7 8 9]
	// [1 2 3 4 5 0 0 0 0 0]
	// []
	// [0 0 0 0 0 0 0 0 0 0]
	// [aaa bbb ccc]
	// [1 2 3 4 5 6]
	// [[1 2] [3 4] [5 6]]
	// [[1 2] [3 4] [5 6]]
}

func ExampleUnpack_map() {
	var err error
	var buf bytes.Buffer
	var msi map[string]int
	var mapmap map[string]map[string]byte
	var unknown interface{}

	Pack(&buf, map[string]int16{"a": 1, "b": 2})
	err = UnpackMap(&buf, &msi)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", msi)
	}

	src := map[string]map[string]int{"first": map[string]int{"sub1": 1}, "second": map[string]int{"sub2": 2}}
	Pack(&buf, src)
	err = UnpackMap(&buf, &mapmap)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", mapmap)
	}

	Pack(&buf, src)
	err = Unpack(&buf, &unknown)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", unknown)
	}

	// Output:
	// map[a:1 b:2]
	// map[first:map[sub1:1] second:map[sub2:2]]
	// map[first:map[sub1:1] second:map[sub2:2]]
}

func ExampleUnpack_pointer() {
	var err error
	var buf bytes.Buffer
	var ptr *string

	Pack(&buf, nil)
	Pack(&buf, "some text")

	err = UnpackPtr(&buf, &ptr)
	if err != nil {
		fmt.Println(err)
	} else {
		if ptr == nil {
			fmt.Printf("%v\n", ptr)
		} else {
			fmt.Printf("%v\n", *ptr)
		}
	}

	err = UnpackPtr(&buf, &ptr)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", *ptr)
	}

	// Output:
	// <nil>
	// some text
}

func ExampleUnpack_struct() {
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

	var err error
	var buf bytes.Buffer
	var st myStruct
	var stp *myStruct

	src := myStruct{"1234567890", 255, "12345", 17, 34, 51, 100, 0}

	Pack(&buf, nil)
	Pack(&buf, nil)
	Pack(&buf, src)
	Pack(&buf, src)

	err = UnpackStruct(&buf, &st)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", st)
	}

	err = UnpackPtr(&buf, &stp)
	if err != nil {
		fmt.Println(err)
	} else {
		if stp == nil {
			fmt.Printf("%v\n", stp)
		} else {
			fmt.Printf("%v\n", *stp)
		}
	}

	err = UnpackStruct(&buf, &st)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", st)
	}

	err = UnpackPtr(&buf, &stp)
	if err != nil {
		fmt.Println(err)
	} else {
		if stp == nil {
			fmt.Printf("%v\n", stp)
		} else {
			fmt.Printf("%v\n", *stp)
		}
	}

	// Output:
	// { 0  0 0 0 0 0}
	// <nil>
	// {1234567890 255 12345 0 34 51 100 0}
	// {1234567890 255 12345 0 34 51 100 0}
}
