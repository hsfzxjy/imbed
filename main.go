package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/hsfzxjy/imbed/app"
	"github.com/hsfzxjy/imbed/cmds"
	"github.com/hsfzxjy/imbed/contrib"
	"github.com/hsfzxjy/imbed/core/pos"
	"github.com/hsfzxjy/imbed/transform"
	"github.com/spf13/pflag"
)

func main() {
	registry := transform.NewRegistry()
	contrib.Register(registry)
	commands := app.Commands{}.
		Register(cmds.InitCommand{}.Spec()).
		Register(cmds.AddCommand{}.Spec()).
		Register(cmds.QCommand{}.Spec()).
		Register(cmds.RevCommand{}.Spec()).
		Register(cmds.GcCommand{}.Spec())
	err := app.ParseAndRun(os.Args, commands, registry)
	if err != nil {
		switch {
		case errors.Is(err, pflag.ErrHelp):

		default:
			fmt.Fprintf(app.Stderr, "fatal: %s\n", pos.FmtError(err))
		}
		os.Exit(2)
	}
}
