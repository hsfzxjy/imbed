package lang

import (
	"github.com/hsfzxjy/imbed/db/assetq"
	dbq "github.com/hsfzxjy/imbed/db/query"
)

var (
	fuzzyExprOid = &fuzzyExpr{
		"db.oid@", "a hex string with length <= 64",
		dbq.BytesNeedle, assetq.ByOid}
	fuzzyExprFHash = &fuzzyExpr{
		"db.fhash@", "a hex string with length <= 32",
		dbq.BytesNeedle, assetq.ByFHash}
	fuzzyExprUrl = &fuzzyExpr{
		"db.url@", "a URL string",
		dbq.BytesNeedle, assetq.ByUrl}
	fuzzyExprs = fuzzyExprSet{
		fuzzyExprOid,
		fuzzyExprFHash,
		fuzzyExprUrl,
	}
)

type fuzzyExpr struct {
	directive     string
	expected      string
	needleBuilder func(s string, prefix bool) (dbq.Needle, error)
	queryBuilder  func(dbq.Needle) assetq.Query
}

type fuzzyExprSet []*fuzzyExpr

func (s fuzzyExprSet) parse(c *Context) (assetq.Query, error) {
	for _, expr := range s {
		loader, err := expr.parse(c)
		if err != nil {
			return nil, err
		}
		if loader != nil {
			return loader, nil
		}
	}
	return nil, nil
}

func (f *fuzzyExpr) parse(c *Context) (assetq.Query, error) {
	if !c.parser.Term(f.directive) {
		return nil, nil
	}
	exact := c.parser.Byte('=')
	value, ok := c.parser.String("")
	if !ok {
		return nil, c.parser.Expect(f.expected)
	}
	needle, err := f.needleBuilder(value, !exact)
	if err != nil {
		return nil, err
	}
	return f.queryBuilder(needle), nil
}
