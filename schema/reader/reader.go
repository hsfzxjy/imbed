package schemareader

import (
	"fmt"

	"github.com/hsfzxjy/imbed/schema"
)

type ErrWrongType struct {
	Expected string
	Value    any
}

func (e *ErrWrongType) Error() string {
	return fmt.Sprintf("expect type %s, got %T (value=%v)", e.Expected, e.Value, e.Value)
}

func wrongType(v any, expected string) error {
	return &ErrWrongType{expected, v}
}

type Reader = schema.Reader
