package lang

import (
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/db/assetq"
	"github.com/hsfzxjy/imbed/lang/revparse"
)

func (c *Context) ParseRun_RevParseBody() ([]string, error) {
	var results []string
	query, err := fuzzyExprs.parse(c)
	if err != nil {
		return nil, err
	}
	err = c.app.DB().RunR(func(ctx db.Context) error {
		it, err := assetq.Chains(query, -1).RunR(ctx)
		if err != nil {
			return err
		}
		for it.HasNext() {
			parsed, err := revparse.Solve(it.Next(), c.registry)
			if err != nil {
				return err
			}
			results = append(results, parsed)
		}
		return nil
	})
	return results, err
}
