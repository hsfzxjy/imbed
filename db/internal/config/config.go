package config

import (
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db/internal"
	"github.com/hsfzxjy/imbed/db/internal/bucketnames"
)

type TransSeq struct {
	Raw          []byte
	Hash         ref.Sha256Hash
	ConfigHashes []ref.Sha256Hash
}

type ConfigModel struct {
	Raw  []byte
	Hash ref.Sha256Hash
}

func (template ConfigModel) Create() internal.Runnable[*ConfigModel] {
	return internal.R[*ConfigModel](func(h internal.H) (*ConfigModel, error) {
		ret := &ConfigModel{}
		*ret = template
		h.Bucket(bucketnames.CONFIGS).
			UpdateLeaf(ref.AsRaw(ret.Hash), ret.Raw)
		return ret, nil
	})
}
