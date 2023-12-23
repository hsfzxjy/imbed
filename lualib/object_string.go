package lualib

import (
	"errors"
	"unsafe"

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

func (s String) asPtr() unsafe.Pointer {
	return unsafe.Pointer(s.v)
}

type StringChecker struct{}

func (StringChecker) check(L *lua.LState, v lua.LValue) (Object, error) {
	if s, ok := unpackLValue[String](v); ok {
		return s, nil
	}
	return nil, errors.New("not a string")
}

func (StringChecker) ptrAsObject(ptr unsafe.Pointer) Object {
	return String{(*lua.LString)(ptr)}
}
