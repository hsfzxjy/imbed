package schemascanner

import (
	"fmt"
	"math/big"

	"github.com/hsfzxjy/imbed/core/pos"
	"github.com/hsfzxjy/imbed/schema"
)

type anyScanner struct{ value any }

func (r anyScanner) Bool() (bool, pos.P, error) {
	switch v := r.value.(type) {
	case bool:
		return v, pos.P{}, nil
	case nil:
		return false, pos.P{}, schema.ErrRequired
	default:
		return false, pos.P{}, wrongType(v, "bool")
	}
}

func (r anyScanner) Rat() (*big.Rat, pos.P, error) {
	switch v := r.value.(type) {
	case float64:
		return new(big.Rat).SetFloat64(v), pos.P{}, nil
	case string:
		r, ok := new(big.Rat).SetString(v)
		if !ok {
			return nil, pos.P{}, fmt.Errorf("invalid rat %q", v)
		}
		return r, pos.P{}, nil
	case nil:
		return nil, pos.P{}, schema.ErrRequired
	default:
		return nil, pos.P{}, wrongType(v, "rat")
	}
}

func (r anyScanner) Int64() (int64, pos.P, error) {
	switch v := r.value.(type) {
	case int64:
		return v, pos.P{}, nil
	case nil:
		return 0, pos.P{}, schema.ErrRequired
	default:
		return 0, pos.P{}, wrongType(v, "int64")
	}
}

func (r anyScanner) String() (string, pos.P, error) {
	switch v := r.value.(type) {
	case string:
		return v, pos.P{}, nil
	case nil:
		return "", pos.P{}, schema.ErrRequired
	default:
		return "", pos.P{}, wrongType(v, "string")
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

func (r anyScanner) IterField(f func(name string, field Scanner, pos pos.P) error) error {
	switch v := r.value.(type) {
	case map[string]any:
		return NewMapScanner(v).IterKV(func(key string, s schema.Scanner) error {
			return f(key, s, pos.P{})
		})
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

func (r anyScanner) ListSize() (int, pos.P, error) {
	switch v := r.value.(type) {
	case []any:
		return NewSliceScanner(v).ListSize()
	case nil:
		return 0, pos.P{}, schema.ErrRequired
	default:
		return 0, pos.P{}, wrongType(v, "[]any")
	}
}

func (r anyScanner) MapSize() (int, pos.P, error) {
	switch v := r.value.(type) {
	case map[string]any:
		return NewMapScanner(v).MapSize()
	case nil:
		return 0, pos.P{}, schema.ErrRequired
	default:
		return 0, pos.P{}, wrongType(v, "map[string]any")
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
