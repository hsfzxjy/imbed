package schema

import (
	"io"
	"unsafe"

	"github.com/tinylib/msgp/msgp"
)

type _Ptr[T any] struct {
	def        optional[*T]
	elemSchema schema[T]
}

func (s *_Ptr[T]) decodeMsg(r *msgp.Reader, target unsafe.Pointer) *schemaError {
	{
		head, err := r.ReadBool()
		if err != nil {
			return newError(err)
		}
		if !head {
			*(**T)(target) = nil
			return nil
		}
	}
	holder := new(T)
	err := s.elemSchema.decodeMsg(r, unsafe.Pointer(holder))
	if err != nil {
		return err
	}
	*(**T)(target) = holder
	return nil
}

func (s *_Ptr[T]) scanFrom(r Scanner, target unsafe.Pointer) *schemaError {
	holder := new(T)
	err := s.elemSchema.scanFrom(r, unsafe.Pointer(holder))
	if err != nil {
		err = s.setDefault(unsafe.Pointer(&holder))
		if err != nil {
			return err
		}
	}
	*(**T)(target) = holder
	return nil
}

func (s *_Ptr[T]) encodeMsg(w *msgp.Writer, source unsafe.Pointer) *schemaError {
	data := (**T)(source)
	if *data == nil {
		return newError(w.WriteBool(false))
	} else {
		err := w.WriteBool(true)
		if err != nil {
			return newError(err)
		}
	}
	return s.elemSchema.encodeMsg(w, unsafe.Pointer(*data))
}

func (s *_Ptr[T]) visit(v Visitor, source unsafe.Pointer) *schemaError {
	data := *(**T)(source)
	var isNil, isDefault bool
	var def *T
	if s.hasDefault() {
		s.setDefault(unsafe.Pointer(&def))
		isDefault = s.equal(data, def)
	}
	isNil = data == nil
	ev, err := v.VisitPtr(isNil, isDefault)
	if err != nil {
		return newError(err)
	}
	if ev == nil || isNil {
		return nil
	}
	return s.elemSchema.visit(ev, unsafe.Pointer(data))
}

func (s *_Ptr[T]) setDefault(target unsafe.Pointer) *schemaError {
	if !s.hasDefault() {
		return newError(ErrRequired)
	}
	if s.def.IsValid {
		*(**T)(target) = s.def.Value
		return nil
	}
	holder := new(T)
	err := s.elemSchema.setDefault(unsafe.Pointer(holder))
	if err != nil {
		return err
	}
	*(**T)(target) = holder
	return nil
}

func (s *_Ptr[T]) hasDefault() bool {
	return s.def.IsValid || s.elemSchema.hasDefault()
}

func (s *_Ptr[T]) equal(a, b *T) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return s.elemSchema.equal(*a, *b)
}

func (s *_Ptr[T]) writeTypeInfo(w io.Writer) error {
	_, err := w.Write([]byte{'P'})
	return err
}

func (s *_Ptr[T]) _schema_stub(*T) {}
