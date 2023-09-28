package app

import (
	"github.com/hsfzxjy/imbed/core"
	"github.com/spf13/pflag"
)

type Commands struct {
	specs []CommandSpec
}

type CommandRunner func(app *App, command CommandSpec) error

func (cs Commands) Register(spec CommandSpec) Commands {
	fs := pflag.NewFlagSet(spec.Name, pflag.ExitOnError)
	if spec.FlagSet != nil {
		fs.AddFlagSet(spec.FlagSet)
	}
	spec.FlagSet = fs
	cs.specs = append(cs.specs, spec)
	return cs
}

type CommandSpec struct {
	Name string
	*pflag.FlagSet
	core.Mode
	Runner CommandRunner
}
