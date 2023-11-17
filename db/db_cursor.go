package db

import (
	"github.com/hsfzxjy/imbed/util"
	"github.com/hsfzxjy/imbed/util/iter"
	"github.com/hsfzxjy/tipe"
	"go.etcd.io/bbolt"
)

type Cursor struct {
	current util.KV
	cursor  *bbolt.Cursor
}

func NewCursor(cursor *bbolt.Cursor, seekTo []byte) *Cursor {
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

func (c *Cursor) Next() (r tipe.Result[util.KV]) {
	if c == nil || c.current.K == nil {
		return r.FillErr(iter.Stop)
	}
	cur := c.current
	c.current.K, c.current.V = c.cursor.Next()
	return tipe.Ok(cur)
}
