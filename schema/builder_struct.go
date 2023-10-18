package schema

import (
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
	return &structBuilder[S]{"<anonymous>", basePtr, fieldBuilders}
}

func (s *structBuilder[T]) DebugName(name string) *structBuilder[T] {
	s.name = name
	return s
}

func (s *structBuilder[S]) buildStruct() *_Struct[S] {
	fields := make([]*_StructField, len(s.fieldBuilders))
	for i, f := range s.fieldBuilders {
		field := &_StructField{
			name:       f.name,
			offset:     uintptr(f.ptr) - uintptr(s.basePtr),
			elemSchema: f.subBuilder.buildGenericSchema(),
		}
		fields[i] = field
	}
	return new_Struct[S](s.name, fields)
}
func (s *structBuilder[T]) buildGenericSchema() genericSchema { return s.buildStruct() }
func (s *structBuilder[T]) buildSchema() schema[T]            { return s.buildStruct() }

func (s *structBuilder[T]) Build() Schema[*T] {
	return New(s)
}
