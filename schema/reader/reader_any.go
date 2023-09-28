package schemareader

import "github.com/hsfzxjy/imbed/schema"

type anyReader struct{ value any }

// Bool implements Reader.
func (r anyReader) Bool() (bool, error) {
	switch v := r.value.(type) {
	case bool:
		return v, nil
	case nil:
		return false, schema.ErrRequired
	default:
		return false, wrongType(v)
	}
}

// Float64 implements Reader.
func (r anyReader) Float64() (float64, error) {
	switch v := r.value.(type) {
	case float64:
		return v, nil
	case nil:
		return 0, schema.ErrRequired
	default:
		return 0, wrongType(v)
	}
}

// Int64 implements Reader.
func (r anyReader) Int64() (int64, error) {
	switch v := r.value.(type) {
	case int64:
		return v, nil
	case nil:
		return 0, schema.ErrRequired
	default:
		return 0, wrongType(v)
	}
}

// String implements Reader.
func (r anyReader) String() (string, error) {
	switch v := r.value.(type) {
	case string:
		return v, nil
	case nil:
		return "", schema.ErrRequired
	default:
		return "", wrongType(v)
	}
}

// IterElem implements Reader.
func (r anyReader) IterElem(f func(i int, elem Reader) error) error {
	switch v := r.value.(type) {
	case []any:
		return NewSliceReader(v).IterElem(f)
	case nil:
		return nil
	default:
		return wrongType(v)
	}
}

// IterField implements Reader.
func (r anyReader) IterField(f func(name string, field Reader) error) error {
	switch v := r.value.(type) {
	case map[string]any:
		return NewMapReader(v).IterKV(f)
	default:
		return nil
	}
}

// IterKV implements Reader.
func (r anyReader) IterKV(f func(key string, value Reader) error) error {
	switch v := r.value.(type) {
	case map[string]any:
		return NewMapReader(v).IterKV(f)
	default:
		return wrongType(v)
	}
}

// ListSize implements Reader.
func (r anyReader) ListSize() (int, error) {
	switch v := r.value.(type) {
	case []any:
		return NewSliceReader(v).ListSize()
	case nil:
		return 0, schema.ErrRequired
	default:
		return 0, wrongType(v)
	}
}

// MapSize implements Reader.
func (r anyReader) MapSize() (int, error) {
	switch v := r.value.(type) {
	case map[string]any:
		return NewMapReader(v).MapSize()
	case nil:
		return 0, schema.ErrRequired
	default:
		return 0, wrongType(v)
	}
}

var _ Reader = anyReader{}

func Any(value any) Reader { return anyReader{value} }
