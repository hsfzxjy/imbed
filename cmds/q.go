package cmds

import (
	"github.com/hsfzxjy/imbed/app"
	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/formatter"
	"github.com/hsfzxjy/imbed/lang"
	"github.com/hsfzxjy/imbed/parser"
	"github.com/hsfzxjy/imbed/util/iter"
	"github.com/hsfzxjy/tipe"
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
	langCtx := lang.NewContext(parser.New(command.Args()), app)
	query, err := langCtx.ParseRun_QueryBody()
	if err != nil {
		return err
	}

	return app.DB().RunR(func(tx *db.Tx) error {
		it, err := query.RunR(tx)
		if err != nil {
			return err
		}
		sortedIt := iter.Sorted(it, (*db.AssetModel).CompareCreated, true)

		return formatter.
			New(asset.FmtFields, c.fmt.Format, !c.fmt.Raw).
			ExecIter(
				app.Stdout(),
				iter.FilterMap(sortedIt, func(m *db.AssetModel) (r tipe.Result[asset.StockAsset]) {
					return tipe.MakeR(asset.FromDBModel(app, m, nil)(tx))
				}))
	})
}
