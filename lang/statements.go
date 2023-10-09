package lang

import (
	"fmt"

	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/db/assetq"
	"github.com/hsfzxjy/imbed/db/configq"
)

func (c *Context) ParseRun_AddBody() ([]asset.StockAsset, error) {
	initialLoader, err := c.parseAsset()
	if err != nil {
		return nil, err
	}
	var assets []asset.Asset
	err = c.app.DB().RunR(func(h db.Context) error {
		cfg := configq.NewProvider(h, c.app)
		graph, err := c.parseTransSeq(cfg)
		if err != nil {
			return err
		}
		initialAsset, err := initialLoader.Do(h)
		if err != nil {
			return err
		}
		assets, err = graph.Compute(h, c.app, initialAsset)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	stockAssets, err := asset.SaveAll(c.app.DB(), c.app, assets)
	if err != nil {
		return nil, err
	}
	return stockAssets, nil
}

func (c *Context) ParseRun_QBody() (assetq.Query, error) {
	q, err := fuzzyExprs.parse(c)
	if err != nil {
		return nil, err
	}
	if q != nil {
		return q, nil
	}
	c.parser.Space()
	if !c.parser.EOF() {
		return nil, c.parser.Error(fmt.Errorf("invalid query"))
	}
	return assetq.All(), nil
}
