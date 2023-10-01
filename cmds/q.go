package cmds

import (
	"fmt"
	"os"

	"github.com/hsfzxjy/imbed/app"
	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/db/assetq"
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

	p.Space()
	if p.EOF() {
		query = assetq.All()
	} else {
		key, ok := p.Ident()
		if !ok {
			return p.Expect("query key")
		}
		p.Space()
		ok = p.Byte('=')
		if !ok {
			return p.Expect("'='")
		}
		p.Space()
		expr := p.Rest()
		switch key {
		case "url":
			query = assetq.ByUrl(expr)
		default:
			return fmt.Errorf("unknown query '%s=%s'", key, expr)
		}
	}

	return app.DB().RunR(func(ctx db.Context) error {
		it, err := query.RunR(ctx)
		if err != nil {
			return err
		}

		return formatter.
			New(asset.FmtFields, c.fmt.Format, !c.fmt.Raw).
			ExecIter(
				os.Stdout,
				iter.Map(it, func(m *db.AssetModel) asset.Asset {
					return asset.FromDBModel(app, m)
				}))
	})
}