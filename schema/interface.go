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
	Visit(v Visitor, source *T) error
	schema[T]
}

type schema[T any] interface {
	genericSchema
	_schema_stub(T)
}

type genericSchema interface {
	decodeValue(r Reader, target unsafe.Pointer) *schemaError
	decodeMsg(r *msgp.Reader, target unsafe.Pointer) *schemaError
	encodeMsg(w *msgp.Writer, source unsafe.Pointer) *schemaError
	visit(v Visitor, source unsafe.Pointer) *schemaError
	setDefault(target unsafe.Pointer) *schemaError
	writeTypeInfo(w io.Writer) error
}

type fieldSchema interface {
	genericSchema
	_fieldSchema_stub()
	Name() string
}

type Validator interface {
	Validate() error
}

type optional[T any] struct {
	IsValid bool
	Value   T
}
