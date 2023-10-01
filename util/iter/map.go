package iter

import "github.com/hsfzxjy/imbed/core"

type mappedIt[T, U any, It core.Iterator[T]] struct {
	iterator   It
	current    U
	mapFunc    func(T) (U, bool)
	exhaustive bool
}

func (m *mappedIt[T, U, It]) Current() U {
	if m.exhaustive {
		var u U
		return u
	}
	return m.current
}

func (m *mappedIt[T, U, It]) Exhausted() bool {
	return m.exhaustive
}

func (m *mappedIt[T, U, It]) Next() {
	if m.exhaustive {
		return
	}
	m.iterator.Next()
	m.setCurrent()
}

func (m *mappedIt[T, U, It]) setCurrent() {
	if m.iterator.Exhausted() {
		m.exhaustive = true
	} else {
		current, ok := m.mapFunc(m.iterator.Current())
		if !ok {
			m.exhaustive = true
		} else {
			m.exhaustive = false
			m.current = current
		}
	}
}

func Map[T, U any, It core.Iterator[T]](it It, mapFunc func(T) (U, bool)) *mappedIt[T, U, It] {
	m := &mappedIt[T, U, It]{iterator: it, mapFunc: mapFunc}
	m.setCurrent()
	return m
}

type flatMappedIt[T, U any, It1 core.Iterator[T], It2 core.Iterator[It1]] struct {
	it2           It2
	currentU      U
	mapFunc       func(T) (U, bool)
	it2Exhaustive bool
}

func (m *flatMappedIt[T, U, It1, It2]) setCurrent() {
	if m.it2Exhaustive {
		return
	}
	if m.it2.Exhausted() {
		m.it2Exhaustive = true
		return
	}
	it1 := m.it2.Current()
	for it1.Exhausted() {
		m.it2.Next()
		if m.it2.Exhausted() {
			m.it2Exhaustive = true
			return
		}
		it1 = m.it2.Current()
	}
	current, ok := m.mapFunc(it1.Current())
	if !ok {
		m.it2Exhaustive = true
		return
	}
	m.currentU = current
}

func (m *flatMappedIt[T, U, It1, It2]) Current() U {
	return m.currentU
}
func (m *flatMappedIt[T, U, It1, It2]) Exhausted() bool {
	return m.it2Exhaustive
}
func (m *flatMappedIt[T, U, It1, It2]) Next() {
	if m.it2Exhaustive {
		return
	}
	it1 := m.it2.Current()
	it1.Next()
	m.setCurrent()
}

func FlatMap[T, U any, It1 core.Iterator[T], It2 core.Iterator[It1]](it It2, mapFunc func(T) (U, bool)) *flatMappedIt[T, U, It1, It2] {
	m := &flatMappedIt[T, U, It1, It2]{it2: it, mapFunc: mapFunc}
	m.setCurrent()
	return m
}
