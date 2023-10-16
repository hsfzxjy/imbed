package schemascanner

type mapScanner[V any] struct {
	m map[string]V
	Void
}

func NewMapScanner[V any](m map[string]V) mapScanner[V] {
	return mapScanner[V]{m: m}
}

func (m mapScanner[V]) IterKV(f func(string, Scanner) error) error {
	for k, v := range m.m {
		err := f(k, anyScanner{v})
		if err != nil {
			return err
		}
	}
	return nil
}

func (m mapScanner[V]) MapSize() (int, error) { return len(m.m), nil }
