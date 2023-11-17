package db

import (
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

func encodeConfigModel(model *ConfigModel) []byte {
	return fastbuf.Concat(model.SHA.MustSHA().Raw(), model.Data.MustData())
}

func DecodeConfigModel(oid ref.OID, encoded []byte) (*ConfigModel, error) {
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

func DecodeConfigModelSHA(encoded []byte) (ref.Sha256, error) {
	sha, err, _ := ref.FromRaw[ref.Sha256](encoded)
	return sha, err
}

func DecodeConfigModelData(encoded []byte) ([]byte, error) {
	_, err, data := ref.FromRaw[ref.Sha256](encoded)
	return data, err
}
