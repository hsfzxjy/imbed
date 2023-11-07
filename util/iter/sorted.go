package iter

import (
	"cmp"
	"slices"

	"github.com/hsfzxjy/tipe"
)

type cmpFunc[K any] func(k1, k2 K) int

type sortedIt[T any, It Nexter[T]] struct {
	it      It
	cmpFunc cmpFunc[T]
	sliceIt tipe.Result[*sliceIt[T]]
}

func (i *sortedIt[T, It]) eval() {
	if i.cmpFunc == nil {
		return
	}
	cmpFunc := i.cmpFunc
	i.cmpFunc = nil
	var s []T
	for {
		t := i.it.Next()
		if Stopped(t) {
			break
		}
		if t.IsErr() {
			i.sliceIt = i.sliceIt.FillErr(t.UnwrapErr())
			return
		}
		s = append(s, t.Unwrap())
	}
	slices.SortFunc(s, cmpFunc)
	i.sliceIt = i.sliceIt.Fill(&sliceIt[T]{slice: s})
}

func (i *sortedIt[T, It]) HasNext() bool {
	i.eval()
	return !Stopped(i.sliceIt) && (i.sliceIt.IsErr() || i.sliceIt.Unwrap().HasNext())
}

func (i *sortedIt[T, It]) Next() (result tipe.Result[T]) {
	if !i.HasNext() {
		return result.FillErr(Stop)
	}
	if i.sliceIt.IsErr() {
		result = result.FillErr(i.sliceIt.UnwrapErr())
		i.sliceIt = i.sliceIt.FillErr(Stop)
		return
	}
	return i.sliceIt.Unwrap().Next()
}

func SortedKeyFunc[T any, It Nexter[T], K cmp.Ordered](
	it It,
	keyFunc func(T) K,
	reversed bool,
) *sortedIt[T, It] {
	return &sortedIt[T, It]{
		it: it, cmpFunc: func(a, b T) int {
			r := cmp.Compare(keyFunc(a), keyFunc(b))
			if reversed {
				return -r
			} else {
				return r
			}
		}}
}

func Sorted[T any, It Nexter[T]](
	it It,
	cmpFunc cmpFunc[T],
	reversed bool,
) *sortedIt[T, It] {
	if reversed {
		oldCmpFunc := cmpFunc
		cmpFunc = func(a, b T) int {
			return -oldCmpFunc(a, b)
		}
	}
	return &sortedIt[T, It]{it: it, cmpFunc: cmpFunc}
}
