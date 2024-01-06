package schema

import (
	"io"
	"unsafe"

	"github.com/hsfzxjy/imbed/core/pos"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type Schema[T any] interface {
	ScanFrom(r Scanner) (T, pos.P, error)
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
	scanFrom(r Scanner, target unsafe.Pointer) (pos.P, *schemaError)
	decodeMsg(r *fastbuf.R, target unsafe.Pointer) *schemaError
	encodeMsg(w *fastbuf.W, source unsafe.Pointer)
	visit(v Visitor, source unsafe.Pointer) *schemaError
	setDefault(target unsafe.Pointer) *schemaError
	hasDefault() bool
	writeTypeInfo(w io.Writer) error
}

// type fieldSchema interface {
// 	genericSchema
// 	_fieldSchema_stub()
// 	Name() string
// }

type Validator interface {
	Validate() error
}

type optional[T any] struct {
	IsValid bool
	Value   T
}

type defaulter interface {
	hasDefault() bool
	setDefault(target unsafe.Pointer) *schemaError
}

type objectProto interface {
	isEqual(a, b unsafe.Pointer) bool
}

type listProto interface {
	ListInit(target unsafe.Pointer, size int)
	ListLen(target unsafe.Pointer) int
	ListElem(target unsafe.Pointer, i int) unsafe.Pointer
}

type mapProto interface {
	MapInit(target unsafe.Pointer, size int)
	MapNewValue() unsafe.Pointer
	MapSetValue(m unsafe.Pointer, key string, value unsafe.Pointer)
	MapLen(target unsafe.Pointer) int
	MapIter(target unsafe.Pointer, iterFn func(key string, value unsafe.Pointer) bool)
	MapIterOrdered(target unsafe.Pointer, iterFn func(key string, value unsafe.Pointer) bool)
}

type structFieldProto interface {
	FieldPtr(structPtr unsafe.Pointer) unsafe.Pointer
}

type structProto[T any] interface {
	StructValidate(target unsafe.Pointer) error
	StructNewTopLevel(s *_Struct[T]) Schema[*T]
}
