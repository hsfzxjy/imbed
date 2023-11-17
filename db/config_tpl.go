package db

import (
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/core/ref/lazy"
	"github.com/hsfzxjy/imbed/util/fastbuf"
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
		return newConfigModel(oid, tx.CONFIGS().Get(vOid))
	}
	oid, err := findAvailOID(tx.CONFIGS())
	if err != nil {
		return nil, err
	}
	data, err := c.O.GetData()
	if err != nil {
		return nil, err
	}
	err = tx.CONFIGS().Put(oid.Raw(), fastbuf.Concat(sha.Raw(), data))
	if err != nil {
		return nil, err
	}
	err = tx.C_SHA__OID().Put(sha.Raw(), oid.Raw())
	if err != nil {
		return nil, err
	}
	return &ConfigModel{
		OID:  oid,
		SHA:  SHA(sha),
		Data: data,
	}, nil
}
