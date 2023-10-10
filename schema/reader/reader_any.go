package schemareader

import (
	"fmt"
	"math/big"

	"github.com/hsfzxjy/imbed/schema"
)

type anyReader struct{ value any }

func (r anyReader) Bool() (bool, error) {
	switch v := r.value.(type) {
	case bool:
		return v, nil
	case nil:
		return false, schema.ErrRequired
	default:
		return false, wrongType(v, "bool")
	}
}

func (r anyReader) Rat() (*big.Rat, error) {
	switch v := r.value.(type) {
	case float64:
		return new(big.Rat).SetFloat64(v), nil
	case string:
		r, ok := new(big.Rat).SetString(v)
		if !ok {
			return nil, fmt.Errorf("invalid rat %q", v)
		}
		return r, nil
	case nil:
		return nil, schema.ErrRequired
	default:
		return nil, wrongType(v, "rat")
	}
}

func (r anyReader) Int64() (int64, error) {
	switch v := r.value.(type) {
	case int64:
		return v, nil
	case nil:
		return 0, schema.ErrRequired
	default:
		return 0, wrongType(v, "int64")
	}
}

func (r anyReader) String() (string, error) {
	switch v := r.value.(type) {
	case string:
		return v, nil
	case nil:
		return "", schema.ErrRequired
	default:
		return "", wrongType(v, "string")
	}
}

func (r anyReader) IterElem(f func(i int, elem Reader) error) error {
	switch v := r.value.(type) {
	case []any:
		return NewSliceReader(v).IterElem(f)
	case nil:
		return nil
	default:
		return wrongType(v, "[]any")
	}
}

func (r anyReader) IterField(f func(name string, field Reader) error) error {
	switch v := r.value.(type) {
	case map[string]any:
		return NewMapReader(v).IterKV(f)
	default:
		return nil
	}
}

func (r anyReader) IterKV(f func(key string, value Reader) error) error {
	switch v := r.value.(type) {
	case map[string]any:
		return NewMapReader(v).IterKV(f)
	default:
		return wrongType(v, "map[string]any")
	}
}

func (r anyReader) ListSize() (int, error) {
	switch v := r.value.(type) {
	case []any:
		return NewSliceReader(v).ListSize()
	case nil:
		return 0, schema.ErrRequired
	default:
		return 0, wrongType(v, "[]any")
	}
}

func (r anyReader) MapSize() (int, error) {
	switch v := r.value.(type) {
	case map[string]any:
		return NewMapReader(v).MapSize()
	case nil:
		return 0, schema.ErrRequired
	default:
		return 0, wrongType(v, "map[string]any")
	}
}

func (r anyReader) Error(e error) error { return e }

func _() { var _ Reader = anyReader{} }

func Any(value any) Reader { return anyReader{value} }
