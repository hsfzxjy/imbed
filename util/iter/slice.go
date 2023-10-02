package iter

import "github.com/hsfzxjy/imbed/core"

type sliceIt[T any] struct {
	slice []T
}

func (i *sliceIt[T]) Next() (t T, ok bool) {
	if len(i.slice) > 0 {
		result := i.slice[0]
		i.slice = i.slice[1:]
		return result, true
	}
	return t, false
}

func (i *sliceIt[T]) Chain(its ...core.Iterator[T]) *chainedIt[T] {
	var slice = make([]core.Iterator[T], 0, len(its)+1)
	slice = append(slice, i)
	slice = append(slice, its...)
	return &chainedIt[T]{slice}
}

func Slice[T any](data ...T) *sliceIt[T] {
	return &sliceIt[T]{slice: data}
}
