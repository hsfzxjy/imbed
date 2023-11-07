package iter

import (
	"errors"

	"github.com/hsfzxjy/tipe"
)

var (
	ErrEmpty   = errors.New("iterator has length == 0")
	ErrTooMany = errors.New("iterator has length >= 2")
)

func One[T any, It Nexter[T]](it It) tipe.Result[T] {
	t := it.Next()
	if Stopped(t) {
		return t.FillErr(ErrEmpty)
	}

	if !Stopped(it.Next()) {
		return t.FillErr(ErrTooMany)
	}

	return t
}
