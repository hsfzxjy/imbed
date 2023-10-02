package iter

import "github.com/hsfzxjy/imbed/core"

type mapFunc[T, U any] func(T) (U, bool)

type mappedIt[T, U any, It core.Iterator[T]] struct {
	iterator It
	mapFunc  mapFunc[T, U]
	stopped  bool
}

func (m *mappedIt[T, U, It]) Next() (result U, ok bool) {
	var t T
	var u U
	if m.stopped {
		return u, false
	}
	t, ok = m.iterator.Next()
	if !ok {
		return u, ok
	}
	u, ok = m.mapFunc(t)
	if !ok {
		m.stopped = true
		return u, ok
	}
	return u, true
}

func Map[T, U any, It core.Iterator[T]](it It, mapFunc mapFunc[T, U]) *mappedIt[T, U, It] {
	m := &mappedIt[T, U, It]{iterator: it, mapFunc: mapFunc}
	return m
}

type flatMappedIt[T, U any, It1 core.Iterator[T], It2 core.Iterator[It1]] struct {
	it2      It2
	it1      It1
	mapFunc  mapFunc[T, U]
	it1Ready bool
	stopped  bool
}

func (m *flatMappedIt[T, U, It1, It2]) Next() (U, bool) {
	var u U
	var t T
	if m.stopped {
		return u, false
	}

	if !m.it1Ready {
		var ok bool
		var it1 It1
		it1, ok = m.it2.Next()
		if !ok {
			m.stopped = true
			return u, false
		}
		m.it1Ready = true
		m.it1 = it1
	}

	for {
		var ok bool
		t, ok = m.it1.Next()
		if ok {
			u, ok = m.mapFunc(t)
			if !ok {
				m.stopped = true
			}
			return u, ok
		}
		m.it1, ok = m.it2.Next()
		if !ok {
			m.stopped = true
			return u, false
		}
	}

}

func FlatMap[T, U any, It1 core.Iterator[T], It2 core.Iterator[It1]](it2 It2, mapFunc func(T) (U, bool)) *flatMappedIt[T, U, It1, It2] {
	m := &flatMappedIt[T, U, It1, It2]{it2: it2, mapFunc: mapFunc}
	return m
}
