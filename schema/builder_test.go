package schema_test

import (
	"bytes"
	"fmt"

	"github.com/hsfzxjy/imbed/schema"
	schemareader "github.com/hsfzxjy/imbed/schema/reader"
	"github.com/tinylib/msgp/msgp"
)

type X struct {
	int64
	bool
	float64
	string
	m map[string]int64
}

func ExampleNew() {
	var x = X{}
	s := schema.Struct(&x,
		schema.F("int", &x.int64, schema.Int()),
		schema.F("bool", &x.bool, schema.Bool()),
		schema.F("float", &x.float64, schema.Float()),
		schema.F("m", &x.m, schema.Map(schema.Int())),
		schema.F("str", &x.string, schema.String())).
		DebugName("X")
	sch := schema.New(s)
	err := sch.DecodeValue(schemareader.Any(map[string]any{
		"int":   int64(1),
		"bool":  true,
		"m":     map[string]any{"a": (1)},
		"float": float64(3.14),
		"str":   "test",
	}), &x)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", x)

	var buf bytes.Buffer
	w := msgp.NewWriter(&buf)
	err = sch.EncodeMsg(w, &x)
	if err != nil {
		panic(err)
	}
	w.Flush()
	x = X{}
	r := msgp.NewReader(&buf)
	err = sch.DecodeMsg(r, &x)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", x)

	// Output:
	// schema_test.X{int64:1, bool:true, float64:3.14, string:"test", m:map[string]int64{"a":1}}
	// schema_test.X{int64:1, bool:true, float64:3.14, string:"test", m:map[string]int64{"a":1}}
}
