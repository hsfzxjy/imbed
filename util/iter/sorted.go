package iter

import (
	"cmp"
	"slices"

	"github.com/hsfzxjy/imbed/core"
)

type cmpFunc[K any] func(k1, k2 K) int

type sortedIt[T any, It core.Iterator[T]] struct {
	it      It
	cmpFunc cmpFunc[T]
	sliceIt sliceIt[T]
}

func (i *sortedIt[T, It]) eval() {
	if i.cmpFunc == nil {
		return
	}
	cmpFunc := i.cmpFunc
	i.cmpFunc = nil
	var s []T
	for i.it.HasNext() {
		s = append(s, i.it.Next())
	}
	slices.SortFunc(s, cmpFunc)
	i.sliceIt.slice = s
}

func (i *sortedIt[T, It]) HasNext() bool {
	i.eval()
	return i.sliceIt.HasNext()
}

func (i *sortedIt[T, It]) Next() (result T) {
	if !i.HasNext() {
		return
	}
	return i.sliceIt.Next()
}

func SortedKeyFunc[T any, It core.Iterator[T], K cmp.Ordered](
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

func Sorted[T any, It core.Iterator[T]](
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
