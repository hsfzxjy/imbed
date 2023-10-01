package assetq

import (
	"bytes"

	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db/internal"
	"github.com/hsfzxjy/imbed/db/internal/asset"
	"github.com/hsfzxjy/imbed/db/internal/bucketnames"
	"github.com/hsfzxjy/imbed/db/internal/helper"
	"github.com/hsfzxjy/imbed/db/internal/iterator"
	"github.com/hsfzxjy/imbed/util"
)

type Iterator = iterator.It[asset.AssetModel]
type Query = internal.Runnable[*Iterator]

func simpleQuery(getIndex func(internal.H) helper.BucketNode) Query {
	return iterator.New[asset.AssetModel](func(h internal.H) (*iterator.Builder[asset.AssetModel], error) {
		index := getIndex(h)
		cursor, err := index.Cursor()
		if err != nil {
			return nil, err
		}
		return &iterator.Builder[asset.AssetModel]{
			Cursor: cursor,
			SeekTo: nil,
			GetObject: func(k []byte, v []byte) *asset.AssetModel {
				return util.IgnoreErr(asset.New(h, k))
			},
		}, nil
	})

}

func ByFID(fid ref.FID) Query {
	return simpleQuery(func(h internal.H) helper.BucketNode {
		return h.Bucket(bucketnames.INDEX_FID).Bucket(ref.AsRaw(fid))
	})
}

func ByFHash(fhash ref.Murmur3Hash) Query {
	return simpleQuery(func(h internal.H) helper.BucketNode {
		return h.Bucket(bucketnames.INDEX_FHASH).Bucket(ref.AsRaw(fhash))
	})
}

func ByUrl(url string) Query {
	return simpleQuery(func(h internal.H) helper.BucketNode {
		return h.Bucket(bucketnames.INDEX_URL).Bucket([]byte(url))
	})
}

func ByDependency(fhash ref.Murmur3Hash, transSeqHash ref.Sha256Hash) Query {
	return iterator.New[asset.AssetModel](func(h internal.H) (*iterator.Builder[asset.AssetModel], error) {
		cursor, err := h.Bucket(bucketnames.INDEX_TRANSSEQ).Cursor()
		if err != nil {
			return nil, err
		}
		pairBytes := ref.AsRaw(ref.NewPair(fhash, transSeqHash))
		return &iterator.Builder[asset.AssetModel]{
			Cursor: cursor,
			SeekTo: pairBytes,
			GetObject: func(k []byte, v []byte) *asset.AssetModel {
				if !bytes.Equal(k, pairBytes) {
					return nil
				}
				return util.IgnoreErr(asset.New(h, v))
			},
		}, nil
	})
}

func All() Query {
	return iterator.New[asset.AssetModel](func(h internal.H) (*iterator.Builder[asset.AssetModel], error) {
		cursor, err := h.Bucket(bucketnames.FILES).Cursor()
		if err != nil {
			return nil, err
		}
		return &iterator.Builder[asset.AssetModel]{
			Cursor: cursor,
			SeekTo: nil,
			GetObject: func(k, v []byte) *asset.AssetModel {
				a, _ := asset.NewFromKV(k, v)
				return a
			},
		}, nil
	})
}
