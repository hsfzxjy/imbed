package internal

import (
	"github.com/hsfzxjy/imbed/db/internal/helper"
)

type H struct {
	*helper.Helper
}

func (h H) runR(f func(H) error) error  { return f(h) }
func (h H) runRW(f func(H) error) error { return f(h) }
func (h H) DB() Service                 { return Service{h.BBoltDB()} }
