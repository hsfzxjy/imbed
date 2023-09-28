package schema

type mapBuilder[V any] struct {
	defFunc      func() map[string]V
	valueBuilder builder[V]
}

func (b mapBuilder[V]) Default(defFunc func() map[string]V) mapBuilder[V] {
	b.defFunc = defFunc
	return b
}

func (b mapBuilder[V]) buildSchema() schema[map[string]V] {
	return &_Map[V]{b.defFunc, b.valueBuilder.buildSchema()}
}

func (b mapBuilder[V]) buildGenericSchema() genericSchema { return b.buildSchema() }

func Map[V any](valueBuilder builder[V]) mapBuilder[V] {
	return mapBuilder[V]{valueBuilder: valueBuilder}
}
