package luao

import lua "github.com/hsfzxjy/gopher-lua"

type Protocol interface {
	L2G(v lua.LValue) (Object, error)
	checkIntegrity(obj Object) error
}
