package schema

type atomBuilder[T any] struct {
	def    optional[T]
	vtable *_AtomVTable[T]
}

func (b atomBuilder[T]) Default(value T) atomBuilder[T] {
	b.def.Value = value
	b.def.IsValid = true
	return b
}

func (b atomBuilder[T]) buildSchema() schemaTyped[T] { return &_Atom[T]{b.def, b.vtable} }
func (b atomBuilder[T]) buildSchemaUntyped() schema  { return b.buildSchema() }

func Int() atomBuilder[int64] {
	return atomBuilder[int64]{vtable: _VTableInt}
}

func Bool() atomBuilder[bool] {
	return atomBuilder[bool]{vtable: _VTableBool}
}

func String() atomBuilder[string] {
	return atomBuilder[string]{vtable: _VTableString}
}

func Float() atomBuilder[float64] {
	return atomBuilder[float64]{vtable: _VTableFloat}
}
