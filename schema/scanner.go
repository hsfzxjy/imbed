package schema

import "math/big"

type _AtomScanner interface {
	Int64() (int64, error)
	Bool() (bool, error)
	String() (string, error)
	Rat() (*big.Rat, error)
}

type _MapScanner interface {
	MapSize() (int, error)
	IterKV(func(key string, value Scanner) error) error
}

type _ListScanner interface {
	ListSize() (int, error)
	IterElem(func(i int, elem Scanner) error) error
}

type _StructScanner interface {
	IterField(func(name string, field Scanner) error) error
	UnnamedField() (fieldScanner Scanner)
}

type Scanner interface {
	_StructScanner
	_AtomScanner
	_MapScanner
	_ListScanner
	Error(err error) error
	Snapshot() any
	Reset(snapshot any)
}
