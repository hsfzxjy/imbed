package lualib

import (
	"errors"
	"unsafe"

	lua "github.com/yuin/gopher-lua"
)

type Int struct {
	v *lua.LNumber
}

func (i Int) AsLValue(L *lua.LState, readonly bool) lua.LValue {
	return packLNumber(i.v)
}

func (i Int) IsIntegral() error {
	return nil
}

func (i Int) asPtr() unsafe.Pointer {
	return unsafe.Pointer(i.v)
}

type IntChecker struct{}

func (IntChecker) check(L *lua.LState, v lua.LValue) (Object, error) {
	if i, ok := unpackLValue[Int](v); ok {
		return i, nil
	}
	return nil, errors.New("not a number")
}

func (IntChecker) ptrAsObject(ptr unsafe.Pointer) Object {
	return Int{(*lua.LNumber)(ptr)}
}
