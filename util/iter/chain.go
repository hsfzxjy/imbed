package iter

import "github.com/hsfzxjy/imbed/core"

type chainedIt[T any] struct {
	its []core.Iterator[T]
}

func (c *chainedIt[T]) HasNext() bool {
	for len(c.its) > 0 {
		if c.its[0].HasNext() {
			return true
		}
		c.its = c.its[1:]
	}
	return false
}

func (c *chainedIt[T]) Next() (t T) {
	if !c.HasNext() {
		return
	}
	return c.its[0].Next()
}

func Chain[T any](its ...core.Iterator[T]) *chainedIt[T] {
	return &chainedIt[T]{its}
}
