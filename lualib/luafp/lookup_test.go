package luafp_test

import (
	"testing"

	lua "github.com/hsfzxjy/gopher-lua"
	"github.com/hsfzxjy/imbed/lualib/luafp"
	"github.com/stretchr/testify/assert"
)

func Test_LookupLocalFunction(t *testing.T) {
	const code = `
	xxx = 0
	local yyy = 1
	local function add(a)
		if d == nil then
			d = 1
		end
		print(d)
		return a + d
	end
	`
	fn := compile(code)
	AssertFunctionProtoEquals(t, fn.FunctionPrototypes[0], luafp.LookupLocalFunction(fn, "add"))
	assert.Nil(t, luafp.LookupLocalFunction(fn, "add2"))

	state := lua.NewState()
	f := state.NewFunctionFromProto(luafp.LookupLocalFunction(fn, "add"))
	f.Env = lua.AutoEnv()
	state.Push(f.AsLValue())
	state.Push(lua.LNumber(1).AsLValue())
	state.PCall(1, 1, nil)
	assert.True(t, lua.LNumber(2).AsLValue().Equals(state.Get(-1)))
	state.Pop(1)
	assert.True(t, state.Get(-1).EqualsLNil())
	state.Push(f.AsLValue())
	state.Push(lua.LNumber(1).AsLValue())
	state.PCall(1, 1, nil)
	assert.True(t, lua.LNumber(2).AsLValue().Equals(state.Get(-1)))
}
