package db

import (
	"time"

	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db/internal"
	"go.etcd.io/bbolt"
)

type AssetTpl struct {
	Origin   *AssetModel
	TransSeq StepListTpl
	FHash    ref.Murmur3
	Basename string
	Url      string
	ExtData  []byte
}

func (t *AssetTpl) computeStepList(tx *Tx) ([]*ConfigModel, []byte, error) {
	if t.TransSeq == nil {
		return nil, nil, nil
	}
	return t.TransSeq.create(tx)
}

func (t *AssetTpl) doCreate(tx *Tx) (*AssetModel, error) {
	var flag Flag
	if t.Origin != nil {
		flag |= HasOrigin
	}

	cfgModels, transSeqRaw, err := t.computeStepList(tx)
	if err != nil {
		return nil, err
	}

	ts := t.TransSeq
	if ts != nil && ts.SupportsRemove() {
		flag |= SupportsRemove
	}

	model := &AssetModel{
		Flag:         flag,
		OriginOID:    t.getOriginOID(),
		Created:      ref.NewTime(time.Now()),
		StepListData: transSeqRaw,
		FHash:        t.FHash,
		Basename:     t.Basename,
		Url:          t.Url,
		ExtData:      t.ExtData,
	}
	if t.Origin != nil {
		model.DepSHA = ref.NewPair(t.getOriginFHash(), ts.MustSHA()).Sum()
	}
	encoded := encodeAsset(model)

	var oid ref.OID
	{
		x, err := internal.AssetMetaNextOID(&tx.assetMeta)
		if err != nil {
			return nil, err
		}
		oid = ref.NewOID(x)
	}

	if err = tx.FILES().Put(oid.Raw(), encoded); err != nil {
		return nil, err
	}
	model.OID = oid

	updateIndex := func(buc *bbolt.Bucket, key, key2, value []byte) {
		if err != nil {
			return
		}
		if key2 != nil {
			key = append(key, key2...)
		}
		err = buc.Put(key, value)
	}

	updateIndex(tx.F_SHA__OID(), model.SHA.Raw(), nil, oid.Raw())

	if !model.OriginOID.IsZero() {
		updateIndex(
			tx.F_FHASH_TSSHA__OID(),
			model.DepSHA.Raw(),
			nil,
			oid.Raw(),
		)

		for _, cfgModel := range cfgModels {
			err = tx.T_COID_FOID().Put(
				ref.NewPair(cfgModel.OID, oid).Raw(),
				oneBytes,
			)
			if err != nil {
				return nil, err
			}
		}

	}

	if !model.FHash.IsZero() {
		updateIndex(tx.F_FHASH_OID(), model.FHash.Raw(), oid.Raw(), oneBytes)
	}

	if model.Basename != "" {
		updateIndex(tx.F_BASENAME_OID(), []byte(model.Basename), oid.Raw(), oneBytes)
	}

	if model.Url != "" {
		updateIndex(tx.F_URL_OID(), []byte(model.Url), oid.Raw(), oneBytes)
	}

	if err != nil {
		return nil, err
	}

	return model, nil
}

func (template AssetTpl) Create() Task[*AssetModel] {
	return func(tx *Tx) (*AssetModel, error) {
		return template.doCreate(tx)
	}
}

func (t *AssetTpl) getOriginOID() (ret ref.OID) {
	if t.Origin != nil {
		ret = t.Origin.OID
	}
	return
}

func (t *AssetTpl) getOriginFHash() (ret ref.Murmur3) {
	if t.Origin != nil {
		ret = t.Origin.FHash
	}
	return
}
