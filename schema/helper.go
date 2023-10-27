package schema

import (
	"bytes"

	"github.com/tinylib/msgp/msgp"
)

func EncodeBytes[S any](schema Schema[S], value S) ([]byte, error) {
	var buf bytes.Buffer
	var w = msgp.NewWriter(&buf)
	err := schema.EncodeMsg(w, value)
	if err != nil {
		return nil, err
	}
	if err = w.Flush(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func EncodeBytesAny(schema GenericSchema, value any) ([]byte, error) {
	var buf bytes.Buffer
	var w = msgp.NewWriter(&buf)
	err := schema.EncodeMsgAny(w, value)
	if err != nil {
		return nil, err
	}
	if err = w.Flush(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type GenericValue interface {
	EncodeMsg(w *msgp.Writer) error
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

func (v *Wrapped[T]) EncodeMsg(w *msgp.Writer) error {
	return v.schema.EncodeMsg(w, v.data)
}

func (v *Wrapped[T]) Visit(visitor Visitor) error {
	return v.schema.Visit(visitor, v.data)
}

func (v *Wrapped[T]) Get() T {
	return v.data
}

func Wrap[T any](schema Schema[T], data T) *Wrapped[T] {
	return &Wrapped[T]{schema, data}
}
