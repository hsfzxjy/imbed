package schema

import (
	"io"
	"unsafe"

	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type Schema[T any] interface {
	ScanFrom(r Scanner) (T, error)
	DecodeMsg(r *fastbuf.R) (T, error)
	EncodeMsg(w *fastbuf.W, source T)
	Visit(v Visitor, source T) error
}

type schema[T any] interface {
	genericSchema
	equal(a, b T) bool
	_schema_stub(T)
}

type genericSchema interface {
	scanFrom(r Scanner, target unsafe.Pointer) *schemaError
	decodeMsg(r *fastbuf.R, target unsafe.Pointer) *schemaError
	encodeMsg(w *fastbuf.W, source unsafe.Pointer)
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
