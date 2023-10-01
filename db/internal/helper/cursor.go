package helper

import (
	"github.com/hsfzxjy/imbed/util"
	"go.etcd.io/bbolt"
)

type Cursor struct {
	n       BucketNode
	cursor  *bbolt.Cursor
	current util.KV
}

func (c *Cursor) Exhausted() bool {
	return c.current.K == nil
}

func (c *Cursor) Current() util.KV {
	return c.current
}

func (c *Cursor) Next() {
	k, v := c.cursor.Next()
	c.current = util.KV{K: k, V: v}
}
