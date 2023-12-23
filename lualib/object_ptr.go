package lualib

import (
	"unsafe"

	lua "github.com/yuin/gopher-lua"
)

type Ptr struct {
	object Object
}

func (p *Ptr) AsLValue(L *lua.LState, readonly bool) lua.LValue {
	if p.object == nil {
		return lua.LNil
	}
	return p.object.AsLValue(L, readonly)
}

func (p *Ptr) IsIntegral() error {
	if p.object == nil {
		return nil
	}
	return p.object.IsIntegral()
}

func (p *Ptr) asPtr() unsafe.Pointer {
	return unsafe.Pointer(p)
}

type PtrChecker struct {
	C objectChecker
}

func (c *PtrChecker) check(L *lua.LState, v lua.LValue) (Object, error) {
	if v == lua.LNil {
		return &Ptr{}, nil
	}
	if o, err := c.C.check(L, v); err != nil {
		return nil, err
	} else {
		return &Ptr{object: o}, nil
	}
}

func (c *PtrChecker) ptrAsObject(ptr unsafe.Pointer) Object {
	return (*Ptr)(ptr)
}
