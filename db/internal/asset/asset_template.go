package asset

import (
	"time"

	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db/internal"
	"github.com/hsfzxjy/imbed/db/internal/bucketnames"
)

type Flag uint8

const (
	HasOrigin Flag = 1 << iota
	SupportsRemove
)

type TransSeq struct {
	Raw            []byte
	Hash           ref.Sha256Hash
	ConfigHashes   []ref.Sha256Hash
	SupportsRemove bool
}

type AssetTemplate struct {
	Origin   *AssetModel
	TransSeq TransSeq
	FID      ref.FID
	Url      string
	ExtData  []byte
}

var oneBytes = []byte{1}

func (t *AssetTemplate) doCreate(h internal.H) (*AssetModel, error) {
	var flag Flag
	if t.Origin != nil {
		flag |= HasOrigin
	}
	if t.TransSeq.SupportsRemove {
		flag |= SupportsRemove
	}
	model := &AssetModel{
		Flag:        flag,
		OriginOID:   t.getOriginOID(),
		Created:     ref.NewTime(time.Now()),
		TransSeqRaw: t.TransSeq.Raw,
		FID:         t.FID,
		Url:         t.Url,
		ExtData:     t.ExtData,
	}
	encoded := encodeAsset(model)
	hash := ref.Sha256HashSum(encoded)
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
			b.UpdateLeaf(ref.AsRaw(ref.NewPair(cfgHash, oid)), oneBytes)
		}
	}

	h.Bucket(bucketnames.INDEX_TIME).
		UpdateLeaf(ref.AsRaw(ref.NewPair(model.Created, oid)), oneBytes)

	if !model.FID.IsZero() {
		h.Bucket(bucketnames.INDEX_FID).
			UpdateLeaf(ref.AsRaw(ref.NewPair(model.FID, oid)), oneBytes)
		h.Bucket(bucketnames.INDEX_FHASH).
			UpdateLeaf(ref.AsRaw(ref.NewPair(model.FID.Hash(), oid)), oneBytes)
	}

	if model.Url != "" {
		h.Bucket(bucketnames.INDEX_URL).
			UpdateLeaf([]byte(model.Url+ref.AsRawString(oid)), oneBytes)
	}

	return model, nil
}
