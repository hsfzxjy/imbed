package cmds

import (
	"os"

	"github.com/hsfzxjy/imbed/app"
	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/formatter"
	"github.com/hsfzxjy/imbed/lang"
	"github.com/hsfzxjy/imbed/parser"
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
	langCtx := lang.NewContext(parser.New(command.Args()), app)
	assets, err := langCtx.ParseRun_AddBody()
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
