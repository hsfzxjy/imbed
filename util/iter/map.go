package iter

import "github.com/hsfzxjy/imbed/core"

type mapFunc[T, U any] func(T) (U, bool)

type mappedIt[T, U any, It core.Iterator[T]] struct {
	iterator It
	mapFunc  mapFunc[T, U]
	stopped  bool
	first    bool
	result   U
}

func (m *mappedIt[T, U, It]) HasNext() bool {
	if m.first {
		m.first = false
		m.next()
	}
	return !m.stopped
}

func (m *mappedIt[T, U, It]) Next() (result U) {
	if !m.HasNext() {
		return
	}
	result = m.result
	m.next()
	return result
}

func (m *mappedIt[T, U, It]) next() {
	if m.stopped || !m.iterator.HasNext() {
		m.stopped = true
		return
	}
	t := m.iterator.Next()
	u, ok := m.mapFunc(t)
	if !ok {
		m.stopped = true
		return
	}
	m.result = u
}

func Map[T, U any, It core.Iterator[T]](it It, mapFunc mapFunc[T, U]) *mappedIt[T, U, It] {
	m := &mappedIt[T, U, It]{iterator: it, mapFunc: mapFunc, first: true}
	return m
}

type flatMappedIt[T, U any, It1 core.Iterator[T], It2 core.Iterator[It1]] struct {
	it2     It2
	it1     It1
	mapFunc mapFunc[T, U]
	first   bool
	stopped bool
	result  U
}

func (m *flatMappedIt[T, U, It1, It2]) HasNext() bool {
	if m.first {
		m.first = false
		m.next(true)
	}
	return !m.stopped
}

func (m *flatMappedIt[T, U, It1, It2]) next(init bool) {
	if m.stopped {
		return
	}
	var it1 It1
	if init {
		if !m.it2.HasNext() {
			m.stopped = true
			return
		}
		it1 = m.it2.Next()
	} else {
		it1 = m.it1
	}
	for !it1.HasNext() && m.it2.HasNext() {
		it1 = m.it2.Next()
	}
	if !it1.HasNext() {
		m.stopped = true
		return
	}
	t := it1.Next()
	u, ok := m.mapFunc(t)
	if !ok {
		m.stopped = true
		return
	}
	m.result = u
}

func (m *flatMappedIt[T, U, It1, It2]) Next() (result U) {
	if !m.HasNext() {
		return
	}
	result = m.result
	m.next(false)
	return result
}

func FlatMap[T, U any, It1 core.Iterator[T], It2 core.Iterator[It1]](it2 It2, mapFunc func(T) (U, bool)) *flatMappedIt[T, U, It1, It2] {
	m := &flatMappedIt[T, U, It1, It2]{it2: it2, mapFunc: mapFunc, first: true}
	return m
}
