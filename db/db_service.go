package db

import (
	"path"
	"time"

	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/db/internal"
	bolt "go.etcd.io/bbolt"
)

type App interface {
	core.App
	DB() Service
}

type Service struct {
	db *bolt.DB
}

const _DB_FILENAME = "db"

func Open(app core.App) (Service, error) {
	dbPath := path.Join(app.DBDir(), _DB_FILENAME)
	db, err := bolt.Open(dbPath, 0o600, &bolt.Options{
		Timeout:  1 * time.Second,
		ReadOnly: app.Mode() == core.ModeReadonly,
	})
	if err != nil {
		return Service{}, err
	}

	if app.Mode() == core.ModeReadWrite {
		err = db.Update(func(tx *bolt.Tx) error {
			return newTx(tx).createAllBuckets()
		})

		if err != nil {
			return Service{}, err
		}
	}

	return Service{db}, nil
}

func (s Service) RunR(f func(*Tx) error) error {
	return s.runR(f)
}
func (s Service) RunRW(f func(h *Tx) error) error {
	return s.runRW(f)
}
func (s Service) DB() Service { return s }

func (s Service) runR(f func(*Tx) error) error {
	return s.db.View(func(bbtx *bolt.Tx) error {
		tx := newTx(bbtx)
		if err := internal.DecodeAssetMeta(&tx.assetMeta, tx.f_meta()); err != nil {
			return err
		}
		return f(tx)
	})
}

func (s Service) runRW(f func(*Tx) error) error {
	return s.db.Update(func(bbtx *bolt.Tx) error {
		tx := newTx(bbtx)
		if err := internal.DecodeAssetMeta(&tx.assetMeta, tx.f_meta()); err != nil {
			return err
		}
		if err := f(tx); err != nil {
			return err
		}
		return internal.WriteAssetMeta(&tx.assetMeta, tx.f_meta())
	})
}

func (s Service) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}
