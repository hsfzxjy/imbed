package schema

import (
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

func EncodeBytes[S any](schema Schema[S], value S) []byte {
	var w fastbuf.W
	schema.EncodeMsg(&w, value)
	return w.Result()
}

func EncodeBytesAny(schema GenericSchema, value any) []byte {
	var w fastbuf.W
	schema.EncodeMsgAny(&w, value)
	return w.Result()
}

type GenericValue interface {
	EncodeMsg(w *fastbuf.W)
	EncodeBytes() []byte
	Visit(visitor Visitor) error
}

type Value[T any] interface {
	GenericValue
	Get() T
}

type Wrapped[T any] struct {
	schema Schema[T]
	data   T
}

func (v *Wrapped[T]) EncodeMsg(w *fastbuf.W) {
	v.schema.EncodeMsg(w, v.data)
}

func (v *Wrapped[T]) Visit(visitor Visitor) error {
	return v.schema.Visit(visitor, v.data)
}

func (v *Wrapped[T]) EncodeBytes() []byte {
	var b fastbuf.W
	v.schema.EncodeMsg(&b, v.data)
	return b.Result()
}

func (v *Wrapped[T]) Get() T {
	return v.data
}

func Wrap[T any](schema Schema[T], data T) *Wrapped[T] {
	return &Wrapped[T]{schema, data}
}
