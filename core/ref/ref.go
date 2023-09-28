package ref

import (
	"io"
	"strings"
	"unsafe"
)

type Ref interface {
	comparable
	getInternalString() string
	Len() int
}

func (fid FID) getInternalString() string { return fid.raw }

func (h Murmur3Hash) getInternalString() string { return h.raw }

func (h Sha256Hash) getInternalString() string {
	return h.raw
}

func Compare[T Ref](a, b T) int {
	return strings.Compare(a.getInternalString(), b.getInternalString())
}

func IsEqual[T Ref](a, b T) bool {
	return a.getInternalString() == b.getInternalString()
}

type fromRaw[T any] interface {
	Ref
	fromBytes([]byte) (T, []byte)
}

func FromRaw[T fromRaw[T]](p []byte) T {
	var v T
	v, p = v.fromBytes(p)
	return v
}

func FromRawString[T fromRaw[T]](p string) T {
	var v T
	v, _ = v.fromBytes(unsafe.Slice(unsafe.StringData(p), len(p)))
	return v
}

func AsRaw[T Ref](v T) []byte {
	s := v.getInternalString()
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func AsRawString[T Ref](v T) string {
	return v.getInternalString()
}

func FromReader[T interface {
	fromRaw[T]
	expectedSize() int
}](r io.Reader) (T, error) {
	var v T
	var buf = make([]byte, v.expectedSize())
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return v, err
	}
	v, _ = v.fromBytes(buf)
	return v, nil
}
