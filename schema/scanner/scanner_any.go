package schemascanner

import (
	"fmt"
	"math/big"

	"github.com/hsfzxjy/imbed/schema"
)

type anyScanner struct{ value any }

func (r anyScanner) Bool() (bool, error) {
	switch v := r.value.(type) {
	case bool:
		return v, nil
	case nil:
		return false, schema.ErrRequired
	default:
		return false, wrongType(v, "bool")
	}
}

func (r anyScanner) Rat() (*big.Rat, error) {
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

func (r anyScanner) Int64() (int64, error) {
	switch v := r.value.(type) {
	case int64:
		return v, nil
	case nil:
		return 0, schema.ErrRequired
	default:
		return 0, wrongType(v, "int64")
	}
}

func (r anyScanner) String() (string, error) {
	switch v := r.value.(type) {
	case string:
		return v, nil
	case nil:
		return "", schema.ErrRequired
	default:
		return "", wrongType(v, "string")
	}
}

func (r anyScanner) IterElem(f func(i int, elem Scanner) error) error {
	switch v := r.value.(type) {
	case []any:
		return NewSliceScanner(v).IterElem(f)
	case nil:
		return nil
	default:
		return wrongType(v, "[]any")
	}
}

func (r anyScanner) IterField(f func(name string, field Scanner) error) error {
	switch v := r.value.(type) {
	case map[string]any:
		return NewMapScanner(v).IterKV(f)
	default:
		return nil
	}
}

func (r anyScanner) UnnamedField() Scanner {
	return nil
}

func (r anyScanner) IterKV(f func(key string, value Scanner) error) error {
	switch v := r.value.(type) {
	case map[string]any:
		return NewMapScanner(v).IterKV(f)
	default:
		return wrongType(v, "map[string]any")
	}
}

func (r anyScanner) ListSize() (int, error) {
	switch v := r.value.(type) {
	case []any:
		return NewSliceScanner(v).ListSize()
	case nil:
		return 0, schema.ErrRequired
	default:
		return 0, wrongType(v, "[]any")
	}
}

func (r anyScanner) MapSize() (int, error) {
	switch v := r.value.(type) {
	case map[string]any:
		return NewMapScanner(v).MapSize()
	case nil:
		return 0, schema.ErrRequired
	default:
		return 0, wrongType(v, "map[string]any")
	}
}

func (r anyScanner) Error(e error) error { return e }

func (anyScanner) Snapshot() any {
	return nil
}

func (anyScanner) Reset(snapshot any) {
}

func _() { var _ Scanner = anyScanner{} }

func Any(value any) Scanner { return anyScanner{value} }
