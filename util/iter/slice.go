package iter

import "github.com/hsfzxjy/tipe"

type sliceIt[T any] struct {
	slice []T
}

func (i *sliceIt[T]) HasNext() bool { return len(i.slice) > 0 }

func (i *sliceIt[T]) Next() (t tipe.Result[T]) {
	if i.HasNext() {
		result := i.slice[0]
		i.slice = i.slice[1:]
		return tipe.Ok(result)
	} else {
		return tipe.Err[T](Stop)
	}
}

func (i *sliceIt[T]) Chain(its ...Nexter[T]) *chainIt[T] {
	var slice = make([]Nexter[T], 0, len(its)+1)
	slice = append(slice, i)
	slice = append(slice, its...)
	return Chain(slice...)
}

func Slice[T any](data ...T) *sliceIt[T] {
	return &sliceIt[T]{slice: data}
}
