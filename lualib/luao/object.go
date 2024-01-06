package luao

import (
	lua "github.com/hsfzxjy/gopher-lua"
)

type lv = lua.LValue

type Object struct {
	lv
}

func (o Object) AsLValue() lua.LValue {
	return o.lv
}
