package db

import (
	"github.com/hsfzxjy/imbed/core/ref"
)

type Flag uint8

const (
	HasOrigin Flag = 1 << iota
	SupportsRemove
)

type AssetModel struct {
	Flag
	OID       ref.OID
	OriginOID ref.OID
	SHA       ref.Sha256

	Created      ref.Time
	StepListData StepListData
	FHash        ref.Murmur3
	Basename     string
	Url          string
	ExtData      []byte

	Tags []string
}

func (a *AssetModel) Filename() string {
	return a.FHash.WithName(a.Basename)
}

func (a *AssetModel) CompareCreated(other *AssetModel) int {
	return a.Created.Compare(other.Created)
}

type AssetOpt func(tx *Tx, model *AssetModel) error

func applyOptions(tx *Tx, options []AssetOpt, model *AssetModel) (*AssetModel, error) {
	for _, opt := range options {
		if err := opt(tx, model); err != nil {
			return nil, err
		}
	}
	return model, nil
}

func WithTags() AssetOpt {
	return func(tx *Tx, model *AssetModel) error {
		if model.Tags != nil {
			return nil
		}
		tags, err := TagByOid(model.OID).RunR(tx)
		if err != nil {
			return err
		}
		model.Tags = tags
		return nil
	}
}

func New(tx *Tx, oid []byte, opts ...AssetOpt) (*AssetModel, error) {
	data := tx.FILES().Get(oid)
	a, err := DecodeAsset(data)
	if err != nil {
		return nil, err
	}
	a.OID, err = ref.FromRawExact[ref.OID](oid)
	if err != nil {
		return nil, err
	}
	return applyOptions(tx, opts, a)
}

func NewFromKV(tx *Tx, k, v []byte, opts ...AssetOpt) (*AssetModel, error) {
	a, err := DecodeAsset(v)
	if err != nil {
		return nil, err
	}
	a.OID, err = ref.FromRawExact[ref.OID](k)
	if err != nil {
		return nil, err
	}
	return applyOptions(tx, opts, a)
}
