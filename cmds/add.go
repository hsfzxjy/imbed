package cmds

import (
	"fmt"
	"os"

	"github.com/hsfzxjy/imbed/app"
	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/formatter"
	"github.com/hsfzxjy/imbed/transform"
)

type AddCommand struct {
	fmt *fmtOption
}

func (AddCommand) Spec() app.CommandSpec {
	fs, fmt := getFmtFlagSet()
	cmd := AddCommand{fmt}
	return app.CommandSpec{
		Name:    "add",
		FlagSet: fs,
		Mode:    core.ModeReadWrite,
		Runner:  cmd.Run,
	}
}

func (c AddCommand) Run(app *app.App, command app.CommandSpec) error {
	flagSet := command.FlagSet
	if flagSet.NArg() < 1 {
		return fmt.Errorf("no arguments")
	}
	initialAsset := asset.LoadFile(flagSet.Arg(0))
	graph, err := transform.DefaultRegistry().
		Parse(app.Config(), flagSet.Args()[1:])
	if err != nil {
		return err
	}
	assets, err := graph.Compute(app, initialAsset)
	if err != nil {
		return err
	}
	err = asset.SaveAll(app.DB(), app, assets)
	if err != nil {
		return err
	}
	fmter := formatter.New(asset.FmtFields, c.fmt.Format, !c.fmt.Raw)
	err = fmter.Exec(os.Stdout, assets)
	if err != nil {
		return err
	}
	return nil
}
