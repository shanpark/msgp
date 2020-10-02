package msgp

import "io"

type PeekableReader struct {
	rd   io.Reader // reader provided by the client
	full bool
	byt  byte
}

func NewPeekableReader(rd io.Reader) *PeekableReader {
	// Is it already a Reader?
	b, ok := rd.(*PeekableReader)
	if ok {
		return b
	}
	r := new(PeekableReader)
	r.rd = rd
	return r
}

func (b *PeekableReader) Read(p []byte) (n int, err error) {
	len := len(p)
	if b.full {
		if len == 1 {
			p[0] = b.byt
			b.full = false
			return 1, nil
		} else if len > 1 {
			p[0] = b.byt
			b.full = false
			read, err := b.rd.Read(p[1:])
			return read + 1, err
		}
	}
	return b.rd.Read(p)
}

func (b *PeekableReader) Peek() (byte, error) {
	if !b.full {
		buf := []byte{0}
		if _, err := b.rd.Read(buf); err != nil {
			return 0, err
		}
		b.byt = buf[0]
		b.full = true
	}
	return b.byt, nil
}
