package schemareader

type sliceReader[E any] struct {
	value []E
	Void
}

func NewSliceReader[E any](value []E) sliceReader[E] {
	return sliceReader[E]{value: value}
}

func (r sliceReader[E]) ListSize() (int, error) { return len(r.value), nil }
func (r sliceReader[E]) IterElem(f func(int, Reader) error) error {
	for i, e := range r.value {
		err := f(i, anyReader{e})
		if err != nil {
			return err
		}
	}
	return nil
}
