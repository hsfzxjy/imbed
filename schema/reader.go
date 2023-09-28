package schema

type _AtomReader interface {
	Int64() (int64, error)
	Bool() (bool, error)
	String() (string, error)
	Float64() (float64, error)
}

type _MapReader interface {
	MapSize() (int, error)
	IterKV(func(key string, value Reader) error) error
}

type _ListReader interface {
	ListSize() (int, error)
	IterElem(func(i int, elem Reader) error) error
}

type _StructReader interface {
	IterField(func(name string, field Reader) error) error
}

type Reader interface {
	_StructReader
	_AtomReader
	_MapReader
	_ListReader
}
