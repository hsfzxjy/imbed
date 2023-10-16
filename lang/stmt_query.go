package lang

import (
	"fmt"

	"github.com/hsfzxjy/imbed/db/assetq"
)

func (c *Context) ParseRun_QueryBody() (assetq.Query, error) {
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
