package iter

import (
	"errors"

	"github.com/hsfzxjy/imbed/core"
)

var (
	ErrEmpty   = errors.New("iterator has length == 0")
	ErrTooMany = errors.New("iterator has length >= 2")
)

func One[T any, It core.Iterator[T]](it It) (T, error) {
	var t, zero T
	if !it.HasNext() {
		return zero, ErrEmpty
	}
	t = it.Next()
	if it.HasNext() {
		return zero, ErrTooMany
	}
	return t, nil
}

func One2[T any, It core.Iterator[T]](it It, err error) (T, error) {
	var zero T
	if err != nil {
		return zero, err
	}
	return One(it)
}
