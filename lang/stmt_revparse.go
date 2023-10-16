package lang

import (
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/db/assetq"
	"github.com/hsfzxjy/imbed/lang/revparse"
)

func (c *Context) ParseRun_RevParseBody() (string, error) {
	var result string
	query, err := c.parseDBAsset()
	if err != nil {
		return "", err
	}
	err = c.app.DB().RunR(func(ctx db.Context) error {
		asset, err := query(ctx)
		if err != nil {
			return err
		}
		models, err := assetq.Chain(asset.Model(), -1).RunR(ctx)
		if err != nil {
			return err
		}
		parsed, err := revparse.Solve(models, c.registry)
		if err != nil {
			return err
		}
		result = parsed
		return nil
	})
	return result, err
}
