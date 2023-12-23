package lualib

import (
	"sync"

	lua "github.com/yuin/gopher-lua"
)

var autoEnv = sync.OnceValue(func() *lua.LTable {
	T := new(lua.LTable)
	mt := new(lua.LTable)
	mt.RawSetString("__index", &lua.LFunction{
		IsG:       true,
		Env:       T,
		GFunction: autoEnv__index,
	})
	T.Metatable = mt
	return T
})

func autoEnv__index(L *lua.LState) int {
	key := L.CheckString(2)
	L.Push(L.GetField(L.Env, key))
	return 1
}

func newGFunction(gfunc lua.LGFunction) *lua.LFunction {
	return &lua.LFunction{
		IsG:       true,
		Env:       autoEnv(),
		GFunction: gfunc,
	}
}

func newUserData(value any) *lua.LUserData {
	return &lua.LUserData{Value: value, Env: autoEnv()}
}

func newTable() *lua.LTable {
	return (*lua.LState).NewTable(nil)
}
