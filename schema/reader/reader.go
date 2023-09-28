package schemareader

import (
	"fmt"

	"github.com/hsfzxjy/imbed/schema"
)

func wrongType(v any, expected string) error {
	return fmt.Errorf("expect type %s, got %T (value=%v)", expected, v, v)
}

type Reader = schema.Reader
