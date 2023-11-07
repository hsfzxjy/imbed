package iter

import "github.com/hsfzxjy/tipe"

type chainIt[T any] struct {
	first bool
	its   []Nexter[T]
	t     tipe.Result[T]
}

func (i *chainIt[T]) HasNext() bool {
	if i.first {
		i.first = false
		i.next()
	}
	return len(i.its) > 0
}

func (i *chainIt[T]) next() {
	if !i.HasNext() {
		return
	}
	for len(i.its) > 0 {
		it := i.its[0]
		t := it.Next()
		if Stopped(t) {
			i.its = i.its[1:]
			continue
		}
		i.t = t
		break
	}
}

func (i *chainIt[T]) Next() (t tipe.Result[T]) {
	if !i.HasNext() {
		return tipe.Err[T](Stop)
	}
	t = i.t
	i.next()
	return t
}

func Chain[T any](its ...Nexter[T]) *chainIt[T] {
	return &chainIt[T]{
		first: true,
		its:   its,
	}
}
