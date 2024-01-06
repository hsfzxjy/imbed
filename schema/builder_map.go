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
	return nil
}

func (b mapBuilder[V]) buildGenericSchema() genericSchema {
	return &_Map{b.valueBuilder.buildSchema(), goMapProto[V]{}, goMapDefaulter[V]{b.defFunc}}
}

func Map[V any](valueBuilder builder[V]) mapBuilder[V] {
	return mapBuilder[V]{valueBuilder: valueBuilder}
}
