package iterator

import (
	"github.com/hsfzxjy/imbed/db/internal"
	"go.etcd.io/bbolt"
)

type It[T any] struct {
	h         internal.H
	cursor    *bbolt.Cursor
	current   *T
	getObject func(k, v []byte) *T
}

type Builder[T any] struct {
	Cursor    *bbolt.Cursor
	SeekTo    []byte
	GetObject func(k, v []byte) *T
}

func New[T any](constructor func(h internal.H) (*Builder[T], error)) internal.Runnable[*It[T]] {
	return internal.R[*It[T]](func(h internal.H) (*It[T], error) {
		b, err := constructor(h)
		if err != nil {
			return nil, err
		}
		if b == nil {
			return nil, nil
		}
		cursor := b.Cursor
		seekTo := b.SeekTo
		getObject := b.GetObject
		it := new(It[T])
		it.h = h
		it.cursor = cursor
		it.getObject = getObject
		var k, v []byte
		if b.SeekTo != nil {
			k, v = cursor.Seek(seekTo)
		} else {
			k, v = cursor.First()
		}
		if k != nil {
			it.current = getObject(k, v)
		}
		return it, nil
	})
}

func (it *It[T]) Current() *T {
	if it == nil {
		return nil
	}
	return it.current
}

func (it *It[T]) Next() *T {
	if it == nil || it.current == nil {
		return nil
	}
	it.current = nil
	k, v := it.cursor.Next()
	if k == nil {
		return nil
	}
	it.current = it.getObject(k, v)
	return it.current
}
