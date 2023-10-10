package schema

import (
	"errors"
	"io"
	"math/big"
	"unsafe"

	"github.com/tinylib/msgp/msgp"
)

type _AtomVTable[T any] struct {
	typeName        string
	decodeMsgFunc   func(r *msgp.Reader) (T, error)
	decodeValueFunc func(r Reader) (T, error)
	encodeMsgFunc   func(w *msgp.Writer, value T) error
	visitFunc       func(v Visitor, value T) error
}

type _Atom[T any] struct {
	def optional[T]
	*_AtomVTable[T]
}

func (s *_Atom[T]) decodeMsg(r *msgp.Reader, target unsafe.Pointer) *schemaError {
	val, err := s.decodeMsgFunc(r)
	if err != nil {
		return newError(err)
	}
	*(*T)(target) = val
	return nil
}

func (s *_Atom[T]) decodeValue(r Reader, target unsafe.Pointer) *schemaError {
	val, err := s.decodeValueFunc(r)
	if err != nil {
		if !errors.Is(err, ErrRequired) {
			return newError(err)
		}
		if !s.def.IsValid {
			return newError(ErrRequired)
		}
		val = s.def.Value
	}
	*(*T)(target) = val
	return nil
}

func (s *_Atom[T]) encodeMsg(w *msgp.Writer, source unsafe.Pointer) *schemaError {
	return newError(s.encodeMsgFunc(w, (*(*T)(source))))
}

func (s *_Atom[T]) visit(v Visitor, source unsafe.Pointer) *schemaError {
	return newError(s.visitFunc(v, (*(*T)(source))))
}

func (s *_Atom[T]) setDefault(target unsafe.Pointer) *schemaError {
	if !s.def.IsValid {
		return newError(ErrRequired)
	}
	*(*T)(target) = s.def.Value
	return nil
}

func (s *_Atom[T]) writeTypeInfo(w io.Writer) error {
	_, err := w.Write([]byte(s.typeName))
	return err
}

func (s *_Atom[T]) _schema_stub(T) {}

type _Int = _Atom[int64]

var _VTableInt = &_AtomVTable[int64]{
	typeName:        "int64",
	decodeMsgFunc:   (*msgp.Reader).ReadInt64,
	decodeValueFunc: (Reader).Int64,
	encodeMsgFunc:   (*msgp.Writer).WriteInt64,
	visitFunc:       Visitor.VisitInt64,
}

func new_Int(def optional[int64]) *_Int { return &_Int{def, _VTableInt} }

type _String = _Atom[string]

var _VTableString = &_AtomVTable[string]{
	typeName:        "string",
	decodeMsgFunc:   (*msgp.Reader).ReadString,
	decodeValueFunc: (Reader).String,
	encodeMsgFunc:   (*msgp.Writer).WriteString,
	visitFunc:       Visitor.VisitString,
}

func new_String(def optional[string]) *_String { return &_String{def, _VTableString} }

type _Bool = _Atom[bool]

var _VTableBool = &_AtomVTable[bool]{
	typeName:        "bool",
	decodeMsgFunc:   (*msgp.Reader).ReadBool,
	decodeValueFunc: (Reader).Bool,
	encodeMsgFunc:   (*msgp.Writer).WriteBool,
	visitFunc:       Visitor.VisitBool,
}

func new_Bool(def optional[bool]) *_Bool { return &_Bool{def, _VTableBool} }

func _() {
	var _ schema[int64] = &_Int{}
	var _ schema[bool] = &_Bool{}
	var _ schema[string] = &_String{}
	var _ schema[*big.Rat] = &_Rat{}
}
