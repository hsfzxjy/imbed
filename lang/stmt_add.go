package lang

import (
	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/db/configq"
)

func (c *Context) ParseRun_AddBody() ([]asset.StockAsset, error) {
	initialLoader, err := c.parseAsset()
	if err != nil {
		return nil, err
	}
	var assets []asset.Asset
	err = c.app.DB().RunR(func(tx *db.Tx) error {
		cfg := configq.NewProvider(tx, c.app)
		graph, err := c.parseTransSeq(cfg)
		if err != nil {
			return err
		}
		initialAsset, err := initialLoader.Do(tx)
		if err != nil {
			return err
		}
		assets, err = graph.Run(tx, c.app, initialAsset)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	var stockAssets []asset.StockAsset
	err = c.app.DB().RunRW(func(tx *db.Tx) error {
		stockAssets, err = asset.SaveAll(tx, c.app, assets)
		return err
	})
	if err != nil {
		return nil, err
	}
	return stockAssets, nil
}
