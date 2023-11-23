package cmds

import (
	"github.com/hsfzxjy/imbed/app"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/db/gc"
)

type GcCommand struct{}

func (GcCommand) Spec() app.CommandSpec {
	return app.CommandSpec{
		Name:   "gc",
		Mode:   core.ModeReadWrite,
		Runner: GcCommand{}.Run,
	}
}

type remover struct {
	db.AssetRemover
}

func (r remover) RemoveUrl(*db.AssetModel) error { return nil }

func (GcCommand) Run(app *app.App, command app.CommandSpec) error {
	return app.DB().RunRW(func(tx *db.Tx) error {
		return gc.GC(tx, remover{})
	})
}
