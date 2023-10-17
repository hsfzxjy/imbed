package cmds

import (
	"github.com/hsfzxjy/imbed/app"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/lang"
	"github.com/hsfzxjy/imbed/parser"
)

type RevCommand struct {
	fmt *fmtOption
}

func (RevCommand) Spec() app.CommandSpec {
	fs, fmt := getFmtFlagSet()
	cmd := RevCommand{fmt}
	return app.CommandSpec{
		Name:    "rev",
		FlagSet: fs,
		Mode:    core.ModeReadonly,
		Runner:  cmd.Run,
	}
}

func (c RevCommand) Run(app *app.App, command app.CommandSpec) error {
	langCtx := lang.NewContext(parser.New(command.Args()), app)
	result, err := langCtx.ParseRun_RevParseBody()
	if err != nil {
		return err
	}
	for _, s := range result {
		println(s)
	}
	return nil
}
