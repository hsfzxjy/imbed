package iter

import "github.com/hsfzxjy/tipe"

type stop struct{}

func (stop) Error() string { return "iter: stopped" }

var Stop = error(stop{})

type Ator[T any] interface {
	HasNext() bool
	Nexter[T]
}

type Nexter[T any] interface {
	Next() tipe.Result[T]
}

func Stopped[T any](r tipe.Result[T]) bool {
	return r.IsErr() && r.UnwrapErr() == Stop
}

type nexterFunc[T any] func() tipe.Result[T]

func (f nexterFunc[T]) Next() tipe.Result[T] {
	return f()
}

func NewFunc[T any](f func() tipe.Result[T]) Nexter[T] {
	return nexterFunc[T](f)
}
