package lang

import (
	"fmt"

	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/db/assetq"
)

func (c *Context) ParseRun_AddBody() ([]asset.StockAsset, error) {
	initialLoader, err := c.parseAsset()
	if err != nil {
		return nil, err
	}
	graph, err := c.parseTransSeq()
	if err != nil {
		return nil, err
	}
	initialAsset, err := initialLoader.Do(c.app.DB())
	if err != nil {
		return nil, err
	}
	assets, err := graph.Compute(c.app, initialAsset)
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
