package iter

import (
	"github.com/hsfzxjy/tipe"
)

type mapFunc[T, U any] func(T) tipe.Result[U]

type mappedIt[T, U any, It Nexter[T]] struct {
	iterator It
	mapFunc  mapFunc[T, U]
	stopped  bool
	first    bool
	u        tipe.Result[U]
}

func (m *mappedIt[T, U, It]) HasNext() bool {
	if m.first {
		m.first = false
		m.next()
	}
	return !m.stopped
}

func (m *mappedIt[T, U, It]) Next() (result tipe.Result[U]) {
	if !m.HasNext() {
		return result.FillErr(Stop)
	}
	result = m.u
	m.next()
	return result
}

func (m *mappedIt[T, U, It]) next() {
	if m.stopped {
		return
	}
	if m.u.IsErr() {
		m.stopped = true
		return
	}
	t := m.iterator.Next()
	u := tipe.BindR(t, m.mapFunc)
	if Stopped(u) {
		m.stopped = true
		return
	}
	m.u = u
}

func FilterMap[T, U any, It Nexter[T]](it It, mapFunc mapFunc[T, U]) *mappedIt[T, U, It] {
	m := &mappedIt[T, U, It]{iterator: it, mapFunc: mapFunc, first: true}
	return m
}

type flatMappedIt[T, U any, It1 Nexter[T], It2 Nexter[It1]] struct {
	it2      It2
	it1      It1
	u        tipe.Result[U]
	stopped  bool
	it1Valid bool
	first    bool
	mapFunc  mapFunc[T, U]
}

func (m *flatMappedIt[T, U, It1, It2]) HasNext() bool {
	if m.first {
		m.first = false
		m.next()
	}
	return !m.stopped
}

func (m *flatMappedIt[T, U, It1, It2]) next() {
	if m.stopped {
		return
	}
	if m.u.IsErr() {
		m.stopped = true
		return
	}
	if m.it1Valid {
		goto FETCH_U
	}
FETCH_IT1:
	{
		out := m.it2.Next()
		if Stopped(out) {
			m.stopped = true
			return
		}
		if out.IsErr() {
			m.it1Valid = false
			m.u = tipe.Err[U](out.UnwrapErr())
			return
		}
		m.it1 = out.Unwrap()
		m.it1Valid = true
	}
FETCH_U:
	t := m.it1.Next()
	if Stopped(t) {
		goto FETCH_IT1
	}
	u := tipe.BindR(t, m.mapFunc)
	if Stopped(u) {
		m.stopped = true
		return
	}
	m.u = u
}

func (m *flatMappedIt[T, U, It1, It2]) Next() (result tipe.Result[U]) {
	if !m.HasNext() {
		return result.FillErr(Stop)
	}
	result = m.u
	m.next()
	return result
}

func FlatFilterMap[T, U any, It1 Nexter[T], It2 Nexter[It1]](it2 It2, mapFunc mapFunc[T, U]) *flatMappedIt[T, U, It1, It2] {
	m := &flatMappedIt[T, U, It1, It2]{
		it2:     it2,
		mapFunc: mapFunc,
		first:   true,
	}
	return m
}
