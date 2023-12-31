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

func (c *assetCache) Lookup(a asset.Asset, transform *step) (asset.StockAsset, error) {
	hash, err := a.Content().GetHash()
	if err != nil {
		return nil, err
	}
	it, err := assetq.ByDependency(hash, transform.Model()).RunR(c.ctx)
	if err != nil {
		return nil, err
	}
	if !it.HasNext() {
		return nil, nil
	}
	model := it.Next()
	if model.IsErr() {
		return nil, model.UnwrapErr()
	}
	return asset.FromDBModel(c.app, model.Unwrap(), a)(nil)
}
