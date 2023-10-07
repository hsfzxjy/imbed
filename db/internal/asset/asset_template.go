package asset

import (
	"time"

	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db/internal"
	"github.com/hsfzxjy/imbed/db/internal/bucketnames"
)

type TransSeq struct {
	Raw          []byte
	Hash         ref.Sha256Hash
	ConfigHashes []ref.Sha256Hash
}

type AssetTemplate struct {
	Origin   *AssetModel
	TransSeq TransSeq
	FID      ref.FID
	Url      string
	ExtData  []byte
}

func (t *AssetTemplate) doCreate(h internal.H) (*AssetModel, error) {
	model := &AssetModel{
		OriginOID:   t.getOriginOID(),
		Created:     ref.NewTime(time.Now()),
		TransSeqRaw: t.TransSeq.Raw,
		FID:         t.FID,
		Url:         t.Url,
		ExtData:     t.ExtData,
	}
	hash, encoded, err := ref.HashEncodeFunc2(encodeAsset, model)
	if err != nil {
		return nil, err
	}
	oid := ref.OIDFromSha256Hash(hash)
	model.OID = oid

	h.Bucket(bucketnames.FILES).
		UpdateLeaf(ref.AsRaw(oid), encoded)

	if !model.OriginOID.IsZero() {
		// Pair(orig_fhash, transseq_hash) -> []oid
		pair := ref.NewPair(t.getOriginFHash(), t.TransSeq.Hash)
		h.Bucket(bucketnames.INDEX_TRANSSEQ).
			UpdateLeaf(ref.AsRaw(pair), ref.AsRaw(oid))

		b := h.Bucket(bucketnames.INDEX_CONFIG_HASHES)
		for _, cfgHash := range t.TransSeq.ConfigHashes {
			b.BucketOrCreate(ref.AsRaw(cfgHash)).
				UpdateLeaf(ref.AsRaw(oid), []byte{1})
		}
	}

	h.Bucket(bucketnames.INDEX_TIME).
		BucketOrCreate(ref.AsRaw(model.Created)).
		UpdateLeaf(ref.AsRaw(oid), []byte{1})

	if !model.FID.IsZero() {
		h.Bucket(bucketnames.INDEX_FID).
			BucketOrCreate(ref.AsRaw(model.FID)).
			UpdateLeaf(ref.AsRaw(oid), []byte{1})

		h.Bucket(bucketnames.INDEX_FHASH).
			BucketOrCreate(ref.AsRaw(model.FID.Hash())).
			UpdateLeaf(ref.AsRaw(oid), []byte{1})
	}

	if model.Url != "" {
		h.Bucket(bucketnames.INDEX_URL).
			BucketOrCreate([]byte(model.Url)).
			UpdateLeaf(ref.AsRaw(oid), []byte{1})
	}

	return model, nil
}
