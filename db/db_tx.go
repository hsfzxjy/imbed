package db

import (
	"sync"

	"github.com/hsfzxjy/imbed/db/internal"
	"go.etcd.io/bbolt"
)

//go:generate go run github.com/hsfzxjy/imbed/db/bucketgen

type Tx struct {
	*bbolt.Tx
	buckets [len(bucketNames)]struct {
		sync.Once
		*bbolt.Bucket
	}
	assetMeta internal.AssetMeta
}

func newTx(bbtx *bbolt.Tx) *Tx {
	return &Tx{Tx: bbtx}
}

func (tx *Tx) AssetMetadata() *internal.AssetMeta {
	return &tx.assetMeta
}

func (tx *Tx) runR(f func(*Tx) error) error  { return f(tx) }
func (tx *Tx) runRW(f func(*Tx) error) error { return f(tx) }
func (tx *Tx) DB() Service                   { return Service{tx.Tx.DB()} }
