package assetq

import (
	"bytes"

	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db/internal"
	"github.com/hsfzxjy/imbed/db/internal/asset"
	"github.com/hsfzxjy/imbed/db/internal/bucketnames"
	"github.com/hsfzxjy/imbed/util"
	"github.com/hsfzxjy/imbed/util/iter"
	"github.com/hsfzxjy/tipe"
)

type Model = asset.AssetModel
type Iterator = iter.Ator[*Model]
type Query = internal.Runnable[Iterator]

type Option = asset.NewOpt

var WithTags = asset.WithTags

func simpleQuery(indexName []byte, needle ndl.Needle, options []Option) Query {
	return func(h internal.H) (Iterator, error) {
		index := h.Bucket(indexName)
		cursor, err := index.Cursor(needle.Bytes())
		if err != nil {
			return nil, err
		}
		it := iter.FilterMap(cursor, func(kv util.KV) (r tipe.Result[*Model]) {
			splitAt := len(kv.K) - ref.OID_LEN
			if !needle.Match(kv.K[:splitAt]) {
				return r.FillErr(iter.Stop)
			}
			return tipe.MakeR(asset.NewFromKV(h, kv.K, kv.V, options...))
		})
		return it, nil
	}
}

func ByFID(needle ndl.Needle, options ...Option) Query {
	return simpleQuery(bucketnames.INDEX_FID, needle, options)
}

func ByFHash(needle ndl.Needle, options ...Option) Query {
	return simpleQuery(bucketnames.INDEX_FHASH, needle, options)
}

func ByUrl(needle ndl.Needle, options ...Option) Query {
	return simpleQuery(bucketnames.INDEX_URL, needle, options)
}

func ByOid(needle ndl.Needle, options ...Option) Query {
	return func(h internal.H) (Iterator, error) {
		cursor, err := h.Bucket(bucketnames.FILES).Cursor(needle.Bytes())
		if err != nil {
			return nil, err
		}
		return iter.FilterMap(cursor, func(kv util.KV) (r tipe.Result[*Model]) {
			if !needle.Match(kv.K) {
				return r.FillErr(iter.Stop)
			}
			return tipe.MakeR(asset.NewFromKV(h, kv.K, kv.V, options...))
		}), nil
	}
}

func ByDependency(fhash ref.Murmur3Hash, transSeqHash ref.Sha256Hash, options ...Option) Query {
	return func(h internal.H) (Iterator, error) {
		pairBytes := ref.AsRaw(ref.NewPair(fhash, transSeqHash))
		cursor, err := h.Bucket(bucketnames.INDEX_TRANSSEQ).
			Cursor(pairBytes)
		if err != nil {
			return nil, err
		}
		return iter.FilterMap(cursor, func(kv util.KV) (r tipe.Result[*asset.AssetModel]) {
			if !bytes.Equal(kv.K, pairBytes) {
				return r.FillErr(iter.Stop)
			}
			return tipe.MakeR(asset.New(h, kv.V, options...))
		}), nil
	}
}

func All(options ...Option) Query {
	return func(h internal.H) (Iterator, error) {
		cursor, err := h.Bucket(bucketnames.FILES).Cursor(nil)
		if err != nil {
			return nil, err
		}
		return iter.FilterMap(cursor, func(kv util.KV) (r tipe.Result[*Model]) {
			return tipe.MakeR(asset.NewFromKV(h, kv.K, kv.V, options...))
		}), nil
	}
}
