package db

import (
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/core/ref/lazy"
)

type ConfigTpl struct {
	O lazy.MustDataSHAObject
}

func (c ConfigTpl) create(tx *Tx) (*ConfigModel, error) {
	if model, ok := c.O.(*ConfigModel); ok {
		return model, nil
	}

	sha, err := c.O.GetSHA()
	if err != nil {
		return nil, err
	}
	if vOid := tx.C_SHA__OID().Get(sha.Raw()); vOid != nil {
		oid, err := ref.FromRawExact[ref.OID](vOid)
		if err != nil {
			return nil, err
		}
		return DecodeConfigModel(oid, tx.CONFIGS().Get(vOid))
	}
	data, err := c.O.GetData()
	if err != nil {
		return nil, err
	}
	oid, err := findAvailOID(tx.CONFIGS())
	if err != nil {
		return nil, err
	}
	model := &ConfigModel{
		OID:  oid,
		SHA:  SHA(sha),
		Data: data,
	}
	err = tx.CONFIGS().Put(oid.Raw(), encodeConfigModel(model))
	if err != nil {
		return nil, err
	}
	err = tx.C_SHA__OID().Put(sha.Raw(), oid.Raw())
	if err != nil {
		return nil, err
	}
	return model, nil
}
