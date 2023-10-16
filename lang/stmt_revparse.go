package lang

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/db/assetq"
	"github.com/hsfzxjy/imbed/schema"
	schemavisitor "github.com/hsfzxjy/imbed/schema/visitor"
	"github.com/hsfzxjy/imbed/transform"
	"github.com/hsfzxjy/imbed/util/iter"
)

func (c *Context) ParseRun_RevParseBody() (string, error) {
	var builder strings.Builder
	_ = builder
	err := c.app.DB().RunR(func(ctx db.Context) error {
		query, err := c.parseDBAsset()
		if err != nil {
			return err
		}
		asset, err := query(ctx)
		if err != nil {
			return err
		}
		model := asset.Model()
		var models = []*db.AssetModel{model}
		for !model.OriginOID.IsZero() {
			it, err := assetq.ByOid(core.StringFull(ref.AsRawString(model.OriginOID))).RunR(ctx)
			if err != nil {
				return err
			}
			origModel, err := iter.One(it)
			if err != nil {
				return err
			}
			models = append(models, origModel)
			model = origModel
		}
		var builders []transform.Builder
		for i := len(models) - 1; i >= 0; i-- {
			model := models[i]
			vs, err := c.registry.DecodeParams(model.TransSeqRaw)
			if err != nil {
				return err
			}
			builders = append(builders, vs...)
		}
		var v paramsVisitor
		for i, b := range builders {
			if i != 0 {
				v.b.WriteString(", ")
			}
			v.b.WriteString(b.Metadata().Name())
			v.b.WriteByte('@')
			v.b.WriteString(b.ConfigHash().FmtHumanize())
			err := b.VisitParams(&v)
			if err != nil {
				return err
			}
		}
		fmt.Printf("%#v\n", builders)
		println(v.GetString())
		return nil
	})
	return "", err
}

type paramsVisitor struct {
	b strings.Builder
	schemavisitor.Void
}

func (v *paramsVisitor) GetString() string { return v.b.String() }

func (v *paramsVisitor) VisitStructFieldBegin(name string) error {
	v.b.WriteString(":")
	v.b.WriteString(name)
	v.b.WriteByte('=')
	return nil
}

func (v *paramsVisitor) VisitBool(x bool) error {
	var s string
	if x {
		s = "true"
	} else {
		s = "false"
	}
	v.b.WriteString(s)
	return nil
}

func (v *paramsVisitor) VisitInt64(x int64) error {
	s := strconv.FormatInt(x, 10)
	v.b.WriteString(s)
	return nil
}

func (*paramsVisitor) VisitRat(x *big.Rat) error {
	panic("unimplemented")
}

func (v *paramsVisitor) VisitString(x string) error {
	v.b.WriteString(x)
	return nil
}

func (v *paramsVisitor) VisitStruct(size int) (sv schema.StructVisitor, elem schema.Visitor, err error) {
	return v, v, nil
}
