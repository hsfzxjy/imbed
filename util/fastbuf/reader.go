package fastbuf

import (
	"encoding/binary"
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

func (r *R) ReadUint8() (byte, error) {
	if len(r.Buf) == 0 {
		return 0, msgp.ErrShortBytes
	}
	ret := r.Buf[0]
	r.Buf = r.Buf[1:]
	return ret, nil
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

func (r *R) ReadInt() (int, error) {
	return read(msgp.ReadIntBytes, r)
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

func (r *R) ReadUint32Array() ([]uint32, error) {
	sz, err := read(msgp.ReadArrayHeaderBytes, r)
	if err != nil {
		return nil, err
	}
	if len(r.Buf) < int(sz)*4 {
		return nil, msgp.ErrShortBytes
	}
	ret := make([]uint32, sz)
	for i := range ret {
		ret[i] = binary.BigEndian.Uint32(r.Buf)
		r.Buf = r.Buf[4:]
	}
	return ret, nil
}

func (r *R) ReadIntArray() ([]int, error) {
	sz, err := read(msgp.ReadArrayHeaderBytes, r)
	if err != nil {
		return nil, err
	}
	ret := make([]int, sz)
	for i := range ret {
		ret[i], err = read(msgp.ReadIntBytes, r)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (r *R) ReadArrayHeader() (uint32, error) {
	return read(msgp.ReadArrayHeaderBytes, r)
}

func (r *R) ReadFloat64() (float64, error) {
	return read(msgp.ReadFloat64Bytes, r)
}

func (r *R) ReadUsize() (uint64, error) {
	if len(r.Buf) < 8 {
		return 0, msgp.ErrShortBytes
	}
	ret := binary.BigEndian.Uint64(r.Buf)
	r.Buf = r.Buf[8:]
	return ret, nil
}

func (r *R) SplitAt(n uint64) (r1, r2 R, err error) {
	if uint64(len(r.Buf)) < n {
		err = msgp.ErrShortBytes
		return
	}
	r1.Buf = r.Buf[:n:n]
	r2.Buf = r.Buf[n:]
	return
}
