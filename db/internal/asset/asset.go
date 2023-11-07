package asset

import (
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db/internal"
	"github.com/hsfzxjy/imbed/db/internal/bucketnames"
	"github.com/hsfzxjy/imbed/db/internal/helper"
	"github.com/hsfzxjy/imbed/db/tagq"
)

type AssetModel struct {
	Flag
	OID       ref.OID
	OriginOID ref.OID

	Created     ref.Time
	TransSeqRaw []byte
	FID         ref.FID
	Url         string
	ExtData     []byte

	Tags []string
}

func (a *AssetModel) CompareCreated(other *AssetModel) int {
	return a.Created.Compare(other.Created)
}

type NewOpt func(h internal.H, model *AssetModel) error

func applyOptions(h internal.H, options []NewOpt, model *AssetModel) (*AssetModel, error) {
	for _, opt := range options {
		if err := opt(h, model); err != nil {
			return nil, err
		}
	}
	return model, nil
}

func WithTags() NewOpt {
	return func(h internal.H, model *AssetModel) error {
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

func NewFromLeaf(h internal.H, n helper.LeafNode, opts ...NewOpt) (*AssetModel, error) {
	a, err := DecodeAsset(n.Data())
	if err != nil {
		return nil, err
	}
	a.OID = ref.FromRawString[ref.OID](n.NodeName())
	return applyOptions(h, opts, a)
}

func New(h internal.H, oid []byte, opts ...NewOpt) (*AssetModel, error) {
	return NewFromLeaf(h, h.Bucket(bucketnames.FILES).Leaf(oid))
}

func NewFromKV(h internal.H, k, v []byte, opts ...NewOpt) (*AssetModel, error) {
	a, err := DecodeAsset(v)
	if err != nil {
		return nil, err
	}
	a.OID = ref.FromRaw[ref.OID](k)
	return applyOptions(h, opts, a)
}

func (template AssetTemplate) Create() internal.Runnable[*AssetModel] {
	return func(h internal.H) (*AssetModel, error) {
		return template.doCreate(h)
	}
}

func (t *AssetTemplate) getOriginOID() (ret ref.OID) {
	if t.Origin != nil {
		ret = t.Origin.OID
	}
	return
}

func (t *AssetTemplate) getOriginFHash() (ret ref.Murmur3Hash) {
	if t.Origin != nil {
		ret = t.Origin.FID.Hash()
	}
	return
}
