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
	filename, pos, ok := c.parser.String("")
	if !ok {
		return nil, pos.WrapErrorString("illegal file name")
	}
	return asset.LoadFile(filename, pos), nil
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
