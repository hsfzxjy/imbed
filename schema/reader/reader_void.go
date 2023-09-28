package schemareader

type Void struct{}

// Bool implements Reader.
func (Void) Bool() (bool, error) {
	panic("unimplemented")
}

// Float64 implements Reader.
func (Void) Float64() (float64, error) {
	panic("unimplemented")
}

// Int64 implements Reader.
func (Void) Int64() (int64, error) {
	panic("unimplemented")
}

// IterElem implements Reader.
func (Void) IterElem(func(i int, elem Reader) error) error {
	panic("unimplemented")
}

// IterKV implements Reader.
func (Void) IterKV(func(key string, value Reader) error) error {
	panic("unimplemented")
}

// IterField implements Reader.
func (Void) IterField(func(name string, field Reader) error) error {
	panic("unimplemented")
}

// ListSize implements Reader.
func (Void) ListSize() (int, error) {
	panic("unimplemented")
}

// MapSize implements Reader.
func (Void) MapSize() (int, error) {
	panic("unimplemented")
}

// String implements Reader.
func (Void) String() (string, error) {
	panic("unimplemented")
}

var _ Reader = Void{}
