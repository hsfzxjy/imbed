package internal

import (
	"path"
	"time"

	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/db/internal/bucketnames"
	"github.com/hsfzxjy/imbed/db/internal/helper"
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
			h := helper.New(tx)
			h.BucketOrCreate(bucketnames.FILES)
			h.BucketOrCreate(bucketnames.INDEX_FID)
			h.BucketOrCreate(bucketnames.INDEX_FHASH)
			h.BucketOrCreate(bucketnames.INDEX_URL)
			h.BucketOrCreate(bucketnames.INDEX_CONFIG_HASHES)
			h.BucketOrCreate(bucketnames.INDEX_TRANSSEQ)
			h.BucketOrCreate(bucketnames.INDEX_TIME)
			h.BucketOrCreate(bucketnames.CONFIGS)
			return h.Err()
		})

		if err != nil {
			return Service{}, err
		}
	}

	return Service{db}, nil
}

func (s Service) RunR(f func(Context) error) error {
	return s.runR(func(h H) error { return f(h) })
}
func (s Service) RunRW(f func(h H) error) error {
	return s.runRW(func(h H) error { return f(h) })
}
func (s Service) DB() Service { return s }

func (s Service) runR(f func(H) error) error {
	return s.db.View(func(tx *bolt.Tx) error {
		h := H{helper.New(tx)}
		err := f(h)
		if err != nil {
			return err
		}
		return h.Err()
	})
}

func (s Service) runRW(f func(H) error) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		h := H{helper.New(tx)}
		err := f(h)
		if err != nil {
			return err
		}
		return h.Err()
	})
}

func (s Service) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}
