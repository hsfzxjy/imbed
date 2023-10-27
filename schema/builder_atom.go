package schema

import "math/big"

type atomBuilder[T comparable] struct {
	def    optional[T]
	vtable *_AtomVTable[T]
}

func (b atomBuilder[T]) Default(value T) atomBuilder[T] {
	b.def.Value = value
	b.def.IsValid = true
	return b
}

func (b atomBuilder[T]) buildSchema() schema[T]            { return &_Atom[T]{b.def, b.vtable} }
func (b atomBuilder[T]) buildGenericSchema() genericSchema { return b.buildSchema() }

func Int() atomBuilder[int64] {
	return atomBuilder[int64]{vtable: _VTableInt}
}

func Bool() atomBuilder[bool] {
	return atomBuilder[bool]{vtable: _VTableBool}
}

func String() atomBuilder[string] {
	return atomBuilder[string]{vtable: _VTableString}
}

func Rat() atomBuilder[*big.Rat] {
	return atomBuilder[*big.Rat]{vtable: _VTableRat}
}
