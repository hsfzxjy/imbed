package cmds

import (
	"fmt"
	"os"

	"github.com/hsfzxjy/imbed/app"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/util"
)

type InitCommand struct{}

func (InitCommand) Spec() app.CommandSpec {
	return app.CommandSpec{
		Name:    "init",
		FlagSet: nil,
		Mode:    core.ModeNone,
		Runner:  InitCommand{}.Run,
	}
}

func (InitCommand) Run(app *app.App, spec app.CommandSpec) error {
	dbDir := app.DBDir()
	if util.IsDir(dbDir) {
		return fmt.Errorf("%s has been initialized", app.WorkDir())
	}
	err := os.Mkdir(dbDir, 0o700)
	if err != nil {
		return err
	}
	err = os.Mkdir(app.TmpDir(), 0o700)
	if err != nil {
		return err
	}
	err = os.Mkdir(app.FilePath(""), 0o700)
	if err != nil {
		return err
	}
	configFilePath := app.ConfigFilePath()
	if !util.IsFile(configFilePath) {
		f, err := os.Create(configFilePath)
		if err != nil {
			return nil
		}
		defer f.Close()
	}
	return nil
}
