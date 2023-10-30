package asset

import (
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db/internal"
	"github.com/hsfzxjy/imbed/db/internal/bucketnames"
	"github.com/hsfzxjy/imbed/db/internal/helper"
)

type AssetModel struct {
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

func NewFromLeaf(n helper.LeafNode) (*AssetModel, error) {
	a, err := decodeAsset(n.Data())
	if err != nil {
		return nil, err
	}
	a.OID = ref.FromRawString[ref.OID](n.NodeName())
	return a, nil
}

func New(h internal.H, oid []byte) (*AssetModel, error) {
	return NewFromLeaf(h.Bucket(bucketnames.FILES).Leaf(oid))
}

func NewFromKV(k, v []byte) (*AssetModel, error) {
	a, err := decodeAsset(v)
	if err != nil {
		return nil, err
	}
	a.OID = ref.FromRaw[ref.OID](k)
	return a, nil
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
