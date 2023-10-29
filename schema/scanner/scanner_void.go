package schemascanner

import "math/big"

type Void struct{}

func (Void) Bool() (bool, error) {
	panic("unimplemented")
}

func (Void) Rat() (*big.Rat, error) {
	panic("unimplemented")
}

func (Void) Int64() (int64, error) {
	panic("unimplemented")
}

func (Void) IterElem(func(i int, elem Scanner) error) error {
	panic("unimplemented")
}

func (Void) IterKV(func(key string, value Scanner) error) error {
	panic("unimplemented")
}

func (Void) IterField(func(name string, field Scanner) error) error {
	panic("unimplemented")
}

func (Void) UnnamedField() Scanner {
	return nil
}

func (Void) ListSize() (int, error) {
	panic("unimplemented")
}

func (Void) MapSize() (int, error) {
	panic("unimplemented")
}

func (Void) String() (string, error) {
	panic("unimplemented")
}

func (Void) Error(e error) error {
	return e
}
