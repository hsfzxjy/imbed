package cmds

import (
	"fmt"
	"os"

	"github.com/hsfzxjy/imbed/app"
	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/db/assetq"
	dbq "github.com/hsfzxjy/imbed/db/query"
	"github.com/hsfzxjy/imbed/formatter"
	"github.com/hsfzxjy/imbed/parser"
	"github.com/hsfzxjy/imbed/util/iter"
)

type QCommand struct {
	fmt *fmtOption
}

func (QCommand) Spec() app.CommandSpec {
	fs, fmt := getFmtFlagSet()
	cmd := QCommand{fmt}
	return app.CommandSpec{
		Name:    "q",
		FlagSet: fs,
		Mode:    core.ModeReadonly,
		Runner:  cmd.Run,
	}
}

func (c QCommand) Run(app *app.App, command app.CommandSpec) error {
	flagSet := command.FlagSet
	p := parser.New(flagSet.Args())

	var query assetq.Query
	var usePrefix bool

	p.Space()
	if p.EOF() {
		query = assetq.All()
	} else {
		key, ok := p.Ident()
		if !ok {
			return p.Expect("query key")
		}
		p.Space()
		if p.Term("^=") {
			usePrefix = true
		} else if p.Byte('=') {
		} else {
			if !ok {
				return p.Expect("'='")
			}
		}
		p.Space()
		expr := p.Rest()
		switch key {
		case "url":
			query = assetq.ByUrl(dbq.StringNeedle(expr, usePrefix))
		case "fhash":
			needle, err := dbq.BytesNeedle(expr, usePrefix)
			if err != nil {
				return err
			}
			query = assetq.ByFHash(needle)
		default:
			return fmt.Errorf("unknown query '%s=%s'", key, expr)
		}
	}

	return app.DB().RunR(func(ctx db.Context) error {
		it, err := query.
			TransformR(assetq.SortByNewest).
			RunR(ctx)
		if err != nil {
			return err
		}
		sortedIt := iter.Sorted(it, (*db.AssetModel).CompareCreated, true)

		return formatter.
			New(asset.FmtFields, c.fmt.Format, !c.fmt.Raw).
			ExecIter(
				os.Stdout,
				iter.FilterMap(sortedIt, func(m *db.AssetModel) (asset.StockAsset, bool) {
					a, _ := asset.FromDBModel(app, m)(ctx)
					return a, true
				}))
	})
}
