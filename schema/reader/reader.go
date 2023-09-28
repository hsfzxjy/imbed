package schemareader

import (
	"fmt"

	"github.com/hsfzxjy/imbed/schema"
)

func wrongType(v any) error {
	return fmt.Errorf("unexpected type: %T", v)
}

type Reader = schema.Reader
