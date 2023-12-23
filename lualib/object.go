package lualib

import (
	"unsafe"

	lua "github.com/yuin/gopher-lua"
)

type Object interface {
	IsIntegral() error
	AsLValue(L *lua.LState, readonly bool) lua.LValue
	asPtr() unsafe.Pointer
}

type objectChecker interface {
	check(L *lua.LState, v lua.LValue) (Object, error)
	ptrAsObject(unsafe.Pointer) Object
}

func packLNumber[T ~float64](x *T) lua.LValue {
	var i lua.LValue = lua.LNumber(0)
	(*struct {
		_, data unsafe.Pointer
	})(unsafe.Pointer(&i)).data = unsafe.Pointer(x)
	return i
}

func packLString[T ~string](x *T) lua.LValue {
	var i lua.LValue = lua.LString("")
	(*struct {
		_, data unsafe.Pointer
	})(unsafe.Pointer(&i)).data = unsafe.Pointer(x)
	return i
}

func packLBool[T ~bool](x *T) lua.LValue {
	var i lua.LValue = lua.LBool(false)
	(*struct {
		_, data unsafe.Pointer
	})(unsafe.Pointer(&i)).data = unsafe.Pointer(x)
	return i
}

func unpackLValue[S ~struct{ v *T }, T lua.LValue](v lua.LValue) (S, bool) {
	if _, ok := v.(T); ok {
		return (*struct {
			_    unsafe.Pointer
			data S
		})(unsafe.Pointer(&v)).data, true
	} else {
		return S{}, false
	}
}
