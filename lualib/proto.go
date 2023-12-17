package lualib

import (
	"slices"

	lua "github.com/yuin/gopher-lua"
)

func opGetCode(code uint32) int {
	return int(code >> 26)
}

func opGetArgA(code uint32) int {
	return int(code>>18) & 0xff
}

func opGetArgBx(code uint32) int {
	return int(code & 0x3ffff)
}

func LookupLocalFunction(p *lua.FunctionProto, name string) *lua.FunctionProto {
	idx := slices.IndexFunc(p.DbgLocals, func(dli *lua.DbgLocalInfo) bool {
		return dli.Name == name
	})
	if idx == -1 {
		return nil
	}
	for _, code := range p.Code {
		if opGetCode(code) == lua.OP_CLOSURE && opGetArgA(code) == idx {
			return p.FunctionPrototypes[opGetArgBx(code)]
		}
	}
	return nil
}
