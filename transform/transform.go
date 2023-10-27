package transform

import "github.com/hsfzxjy/imbed/core/ref"

type Transform struct {
	Name string
	Applier
	Category
	Data   *Data
	Config ref.EncodableObject

	ForceTerminal bool
}
