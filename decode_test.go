package msgp

import (
	"bytes"
	"fmt"
)

func ExampleUnpackValue() {
	var buf bytes.Buffer
	var ptr *int64

	PackInt(&buf, 0x7f)
	UnpackValue(&buf, &ptr)

	fmt.Printf("%x", *ptr)
	// Output:
	// 7f
}
