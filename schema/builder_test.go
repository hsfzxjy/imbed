package schema_test

import (
	"fmt"
	"math/big"

	"github.com/hsfzxjy/imbed/schema"
	schemascanner "github.com/hsfzxjy/imbed/schema/scanner"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type X struct {
	int64
	bool
	*big.Rat
	string
	m map[string]int64
}

func (x X) GoString() string {
	return fmt.Sprintf("%d %v %s %s %#v", x.int64, x.bool, x.RatString(), x.string, x.m)
}

func ExampleNew() {
	var x = new(X)
	s := schema.Struct(x,
		schema.F("int", &x.int64, schema.Int()),
		schema.F("bool", &x.bool, schema.Bool()),
		schema.F("float", &x.Rat, schema.Rat()),
		schema.F("m", &x.m, schema.Map(schema.Int())),
		schema.F("str", &x.string, schema.String())).
		DebugName("X")
	sch := schema.New(s)
	x, _, err := sch.ScanFrom(schemascanner.Any(map[string]any{
		"int":   int64(1),
		"bool":  true,
		"m":     map[string]any{"a": int64(1)},
		"float": "3.14",
		"str":   "test",
	}))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", x)

	var w fastbuf.W
	sch.EncodeMsg(&w, x)
	r := fastbuf.R{Buf: w.Result()}
	x, err = sch.DecodeMsg(&r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", x)

	// Output:
	// 1 true 157/50 test map[string]int64{"a":1}
	// 1 true 157/50 test map[string]int64{"a":1}
}
