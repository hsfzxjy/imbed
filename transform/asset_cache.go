package transform

import (
	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/db/assetq"
)

type assetCache struct {
	ctx db.Context
	app db.App
}

func (c *assetCache) Lookup(a asset.Asset, transform *transSeq) (asset.StockAsset, error) {
	transSeqHash, err := transform.GetSha256Hash()
	if err != nil {
		return nil, err
	}
	hash, err := a.Content().GetHash()
	if err != nil {
		return nil, err
	}
	it, err := assetq.ByDependency(hash, transSeqHash).RunR(c.ctx)
	if err != nil {
		return nil, err
	}
	if !it.HasNext() {
		return nil, nil
	}
	model := it.Next()
	return asset.FromDBModel(c.app, model, a)(nil)
}
