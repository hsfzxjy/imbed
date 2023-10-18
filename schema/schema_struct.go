package schema

import (
	"cmp"
	"errors"
	"io"
	"slices"
	"unsafe"

	"github.com/tinylib/msgp/msgp"
)

type _Struct[T any] struct {
	name   string
	fields []*_StructField
	m      map[string]*_StructField
}

func new_Struct[T any](name string, fields []*_StructField) *_Struct[T] {
	slices.SortFunc(fields, func(a, b *_StructField) int {
		return cmp.Compare(a.name, b.name)
	})
	m := make(map[string]*_StructField, len(fields))
	for _, f := range fields {
		m[f.name] = f
	}
	return &_Struct[T]{name, fields, m}
}

func (s *_Struct[T]) scanFrom(r Scanner, target unsafe.Pointer) *schemaError {
	var seen = map[string]struct{}{}
	err := r.IterField(func(name string, r Scanner) error {
		f, ok := s.m[name]
		if !ok {
			return unexpectedField(name)
		}
		err := f.scanFrom(r, target)
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

func (s *_Struct[T]) visit(v Visitor, source unsafe.Pointer) *schemaError {
	sv, ev, err := v.VisitStruct(len(s.fields))
	if err != nil {
		return newError(err).AppendPath(s.name)
	}
	if err := sv.VisitStructBegin(len(s.fields)); err != nil {
		return newError(err).AppendPath(s.name)
	}
	for _, f := range s.fields {
		if err := sv.VisitStructFieldBegin(f.name); err != nil {
			return newError(err).AppendPath(f.name).AppendPath(s.name)
		}
		if err := f.elemSchema.visit(ev, source); err != nil {
			return err.AppendPath(s.name)
		}
		if err := sv.VisitStructFieldEnd(f.name); err != nil {
			return newError(err).AppendPath(f.name).AppendPath(s.name)
		}
	}
	if err := sv.VisitStructEnd(len(s.fields)); err != nil {
		return newError(err).AppendPath(s.name)
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
