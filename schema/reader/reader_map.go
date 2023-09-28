package schemareader

type mapReader[V any] struct {
	m map[string]V
	Void
}

func NewMapReader[V any](m map[string]V) mapReader[V] {
	return mapReader[V]{m: m}
}

func (m mapReader[V]) IterKV(f func(string, Reader) error) error {
	for k, v := range m.m {
		err := f(k, anyReader{v})
		if err != nil {
			return err
		}
	}
	return nil
}

func (m mapReader[V]) MapSize() (int, error) { return len(m.m), nil }
