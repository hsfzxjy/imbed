package helper

import (
	"github.com/hsfzxjy/imbed/util"
	"go.etcd.io/bbolt"
)

type Cursor struct {
	current util.KV
	cursor  *bbolt.Cursor
}

func newCursor(cursor *bbolt.Cursor, seekTo []byte) *Cursor {
	var k, v []byte
	if seekTo != nil {
		k, v = cursor.Seek(seekTo)
	} else {
		k, v = cursor.First()
	}
	if k == nil {
		return nil
	}
	return &Cursor{
		current: util.KV{K: k, V: v},
		cursor:  cursor,
	}
}

func (c *Cursor) HasNext() bool {
	return c != nil && c.current.K != nil
}

func (c *Cursor) Next() (result util.KV) {
	if !c.HasNext() {
		return
	}
	result = c.current
	c.current.K, c.current.V = c.cursor.Next()
	return result
}
