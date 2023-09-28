package schema

import (
	"io"
	"unsafe"

	"github.com/tinylib/msgp/msgp"
)

type Schema[T any] interface {
	DecodeValue(r Reader, target *T) error
	DecodeMsg(r *msgp.Reader, target *T) error
	EncodeMsg(w *msgp.Writer, source *T) error
	schemaTyped[T]
}

type schemaTyped[T any] interface {
	schema
	_schemaTyped_stub(T)
}

type schema interface {
	decodeValue(r Reader, target unsafe.Pointer) *schemaError
	decodeMsg(r *msgp.Reader, target unsafe.Pointer) *schemaError
	encodeMsg(w *msgp.Writer, source unsafe.Pointer) *schemaError
	setDefault(target unsafe.Pointer) *schemaError
	writeTypeInfo(w io.Writer) error
}

type fieldSchema interface {
	schema
	_fieldSchema_stub()
	Name() string
}

type Validator interface{ Validate() error }

type optional[T any] struct {
	IsValid bool
	Value   T
}
