package schema

import (
	"math/big"

	"github.com/hsfzxjy/imbed/core/pos"
)

type _AtomScanner interface {
	Int64() (int64, pos.P, error)
	Bool() (bool, pos.P, error)
	String() (string, pos.P, error)
	Rat() (*big.Rat, pos.P, error)
}

type _MapScanner interface {
	MapSize() (int, pos.P, error)
	IterKV(func(key string, value Scanner) error) error
}

type _ListScanner interface {
	ListSize() (int, pos.P, error)
	IterElem(func(i int, elem Scanner) error) error
}

type _StructScanner interface {
	IterField(func(name string, field Scanner, pos pos.P) error) error
	UnnamedField() (fieldScanner Scanner)
}

type Scanner interface {
	_StructScanner
	_AtomScanner
	_MapScanner
	_ListScanner
	Snapshot() any
	Reset(snapshot any)
}
