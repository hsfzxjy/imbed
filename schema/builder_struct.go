package schema

import (
	"cmp"
	"slices"
	"unsafe"
)

type fieldBuilder struct {
	name       string
	ptr        unsafe.Pointer
	subBuilder genericBuilder
}

func F[T any](name string, ptr *T, subBuilder builder[T]) *fieldBuilder {
	return &fieldBuilder{
		name:       name,
		ptr:        unsafe.Pointer(ptr),
		subBuilder: subBuilder,
	}
}

type structBuilder[T any] struct {
	name          string
	basePtr       unsafe.Pointer
	fieldBuilders []*fieldBuilder
}

func Struct[S any](ptr *S, fieldBuilders ...*fieldBuilder) *structBuilder[S] {
	basePtr := unsafe.Pointer(ptr)
	var x S
	structSize := unsafe.Sizeof(x)
	for _, field := range fieldBuilders {
		if !(uintptr(basePtr) <= uintptr(field.ptr) &&
			uintptr(field.ptr) < uintptr(basePtr)+structSize) {
			panic("ptr of " + field.name + " underflows/overflows")
		}
	}
	slices.SortFunc(fieldBuilders, func(a, b *fieldBuilder) int {
		return cmp.Compare(a.name, b.name)
	})
	return &structBuilder[S]{"<anonymous>", basePtr, fieldBuilders}
}

func (s *structBuilder[T]) DebugName(name string) *structBuilder[T] {
	s.name = name
	return s
}

func (s *structBuilder[S]) buildSchema() schema[S] {
	fields := make([]*_StructField, len(s.fieldBuilders))
	m := make(map[string]*_StructField, len(fields))
	for i, f := range s.fieldBuilders {
		field := &_StructField{
			name:       f.name,
			offset:     uintptr(f.ptr) - uintptr(s.basePtr),
			elemSchema: f.subBuilder.buildGenericSchema(),
		}
		fields[i] = field
		m[f.name] = field
	}
	return &_Struct[S]{s.name, fields, m}
}
func (s *structBuilder[T]) buildGenericSchema() genericSchema { return s.buildSchema() }

func (s *structBuilder[T]) Build() Schema[T] {
	return New(s)
}
