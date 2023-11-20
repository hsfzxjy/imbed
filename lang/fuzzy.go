package lang

import (
	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/core/pos"
	"github.com/hsfzxjy/imbed/db/assetq"
)

var (
	fuzzyExprOid = &fuzzyExpr{
		"sha@", "a hex string with length <= 64",
		ndl.Hex, assetq.BySHA}
	fuzzyExprFHash = &fuzzyExpr{
		"fhash@", "a hex string with length <= 32",
		ndl.Hex, assetq.ByFHash}
	fuzzyExprUrl = &fuzzyExpr{
		"url@", "a URL string",
		ndl.Raw, assetq.ByUrl}
	fuzzyExprs = fuzzyExprSet{
		fuzzyExprOid,
		fuzzyExprFHash,
		fuzzyExprUrl,
	}
)

type fuzzyExpr struct {
	directive     string
	expected      string
	needleBuilder func(s string, prefix bool, pos pos.P) (ndl.Needle, error)
	queryBuilder  func(ndl.Needle, ...assetq.Option) assetq.Query
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
	if _, ok := c.parser.Term(f.directive); !ok {
		return nil, nil
	}
	_, exact := c.parser.Byte('=')
	value, pos, ok := c.parser.String("")
	if !ok {
		return nil, pos.WrapError(c.parser.Error(nil))
	}
	needle, err := f.needleBuilder(value, !exact, pos)
	if err != nil {
		return nil, pos.WrapError(err)
	}
	return f.queryBuilder(needle, assetq.WithTags()), nil
}
