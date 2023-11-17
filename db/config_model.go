package db

import (
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/core/ref/lazy"
)

type SHA = lazy.ConstSHA
type Data = lazy.ConstData

type ConfigModel struct {
	OID ref.OID
	SHA
	Data
}
