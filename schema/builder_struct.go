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

type StructBuilder[T any] struct {
	name          string
	basePtr       unsafe.Pointer
	fieldBuilders []*fieldBuilder
}

func Struct[S any](ptr *S, fieldBuilders ...*fieldBuilder) *StructBuilder[S] {
	basePtr := unsafe.Pointer(ptr)
	var x S
	structSize := unsafe.Sizeof(x)
	for _, field := range fieldBuilders {
		if !(uintptr(basePtr) <= uintptr(field.ptr) &&
			uintptr(field.ptr) < uintptr(basePtr)+structSize) {
			panic("ptr of " + field.name + " underflows/overflows")
		}
	}
	return &StructBuilder[S]{"<anonymous>", basePtr, fieldBuilders}
}

func (s *StructBuilder[T]) DebugName(name string) *StructBuilder[T] {
	s.name = name
	return s
}

func (s *StructBuilder[S]) buildStruct() *_Struct[S] {
	fields := make([]*_StructField, len(s.fieldBuilders))
	for i, f := range s.fieldBuilders {
		field := &_StructField{
			name:       f.name,
			proto:      goStructFieldProto{uintptr(f.ptr) - uintptr(s.basePtr)},
			elemSchema: f.subBuilder.buildGenericSchema(),
		}
		fields[i] = field
	}
	return new_Struct[S](s.name, fields, goStructProto[S]{})
}

func (s *StructBuilder[T]) buildGenericSchema() genericSchema { return s.buildStruct() }
func (s *StructBuilder[T]) buildSchema() schema[T]            { return s.buildStruct() }

func (s *StructBuilder[T]) Build() Schema[*T] {
	return s.buildStruct().Build()
}

func StructFunc[S any](f func(*S) *StructBuilder[S]) *_Struct[S] {
	var prototype S
	return f(&prototype).buildStruct()
}
