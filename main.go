package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/hsfzxjy/imbed/app"
	"github.com/hsfzxjy/imbed/cmds"
	"github.com/hsfzxjy/imbed/contrib"
	"github.com/hsfzxjy/imbed/transform"
	"github.com/spf13/pflag"
)

func main() {
	contrib.Register(transform.DefaultRegistry())
	commands := app.Commands{}.
		Register(cmds.InitCommand{}.Spec()).
		Register(cmds.AddCommand{}.Spec()).
		Register(cmds.QCommand{}.Spec()).
		Register(cmds.RevCommand{}.Spec())
	err := app.ParseAndRun(os.Args, commands)
	if err != nil {
		switch {
		case errors.Is(err, pflag.ErrHelp):

		default:
			fmt.Fprintf(os.Stderr, "fatal: %s\n", err)
		}
		os.Exit(2)
	}
}
