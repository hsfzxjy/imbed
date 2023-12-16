package fastbuf

import (
	"encoding/binary"
	"math/big"

	"github.com/tinylib/msgp/msgp"
)

type W struct {
	buf []byte
}

func (b *W) WriteUint32Array(arr []uint32) *W {
	b.buf = msgp.AppendArrayHeader(b.buf, uint32(len(arr)))
	for _, x := range arr {
		b.buf = binary.BigEndian.AppendUint32(b.buf, x)
	}
	return b
}

func (b *W) WriteIntArray(arr []int) *W {
	b.buf = msgp.AppendArrayHeader(b.buf, uint32(len(arr)))
	for _, x := range arr {
		b.buf = msgp.AppendInt(b.buf, x)
	}
	return b
}

func (b *W) WriteArrayHeader(sz uint32) *W {
	b.buf = msgp.AppendArrayHeader(b.buf, sz)
	return b
}

func (b *W) WriteFloat64(x float64) *W {
	b.buf = msgp.AppendFloat64(b.buf, x)
	return b
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

func (b *W) WriteInt(x int) *W {
	b.buf = msgp.AppendInt(b.buf, x)
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

func (b *W) AppendUsize(x uint64) *W {
	b.buf = binary.BigEndian.AppendUint64(b.buf, x)
	return b
}
