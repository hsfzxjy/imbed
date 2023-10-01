package iter

import "github.com/hsfzxjy/imbed/core"

type mappedIt[T, U any] struct {
	core.Iterator[T]
	mapFunc func(T) U
}

func (m *mappedIt[T, U]) Current() U {
	return m.mapFunc(m.Iterator.Current())
}

func Map[T, U any](it core.Iterator[T], mapFunc func(T) U) core.Iterator[U] {
	return &mappedIt[T, U]{it, mapFunc}
}
