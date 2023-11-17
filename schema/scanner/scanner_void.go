package schemascanner

import (
	"math/big"

	"github.com/hsfzxjy/imbed/core/pos"
)

type Void struct{}

func (Void) Bool() (bool, pos.P, error) {
	panic("unimplemented")
}

func (Void) Rat() (*big.Rat, pos.P, error) {
	panic("unimplemented")
}

func (Void) Int64() (int64, pos.P, error) {
	panic("unimplemented")
}

func (Void) IterElem(func(i int, elem Scanner) error) error {
	panic("unimplemented")
}

func (Void) IterKV(func(key string, value Scanner) error) error {
	panic("unimplemented")
}

func (Void) IterField(func(name string, field Scanner, pos pos.P) error) error {
	panic("unimplemented")
}

func (Void) UnnamedField() Scanner {
	return nil
}

func (Void) ListSize() (int, pos.P, error) {
	panic("unimplemented")
}

func (Void) MapSize() (int, pos.P, error) {
	panic("unimplemented")
}

func (Void) String() (string, pos.P, error) {
	panic("unimplemented")
}

func (Void) Error(e error) error {
	return e
}
