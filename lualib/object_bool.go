package lualib

import (
	"errors"

	lua "github.com/yuin/gopher-lua"
)

type Bool struct {
	v *lua.LBool
}

func (b Bool) AsLValue(L *lua.LState, readonly bool) lua.LValue {
	return packLBool(b.v)
}

func (b Bool) IsIntegral() error {
	return nil
}

type BoolChecker struct{}

func (BoolChecker) check(L *lua.LState, v lua.LValue) (Object, error) {
	if b, ok := unpackLValue[Bool](v); ok {
		return b, nil
	}
	return nil, errors.New("not a bool")
}
