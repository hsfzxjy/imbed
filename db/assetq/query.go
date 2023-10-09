package assetq

import (
	"bytes"

	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db/internal"
	"github.com/hsfzxjy/imbed/db/internal/asset"
	"github.com/hsfzxjy/imbed/db/internal/bucketnames"
	"github.com/hsfzxjy/imbed/db/internal/helper"
	"github.com/hsfzxjy/imbed/util"
	"github.com/hsfzxjy/imbed/util/iter"
)

type Iterator = core.Iterator[*asset.AssetModel]
type Query = internal.Runnable[Iterator]

func simpleQuery(indexName []byte, needle core.Needle) Query {
	return func(h internal.H) (Iterator, error) {
		index := h.Bucket(indexName)
		cursor, err := index.Cursor(needle.Bytes())
		if err != nil {
			return nil, err
		}
		it := iter.FilterMap(cursor, func(kv util.KV) (*helper.Cursor, bool) {
			if !needle.Match(kv.K) {
				return nil, false
			}
			cursor, err := index.Bucket(kv.K).Cursor(nil)
			if err != nil {
				return nil, false
			}
			return cursor, true
		})
		it2 := iter.FlatFilterMap(it, func(kv util.KV) (*asset.AssetModel, bool) {
			a, err := asset.New(h, kv.K)
			if err != nil {
				return nil, false
			}
			return a, true
		})
		return it2, nil
	}
}

func ByFID(needle core.Needle) Query {
	return simpleQuery(bucketnames.INDEX_FID, needle)
}

func ByFHash(needle core.Needle) Query {
	return simpleQuery(bucketnames.INDEX_FHASH, needle)
}

func ByUrl(needle core.Needle) Query {
	return simpleQuery(bucketnames.INDEX_URL, needle)
}

func ByOid(needle core.Needle) Query {
	return func(h internal.H) (Iterator, error) {
		cursor, err := h.Bucket(bucketnames.FILES).Cursor(needle.Bytes())
		if err != nil {
			return nil, err
		}
		return iter.FilterMap(cursor, func(kv util.KV) (*asset.AssetModel, bool) {
			if !needle.Match(kv.K) {
				return nil, false
			}
			a, err := asset.NewFromKV(kv.K, kv.V)
			if err != nil {
				return nil, false
			}
			return a, true
		}), nil
	}
}

func ByDependency(fhash ref.Murmur3Hash, transSeqHash ref.Sha256Hash) Query {
	return func(h internal.H) (Iterator, error) {
		pairBytes := ref.AsRaw(ref.NewPair(fhash, transSeqHash))
		cursor, err := h.Bucket(bucketnames.INDEX_TRANSSEQ).
			Cursor(pairBytes)
		if err != nil {
			return nil, err
		}
		return iter.FilterMap(cursor, func(kv util.KV) (*asset.AssetModel, bool) {
			if !bytes.Equal(kv.K, pairBytes) {
				return nil, false
			}
			asset, err := asset.New(h, kv.V)
			if err != nil {
				return nil, false
			}
			return asset, true
		}), nil
	}
}

func All() Query {
	return func(h internal.H) (Iterator, error) {
		cursor, err := h.Bucket(bucketnames.FILES).Cursor(nil)
		if err != nil {
			return nil, err
		}
		return iter.FilterMap(cursor, func(kv util.KV) (*asset.AssetModel, bool) {
			a, err := asset.NewFromKV(kv.K, kv.V)
			if err != nil {
				return nil, false
			}
			return a, true
		}), nil
	}
}
