package iter

import "github.com/hsfzxjy/imbed/core"

type chainedIt[T any] struct {
	its []core.Iterator[T]
}

func (c *chainedIt[T]) Next() (t T, ok bool) {
	for len(c.its) > 0 {
		it := c.its[0]
		t, ok = it.Next()
		if ok {
			return t, true
		}
		c.its = c.its[1:]
	}
	return t, false
}

func Chain[T any](its ...core.Iterator[T]) *chainedIt[T] {
	return &chainedIt[T]{its}
}
