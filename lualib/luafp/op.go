package luafp

import (
	"slices"
	"strings"

	lua "github.com/hsfzxjy/gopher-lua"
	"github.com/hsfzxjy/gopher-lua/parse"
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

func Compile(code string, filename string) (*lua.FunctionProto, error) {
	reader := strings.NewReader(code)
	chunk, err := parse.Parse(reader, filename)
	if err != nil {
		return nil, err
	}
	return lua.Compile(chunk, filename)
}
