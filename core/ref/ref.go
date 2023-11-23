package ref

import (
	"io"
	"strconv"
	"unsafe"

	"github.com/hsfzxjy/imbed/formatter"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

const HUMANIZED_WIDTH = 12

type Ref interface {
	comparable
	formatter.Humanizer
	formatter.Stringer

	Sizeof() int
	IsZero() bool
	Raw() []byte
	RawString() string
}

type fromRaw[T any] interface {
	Ref
	fromRaw([]byte) (T, error)
}

type errSizeMismatch struct{ expected, actual int }

func (e errSizeMismatch) Error() string {
	return "ref.FromRaw: input size mismatched: expected " + strconv.Itoa(e.expected) + ", actual " + strconv.Itoa(e.actual)
}

func FromRaw[T fromRaw[T]](p []byte) (T, error, []byte) {
	var v T
	if len(p) < v.Sizeof() {
		return v, errSizeMismatch{expected: v.Sizeof(), actual: len(p)}, p
	}
	v, err := v.fromRaw(p[:v.Sizeof()])
	return v, err, p[v.Sizeof():]
}

func FromRawString[T fromRaw[T]](p string) (T, error, string) {
	v, err, rest := FromRaw[T](unsafe.Slice(unsafe.StringData(p), len(p)))
	return v, err, unsafe.String(unsafe.SliceData(rest), len(rest))
}

func FromRawExact[T fromRaw[T]](p []byte) (T, error) {
	v, err, rest := FromRaw[T](p)
	if err == nil && len(rest) != 0 {
		err = errSizeMismatch{expected: v.Sizeof(), actual: len(p)}
	}
	return v, err
}

func FromRawStringExact[T fromRaw[T]](p string) (T, error) {
	v, err, rest := FromRawString[T](p)
	if err == nil && rest != "" {
		err = errSizeMismatch{expected: v.Sizeof(), actual: len(p)}
	}
	return v, err
}

func FromReader[T fromRaw[T]](r io.Reader) (T, error) {
	var v T
	var buf = make([]byte, v.Sizeof())
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return v, err
	}
	return v.fromRaw(buf)
}

func FromFastbuf[T fromRaw[T]](r *fastbuf.R) (T, error) {
	var v T
	buf, err := r.ReadRaw(v.Sizeof())
	if err != nil {
		return v, err
	}
	return v.fromRaw(buf)
}

func AppendRaw[T Ref](b []byte, v T) []byte {
	return append(b, v.Raw()...)
}
