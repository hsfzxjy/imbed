package schema

import (
	"io"
	"unsafe"

	"github.com/tinylib/msgp/msgp"
)

type Schema[T any] interface {
	ScanFrom(r Scanner) (T, error)

	DecodeMsg(r *msgp.Reader) (T, error)
	EncodeMsg(w *msgp.Writer, source T) error
	Visit(v Visitor, source T) error
	New() T
	GenericSchema
}

type GenericSchema interface {
	genericSchema
	ScanFromAny(r Scanner) (any, error)
	DecodeMsgAny(r *msgp.Reader) (any, error)
	EncodeMsgAny(w *msgp.Writer, source any) error
	VisitAny(v Visitor, source any) error
	NewAny() any
	WrapAny(data any) GenericValue
}

type schema[T any] interface {
	genericSchema
	equal(a, b T) bool
	_schema_stub(T)
}

type genericSchema interface {
	scanFrom(r Scanner, target unsafe.Pointer) *schemaError
	decodeMsg(r *msgp.Reader, target unsafe.Pointer) *schemaError
	encodeMsg(w *msgp.Writer, source unsafe.Pointer) *schemaError
	visit(v Visitor, source unsafe.Pointer) *schemaError
	setDefault(target unsafe.Pointer) *schemaError
	hasDefault() bool
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
