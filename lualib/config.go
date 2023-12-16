package lualib

import (
	"unsafe"

	lua "github.com/yuin/gopher-lua"
)

const (
	lconstNumber = iota
	lconstString
)

// sync with: gopher-lua.FunctionProto
type luaFunctionProto struct {
	SourceName         string
	LineDefined        int
	LastLineDefined    int
	NumUpvalues        uint8
	NumParameters      uint8
	IsVarArg           uint8
	NumUsedRegisters   uint8
	Code               []uint32
	Constants          []lua.LValue
	FunctionPrototypes []*lua.FunctionProto

	DbgSourcePositions []int
	DbgLocals          []*lua.DbgLocalInfo
	DbgCalls           []lua.DbgCall
	DbgUpvalues        []string

	stringConstants []string
}

func (p *luaFunctionProto) asLua() *lua.FunctionProto {
	var _ [unsafe.Sizeof(luaFunctionProto{}) - unsafe.Sizeof(lua.FunctionProto{})]struct{}
	var _ [unsafe.Sizeof(lua.FunctionProto{}) - unsafe.Sizeof(luaFunctionProto{})]struct{}

	return (*lua.FunctionProto)(unsafe.Pointer(p))
}
