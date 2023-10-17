package lang

import (
	"github.com/hsfzxjy/imbed/asset"
)

func (c *Context) parseAsset() (asset.Loader, error) {
	q, err := fuzzyExprs.parse(c)
	if err != nil {
		return nil, err
	}
	if q != nil {
		return asset.FromQ(c.app, q), nil
	}
	c.parser.Term("file@")
	filename, ok := c.parser.String("")
	if !ok {
		return nil, c.parser.ErrorString("illegal file name")
	}
	return asset.LoadFile(filename), nil
}

func (c *Context) parseDBAsset() (asset.StockLoader, error) {
	q, err := fuzzyExprs.parse(c)
	if err != nil {
		return nil, err
	}
	if q != nil {
		return asset.FromQ(c.app, q), nil
	}
	return nil, c.parser.ErrorString("expect asset query")
}
