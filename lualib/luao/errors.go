package luao

import (
	"strconv"

	lua "github.com/hsfzxjy/gopher-lua"
)

type typeError struct {
	Expected, Got lua.LValueType
}

func (e *typeError) Error() string {
	return ("expected " + e.Expected.String() + ", got " + e.Got.String())
}

type intError struct{ v lua.LNumber }

func (e *intError) Error() string {
	return ("expected int64, got " + strconv.FormatFloat(float64(e.v), 'g', -1, 64))
}

type noSuchFieldError struct{ name string }

func (e *noSuchFieldError) Error() string {
	return ("no such field: " + e.name)
}
