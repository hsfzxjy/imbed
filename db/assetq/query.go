package assetq

import (
	"bytes"

	"github.com/hsfzxjy/imbed/core"
	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db/internal"
	"github.com/hsfzxjy/imbed/db/internal/asset"
	"github.com/hsfzxjy/imbed/db/internal/bucketnames"
	"github.com/hsfzxjy/imbed/db/tagq"
	"github.com/hsfzxjy/imbed/util"
	"github.com/hsfzxjy/imbed/util/iter"
)

type Iterator = core.Iterator[*asset.AssetModel]
type Query = internal.Runnable[Iterator]

type Option func(h internal.H, model *asset.AssetModel) error

func applyOptions(h internal.H, options []Option, model *asset.AssetModel) error {
	for _, opt := range options {
		if err := opt(h, model); err != nil {
			return err
		}
	}
	return nil
}

func WithTags() Option {
	return func(h internal.H, model *asset.AssetModel) error {
		if model.Tags != nil {
			return nil
		}
		tags, err := tagq.ByOid(model.OID).RunR(h)
		if err != nil {
			return err
		}
		model.Tags = tags
		return nil
	}
}

func simpleQuery(indexName []byte, needle ndl.Needle, options []Option) Query {
	return func(h internal.H) (Iterator, error) {
		index := h.Bucket(indexName)
		cursor, err := index.Cursor(needle.Bytes())
		if err != nil {
			return nil, err
		}
		it := iter.FilterMap(cursor, func(kv util.KV) (*asset.AssetModel, bool) {
			splitAt := len(kv.K) - ref.OID_LEN
			if !needle.Match(kv.K[:splitAt]) {
				return nil, false
			}
			a, err := asset.New(h, kv.K[splitAt:])
			if err != nil {
				return nil, false
			}
			if applyOptions(h, options, a) != nil {
				return nil, false
			}
			return a, true
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
		return iter.FilterMap(cursor, func(kv util.KV) (*asset.AssetModel, bool) {
			if !needle.Match(kv.K) {
				return nil, false
			}
			a, err := asset.NewFromKV(kv.K, kv.V)
			if err != nil {
				return nil, false
			}
			if applyOptions(h, options, a) != nil {
				return nil, false
			}
			return a, true
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
		return iter.FilterMap(cursor, func(kv util.KV) (*asset.AssetModel, bool) {
			if !bytes.Equal(kv.K, pairBytes) {
				return nil, false
			}
			asset, err := asset.New(h, kv.V)
			if err != nil {
				return nil, false
			}
			if applyOptions(h, options, asset) != nil {
				return nil, false
			}
			return asset, true
		}), nil
	}
}

func All(options ...Option) Query {
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
			if applyOptions(h, options, a) != nil {
				return nil, false
			}
			return a, true
		}), nil
	}
}
