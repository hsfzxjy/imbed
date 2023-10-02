package helper

import (
	"github.com/hsfzxjy/imbed/util"
	"go.etcd.io/bbolt"
)

type Cursor struct {
	first        util.KV
	cursor       *bbolt.Cursor
	stopped      bool
	firstEmitted bool
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
		first:        util.KV{K: k, V: v},
		cursor:       cursor,
		stopped:      false,
		firstEmitted: false,
	}
}

func (c *Cursor) Next() (util.KV, bool) {
	if c == nil || c.stopped {
		return util.KV{}, false
	}
	if !c.firstEmitted {
		c.firstEmitted = true
		return c.first, true
	}
	k, v := c.cursor.Next()
	if k == nil {
		c.stopped = true
		return util.KV{}, false
	}
	return util.KV{K: k, V: v}, true
}
