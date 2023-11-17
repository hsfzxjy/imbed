package db

import (
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/core/ref/lazy"
)

type SHA = lazy.ConstSHA
type Data = lazy.ConstData

type ConfigModel struct {
	OID ref.OID
	SHA
	Data
}

func newConfigModel(oid ref.OID, encoded []byte) (*ConfigModel, error) {
	sha, err, rest := ref.FromRaw[ref.Sha256](encoded)
	if err != nil {
		return nil, err
	}
	return &ConfigModel{
		OID:  oid,
		SHA:  SHA(sha),
		Data: rest,
	}, nil
}
