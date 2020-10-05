# github.com/shanpark/msgp
MsgPack Library for Go

# Usage
See URL below for more information<br/>
    - https://godoc.org/github.com/shanpark/msgp

Pack...
<pre><code>err = PackValue(&buf, nil)      // nil
err = PackValue(&buf, true)     // bool
err = PackValue(&buf, 100)      // int
err = PackValue(&buf, 3.14)     // float
err = PackValue(&buf, "aaa")    // string
err = PackValue(&buf, []string{"aaa", "bbb", "ccc"})    // array
err = PackValue(&buf, map[string]int{"aaa": 1})         // map

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
PackValue(&buf, st)             // struct

var a int = 100
err = PackValue(&buf, &a)       // pointer
</code></pre>
Unpack...
<pre><code>// Bool
var b bool
err = UnpackValue(&buf, &b)

// Integer
var i int
err = UnpackValue(&buf, &i)

// Float
var f float32
err = UnpackValue(&buf, &f)

// String
var str string
err = UnpackValue(&buf, &str)

// Array(Slice)
var ia []int
err = UnpackValue(&buf, &ia)

// Map
var msi map[string]int
err = UnpackValue(&buf, &msi)

// Pointer (of pointer)
var ptr *string
err = UnpackValue(&buf, &ptr)

// Unknown type
var unknown interface{}
err = UnpackValue(&buf, &unknown)
</code></pre>