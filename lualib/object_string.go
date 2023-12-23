package lualib

import (
	"errors"

	lua "github.com/yuin/gopher-lua"
)

type String struct {
	v *lua.LString
}

func (s String) AsLValue(L *lua.LState, readonly bool) lua.LValue {
	return packLString(s.v)
}

func (s String) IsIntegral() error {
	return nil
}

type StringChecker struct{}

func (StringChecker) check(L *lua.LState, v lua.LValue) (Object, error) {
	if s, ok := unpackLValue[String](v); ok {
		return s, nil
	}
	return nil, errors.New("not a string")
}
