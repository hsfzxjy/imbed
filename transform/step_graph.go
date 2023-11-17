package transform

import (
	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/db"
)

type Graph struct {
	reg  *Registry
	root stepSet
}

func (g *Graph) Run(h db.Context, app db.App, initial asset.Asset) ([]asset.Asset, error) {
	var result []asset.Asset
	var cache = assetCache{h, app}
	err := g.root.Run(&cache, initial, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
