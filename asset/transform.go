package asset

import (
	"github.com/hsfzxjy/imbed/core/ref"
)

type Transform interface {
	Name() string
	Applier
	ref.EncodableObject
	AssociatedConfigs() []ref.EncodableObject
}

type Applier interface {
	Apply(asset Asset) (Update, error)
}
