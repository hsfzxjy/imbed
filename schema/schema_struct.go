package schema

import (
	"cmp"
	"errors"
	"io"
	"slices"
	"sync"
	"unsafe"

	"github.com/hsfzxjy/imbed/core/pos"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type _Struct[T any] struct {
	name      string
	fields    []*_StructField
	m         map[string]*_StructField
	mainField *_StructField
	hasDef    bool

	buildOnce sync.Once
	toplevel  *_TopLevel[T]
}

func new_Struct[T any](name string, fields []*_StructField) *_Struct[T] {
	slices.SortFunc(fields, func(a, b *_StructField) int {
		return cmp.Compare(a.name, b.name)
	})
	m := make(map[string]*_StructField, len(fields))
	var mainField *_StructField
	var noDefCount int
	for _, f := range fields {
		m[f.name] = f
		if !f.hasDefault() {
			mainField = f
			noDefCount++
		}
	}
	if noDefCount != 1 {
		mainField = nil
	}
	if len(fields) == 1 {
		mainField = fields[0]
	}
	return &_Struct[T]{
		name:      name,
		fields:    fields,
		m:         m,
		mainField: mainField,
		hasDef:    noDefCount == 0,
	}
}

func (s *_Struct[T]) scanFrom(r Scanner, target unsafe.Pointer) (pos.P, *schemaError) {
	var seen = map[string]struct{}{}
	var normalScan = true
	var err *schemaError
	var structPos pos.P
	if s.mainField != nil {
		mainScanner := r.UnnamedField()
		if mainScanner != nil {
			structPos, err = s.mainField.scanFrom(mainScanner, target)
			if err == nil {
				normalScan = false
				seen[s.mainField.Name()] = struct{}{}
			}
		}
	}
	if normalScan {
		var fieldPos pos.P
		err := r.IterField(func(name string, r Scanner, namePos pos.P) error {
			f, ok := s.m[name]
			if !ok {
				return unexpectedField(name)
			}
			fieldPos, err = f.scanFrom(r, target)
			structPos = structPos.Add(fieldPos)
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
			return fieldPos, newError(err).AppendPath(s.name)
		}
	}
	for _, f := range s.fields {
		if _, ok := seen[f.Name()]; !ok {
			err := f.setDefault(target)
			if err != nil {
				return structPos, err.AppendPath(s.name)
			}
		}
	}
	{
		validator, ok := any((*T)(target)).(Validator)
		if ok {
			err := validator.Validate()
			if err != nil {
				return structPos, newError(validation(err)).AppendPath(s.name)
			}
		}
	}
	return structPos, nil
}

func (s *_Struct[T]) decodeMsg(r *fastbuf.R, target unsafe.Pointer) *schemaError {
	for _, f := range s.fields {
		err := f.decodeMsg(r, target)
		if err != nil {
			return err.AppendPath(s.name)
		}
	}
	return nil
}

func (s *_Struct[T]) encodeMsg(w *fastbuf.W, source unsafe.Pointer) {
	for _, f := range s.fields {
		f.encodeMsg(w, source)
	}
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
		if err := f.visit(ev, source); err != nil {
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

func (s *_Struct[T]) hasDefault() bool {
	return s.hasDef
}

func (s *_Struct[T]) equal(a, b T) bool {
	return false
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

func (s *_Struct[T]) buildGenericSchema() genericSchema {
	return s
}

func (s *_Struct[T]) buildSchema() schema[T] {
	return s
}

func (s *_Struct[T]) Build() *_TopLevel[T] {
	s.buildOnce.Do(func() {
		s.toplevel = new_Toplevel(s)
	})
	return s.toplevel
}

func _() {
	type X struct{}
	var _ schema[X] = &_Struct[X]{}
	var _ builder[X] = &_Struct[X]{}
}
