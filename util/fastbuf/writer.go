package fastbuf

import (
	"math/big"

	"github.com/tinylib/msgp/msgp"
)

type W struct {
	buf []byte
}

func (b *W) WriteUint8(x byte) *W {
	b.buf = append(b.buf, x)
	return b
}

func (b *W) WriteString(s string) *W {
	b.buf = msgp.AppendString(b.buf, s)
	return b
}

func (b *W) WriteBytes(p []byte) *W {
	b.buf = msgp.AppendBytes(b.buf, p)
	return b
}

func (b *W) WriteInt64(x int64) *W {
	b.buf = msgp.AppendInt64(b.buf, x)
	return b
}

func (b *W) WriteBool(x bool) *W {
	b.buf = msgp.AppendBool(b.buf, x)
	return b
}

func (b *W) WriteRat(x *big.Rat) *W {
	encoded, err := x.GobEncode()
	if err != nil {
		panic(err)
	}
	b.buf = msgp.AppendBytes(b.buf, encoded)
	return b
}

func (b *W) AppendRaw(p []byte) *W {
	b.buf = append(b.buf, p...)
	return b
}

func (b *W) Result() []byte {
	return b.buf
}
