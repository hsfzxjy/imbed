package db

import (
	"sync"

	"github.com/hsfzxjy/imbed/db/internal"
	"go.etcd.io/bbolt"
)

//go:generate go run github.com/hsfzxjy/imbed/db/bucketgen

type Tx struct {
	*bbolt.Tx
	service *Service
	buckets [len(bucketNames)]struct {
		sync.Once
		*bbolt.Bucket
	}
	assetMeta internal.AssetMeta

	mu          sync.Mutex
	onRollbacks []func()
}

func newTx(service *Service, bbtx *bbolt.Tx) *Tx {
	return &Tx{Tx: bbtx, service: service}
}

func (tx *Tx) AssetMetadata() *internal.AssetMeta {
	return &tx.assetMeta
}

func (tx *Tx) runR(f func(*Tx) error) error  { return f(tx) }
func (tx *Tx) runRW(f func(*Tx) error) error { return f(tx) }
func (tx *Tx) DB() *Service                  { return tx.service }

func (tx *Tx) onRollback(f func()) {
	if f == nil {
		return
	}
	tx.mu.Lock()
	tx.onRollbacks = append(tx.onRollbacks, f)
	tx.mu.Unlock()
}

func (tx *Tx) invokeOnRollback() {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	for _, f := range tx.onRollbacks {
		f()
	}
	tx.onRollbacks = nil
}
