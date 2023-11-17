package asset

import (
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/db"
)

type Transform interface {
	Model() db.StepListTpl
}

type Applier interface {
	Apply(app core.App, asset Asset) (Update, error)
}
