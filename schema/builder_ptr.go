package schema

type ptrBuilder[T any] struct {
	def         optional[*T]
	elemBuilder builder[T]
}

func (b ptrBuilder[T]) Default(def *T) ptrBuilder[T] {
	b.def.IsValid = true
	b.def.Value = def
	return b
}

func (b ptrBuilder[T]) buildSchema() schema[*T] {
	return &_Ptr[T]{
		def:        b.def,
		elemSchema: b.elemBuilder.buildSchema(),
	}
}

func (b ptrBuilder[T]) buildGenericSchema() genericSchema { return b.buildSchema() }

func Ptr[T any](elemBuilder builder[T]) ptrBuilder[T] {
	return ptrBuilder[T]{elemBuilder: elemBuilder}
}
