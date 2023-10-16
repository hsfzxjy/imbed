package schemascanner

type sliceScanner[E any] struct {
	value []E
	Void
}

func NewSliceScanner[E any](value []E) sliceScanner[E] {
	return sliceScanner[E]{value: value}
}

func (r sliceScanner[E]) ListSize() (int, error) { return len(r.value), nil }
func (r sliceScanner[E]) IterElem(f func(int, Scanner) error) error {
	for i, e := range r.value {
		err := f(i, anyScanner{e})
		if err != nil {
			return err
		}
	}
	return nil
}
