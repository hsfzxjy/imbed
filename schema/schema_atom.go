package schema

import (
	"errors"
	"io"
	"unsafe"

	"github.com/tinylib/msgp/msgp"
)

type _AtomVTable[T any] struct {
	typeName        string
	decodeMsgFunc   func(r *msgp.Reader) (T, error)
	decodeValueFunc func(r Reader) (T, error)
	encodeMsgFunc   func(w *msgp.Writer, value T) error
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

func (s *_Atom[T]) _schemaTyped_stub(T) {}

type _Int = _Atom[int64]

var _VTableInt = &_AtomVTable[int64]{
	typeName:        "int64",
	decodeMsgFunc:   (*msgp.Reader).ReadInt64,
	decodeValueFunc: (Reader).Int64,
	encodeMsgFunc:   (*msgp.Writer).WriteInt64,
}

func new_Int(def optional[int64]) *_Int { return &_Int{def, _VTableInt} }

type _String = _Atom[string]

var _VTableString = &_AtomVTable[string]{
	typeName:        "string",
	decodeMsgFunc:   (*msgp.Reader).ReadString,
	decodeValueFunc: (Reader).String,
	encodeMsgFunc:   (*msgp.Writer).WriteString,
}

func new_String(def optional[string]) *_String { return &_String{def, _VTableString} }

type _Bool = _Atom[bool]

var _VTableBool = &_AtomVTable[bool]{
	typeName:        "bool",
	decodeMsgFunc:   (*msgp.Reader).ReadBool,
	decodeValueFunc: (Reader).Bool,
	encodeMsgFunc:   (*msgp.Writer).WriteBool,
}

func new_Bool(def optional[bool]) *_Bool { return &_Bool{def, _VTableBool} }

type _Float = _Atom[float64]

var _VTableFloat = &_AtomVTable[float64]{
	typeName:        "float64",
	decodeMsgFunc:   (*msgp.Reader).ReadFloat64,
	decodeValueFunc: (Reader).Float64,
	encodeMsgFunc:   (*msgp.Writer).WriteFloat64,
}

func new_Float(def optional[float64]) *_Float { return &_Float{def, _VTableFloat} }

func _() {
	var _ schemaTyped[int64] = &_Int{}
	var _ schemaTyped[bool] = &_Bool{}
	var _ schemaTyped[string] = &_String{}
	var _ schemaTyped[float64] = &_Float{}
}
