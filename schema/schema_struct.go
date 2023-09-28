package schema

import (
	"errors"
	"io"
	"unsafe"

	"github.com/tinylib/msgp/msgp"
)

type _Struct[T any] struct {
	name   string
	fields []*_StructField
	m      map[string]*_StructField
}

func (s *_Struct[T]) decodeValue(r Reader, target unsafe.Pointer) *schemaError {
	var seen = map[string]struct{}{}
	err := r.IterField(func(name string, r Reader) error {
		f, ok := s.m[name]
		if !ok {
			return unexpectedField(name)
		}
		err := f.decodeValue(r, target)
		if err != nil {
			if !errors.Is(err.AsError(), ErrRequired) {
				return err.AsError()
			}
		} else {
			seen[name] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return newError(err).AppendPath(s.name)
	}
	for _, f := range s.fields {
		if _, ok := seen[f.Name()]; !ok {
			err := f.setDefault(target)
			if err != nil {
				return err.AppendPath(s.name)
			}
		}
	}
	{
		validator, ok := any((*T)(target)).(Validator)
		if ok {
			err := validator.Validate()
			if err != nil {
				return newError(validation(err)).AppendPath(s.name)
			}
		}
	}
	return nil
}

func (s *_Struct[T]) decodeMsg(r *msgp.Reader, target unsafe.Pointer) *schemaError {
	for _, f := range s.fields {
		err := f.decodeMsg(r, target)
		if err != nil {
			return err.AppendPath(s.name)
		}
	}
	return nil
}

func (s *_Struct[T]) encodeMsg(w *msgp.Writer, source unsafe.Pointer) *schemaError {
	for _, f := range s.fields {
		err := f.encodeMsg(w, source)
		if err != nil {
			return err.AppendPath(s.name)
		}
	}
	return nil
}

func (s *_Struct[T]) setDefault(target unsafe.Pointer) *schemaError {
	for _, f := range s.fields {
		err := f.setDefault(target)
		if err != nil {
			return err.AppendPath(s.name)
		}
	}
	return nil
}

func (s *_Struct[T]) writeTypeInfo(w io.Writer) error {
	for _, f := range s.fields {
		err := f.writeTypeInfo(w)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *_Struct[T]) _schema_stub(T) {}

func _() {
	type X struct{}
	var _ schema[X] = &_Struct[X]{}
}
