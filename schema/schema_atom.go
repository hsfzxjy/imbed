package schema

import (
	"errors"
	"io"
	"math/big"
	"unsafe"

	"github.com/hsfzxjy/imbed/core/pos"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type _AtomVTable[T comparable] struct {
	typeName      string
	decodeMsgFunc func(r *fastbuf.R) (T, error)
	scanFromFunc  func(r Scanner) (T, pos.P, error)
	encodeMsgFunc func(w *fastbuf.W, value T) *fastbuf.W
	visitFunc     func(v Visitor, value T, isDefault bool) error
	cmpFunc       func(a, b T) int
}

type _Atom[T comparable] struct {
	def optional[T]
	*_AtomVTable[T]
}

func (s *_Atom[T]) decodeMsg(r *fastbuf.R, target unsafe.Pointer) *schemaError {
	val, err := s.decodeMsgFunc(r)
	if err != nil {
		return newError(err)
	}
	*(*T)(target) = val
	return nil
}

func (s *_Atom[T]) scanFrom(r Scanner, target unsafe.Pointer) (pos.P, *schemaError) {
	val, pos, err := s.scanFromFunc(r)
	if err != nil {
		if !errors.Is(err, ErrRequired) {
			return pos, newError(err)
		}
		if !s.def.IsValid {
			return pos, newError(ErrRequired)
		}
		val = s.def.Value
	}
	*(*T)(target) = val
	return pos, nil
}

func (s *_Atom[T]) encodeMsg(w *fastbuf.W, source unsafe.Pointer) {
	s.encodeMsgFunc(w, (*(*T)(source)))
}

func (s *_Atom[T]) visit(v Visitor, source unsafe.Pointer) *schemaError {
	value := (*(*T)(source))
	var isDefault bool
	if s.def.IsValid {
		if s.cmpFunc != nil {
			isDefault = s.cmpFunc(value, s.def.Value) == 0
		} else {
			isDefault = value == s.def.Value
		}
	}
	return newError(s.visitFunc(v, value, isDefault))
}

func (s *_Atom[T]) setDefault(target unsafe.Pointer) *schemaError {
	if !s.def.IsValid {
		return newError(ErrRequired)
	}
	*(*T)(target) = s.def.Value
	return nil
}

func (s *_Atom[T]) hasDefault() bool {
	return s.def.IsValid
}

func (s *_Atom[T]) equal(a, b T) bool {
	if s.cmpFunc != nil {
		return s.cmpFunc(a, b) == 0
	} else {
		return a == b
	}
}

func (s *_Atom[T]) writeTypeInfo(w io.Writer) error {
	_, err := w.Write([]byte(s.typeName))
	return err
}

func (s *_Atom[T]) _schema_stub(T) {}

type _Int = _Atom[int64]

var _VTableInt = &_AtomVTable[int64]{
	typeName:      "int64",
	decodeMsgFunc: (*fastbuf.R).ReadInt64,
	scanFromFunc:  (Scanner).Int64,
	encodeMsgFunc: (*fastbuf.W).WriteInt64,
	visitFunc:     Visitor.VisitInt64,
}

func new_Int(def optional[int64]) *_Int { return &_Int{def, _VTableInt} }

type _String = _Atom[string]

var _VTableString = &_AtomVTable[string]{
	typeName:      "string",
	decodeMsgFunc: (*fastbuf.R).ReadString,
	scanFromFunc:  (Scanner).String,
	encodeMsgFunc: (*fastbuf.W).WriteString,
	visitFunc:     Visitor.VisitString,
}

func new_String(def optional[string]) *_String { return &_String{def, _VTableString} }

type _Bool = _Atom[bool]

var _VTableBool = &_AtomVTable[bool]{
	typeName:      "bool",
	decodeMsgFunc: (*fastbuf.R).ReadBool,
	scanFromFunc:  (Scanner).Bool,
	encodeMsgFunc: (*fastbuf.W).WriteBool,
	visitFunc:     Visitor.VisitBool,
}

func new_Bool(def optional[bool]) *_Bool { return &_Bool{def, _VTableBool} }

func _() {
	var _ schema[int64] = &_Int{}
	var _ schema[bool] = &_Bool{}
	var _ schema[string] = &_String{}
	var _ schema[*big.Rat] = &_Rat{}
}
