# github.com/shanpark/msgp
MsgPack Library for Go

# Usage
See URL below for more information<br/>
    - https://godoc.org/github.com/shanpark/msgp

Pack...
<pre><code>err = Pack(&buf, nil)      // nil
err = Pack(&buf, true)     // bool
err = Pack(&buf, 100)      // int
err = Pack(&buf, 3.14)     // float
err = Pack(&buf, "aaa")    // string
err = Pack(&buf, []string{"aaa", "bbb", "ccc"})    // array
err = Pack(&buf, map[string]int{"aaa": 1})         // map

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
Pack(&buf, st)             // struct

var a int = 100
err = Pack(&buf, &a)       // pointer
</code></pre>
Unpack...
<pre><code>// Bool
var b bool
err = Unpack(&buf, &b)

// Integer
var i int
err = Unpack(&buf, &i)

// Float
var f float32
err = Unpack(&buf, &f)

// String
var str string
err = Unpack(&buf, &str)

// Array(Slice)
var ia []int
err = Unpack(&buf, &ia)

// Map
var msi map[string]int
err = Unpack(&buf, &msi)

// Struct
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
err = Unpack(&buf, &st)

// Unknown type
var unknown interface{}
err = Unpack(&buf, &unknown)

// Pointer (of pointer)
var ptr *string
err = Unpack(&buf, &ptr)
</code></pre>