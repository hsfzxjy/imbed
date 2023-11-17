package fastbuf

import (
	"math/big"
	"unsafe"

	"github.com/tinylib/msgp/msgp"
)

type R struct {
	Buf []byte
}

func read[T any](f func([]byte) (T, []byte, error), r *R) (T, error) {
	ret, buf, err := f(r.Buf)
	if err != nil {
		return ret, err
	}
	r.Buf = buf
	return ret, nil
}

func (r *R) EOF() bool {
	return len(r.Buf) == 0
}

func (r *R) ReadByte() (byte, error) {
	return read(msgp.ReadByteBytes, r)
}

func (r *R) ReadRaw(size int) ([]byte, error) {
	if len(r.Buf) < size {
		return nil, msgp.ErrShortBytes
	}
	ret := r.Buf[:size]
	r.Buf = r.Buf[size:]
	return ret, nil
}

func (r *R) ReadFull(p []byte) error {
	if len(r.Buf) < len(p) {
		return msgp.ErrShortBytes
	}
	copy(p, r.Buf[:len(p)])
	r.Buf = r.Buf[len(p):]
	return nil
}

func (r *R) ReadBytes() ([]byte, error) {
	return read(msgp.ReadBytesZC, r)
}

func (r *R) ReadInt64() (int64, error) {
	return read(msgp.ReadInt64Bytes, r)
}

func (r *R) ReadBool() (bool, error) {
	return read(msgp.ReadBoolBytes, r)
}

func (r *R) ReadRat() (*big.Rat, error) {
	p, err := read(msgp.ReadBytesZC, r)
	if err != nil {
		return nil, err
	}
	var x big.Rat
	if err := x.GobDecode(p); err != nil {
		return nil, err
	}
	return &x, nil
}

func (r *R) ReadString() (string, error) {
	v, err := read(msgp.ReadStringZC, r)
	if err != nil {
		return "", err
	}
	return unsafe.String(unsafe.SliceData(v), len(v)), nil
}
