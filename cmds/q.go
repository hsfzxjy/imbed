package cmds

import (
	"os"

	"github.com/hsfzxjy/imbed/app"
	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/formatter"
	"github.com/hsfzxjy/imbed/lang"
	"github.com/hsfzxjy/imbed/parser"
	"github.com/hsfzxjy/imbed/transform"
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
	langCtx := lang.NewContext(parser.New(command.Args()), app, transform.DefaultRegistry())
	query, err := langCtx.ParseRun_QBody()
	if err != nil {
		return err
	}

	return app.DB().RunR(func(ctx db.Context) error {
		it, err := query.RunR(ctx)
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
