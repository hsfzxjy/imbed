package assetq

import (
	"bytes"

	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/util"
	"github.com/hsfzxjy/imbed/util/iter"
	"github.com/hsfzxjy/tipe"
	"go.etcd.io/bbolt"
)

type Model = db.AssetModel
type Iterator = iter.Ator[*Model]
type Query = db.Task[Iterator]

type Option = db.AssetOpt

var WithTags = db.WithTags

func simpleQuery(indexFunc func(*db.Tx) *bbolt.Bucket, needle ndl.Needle, options []Option) Query {
	return func(tx *db.Tx) (Iterator, error) {
		bucIndex := indexFunc(tx)
		cursor := db.NewCursor(bucIndex.Cursor(), needle.Bytes())
		it := iter.FilterMap(cursor, func(kv util.KV) (r tipe.Result[*Model]) {
			splitAt := len(kv.K) - ref.OID{}.Sizeof()
			if !needle.Match(kv.K[:splitAt]) {
				return r.FillErr(iter.Stop)
			}
			oid := kv.K[splitAt:]
			data := tx.FILES().Get(oid)
			return tipe.MakeR(db.NewFromKV(tx, oid, data, options...))
		})
		return it, nil
	}
}

func ByFHash(needle ndl.Needle, options ...Option) Query {
	return simpleQuery((*db.Tx).F_FHASH_OID, needle, options)
}

func ByUrl(needle ndl.Needle, options ...Option) Query {
	return simpleQuery((*db.Tx).F_URL_OID, needle, options)
}

func BySHA(needle ndl.Needle, options ...Option) Query {
	return func(tx *db.Tx) (Iterator, error) {
		cursor := db.NewCursor(tx.F_SHA__OID().Cursor(), needle.Bytes())
		return iter.FilterMap(cursor, func(kv util.KV) (r tipe.Result[*Model]) {
			if !needle.Match(kv.K) {
				return r.FillErr(iter.Stop)
			}
			return tipe.MakeR(db.New(tx, kv.V, options...))
		}), nil
	}
}

func ByDependency(fhash ref.Murmur3, transSeq db.StepListTpl, options ...Option) Query {
	return func(tx *db.Tx) (Iterator, error) {
		key := ref.NewPair(fhash, transSeq.MustSHA()).Sum().Raw()
		cursor := db.NewCursor(tx.F_FHASH_TSSHA__OID().Cursor(), key)
		return iter.FilterMap(cursor, func(kv util.KV) (r tipe.Result[*db.AssetModel]) {
			if !bytes.Equal(kv.K, key) {
				return r.FillErr(iter.Stop)
			}
			return tipe.MakeR(db.New(tx, kv.V, options...))
		}), nil
	}
}

func All(options ...Option) Query {
	return func(tx *db.Tx) (Iterator, error) {
		cursor := db.NewCursor(tx.FILES().Cursor(), nil)
		return iter.FilterMap(cursor, func(kv util.KV) (r tipe.Result[*Model]) {
			return tipe.MakeR(db.NewFromKV(tx, kv.K, kv.V, options...))
		}), nil
	}
}
