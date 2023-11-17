package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hsfzxjy/imbed/parser"
	"github.com/stretchr/testify/assert"
)

func dedent(expected string) string {
	expected = strings.TrimLeft(expected, "\n")
	lastNL := strings.LastIndex(expected, "\n")
	return expected[:lastNL+1]
}

func Test_Parser_Int64(t *testing.T) {
	var out strings.Builder
	test := func(input string) {
		p := parser.NewString(input)
		p.Space()
		x, ok := p.Int64()
		fmt.Fprintln(&out, "===")
		if ok {
			fmt.Fprintf(&out, "%d\n", x)
		} else {
			fmt.Fprintln(&out, p.ErrorString("foo").Error())
		}
	}
	test(" 1234444444444444444")
	test(" 12344444444444444444")
	test(" ")
	expected := `
===
1234444444444444444
===
foo: expect 64-bit integer: value out of range
	|  12344444444444444444
	|  ^^^^^^^^^^^^^^^^^^^^
===
foo: expect 64-bit integer (e.g. '42')
	|  
	|  ^
	`
	assert.Equal(t, dedent(expected), out.String())
}

func Test_Parser_Rat(t *testing.T) {
	var out strings.Builder
	test := func(input string) {
		p := parser.NewString(input)
		p.Space()
		x, ok := p.Rat()
		fmt.Fprintln(&out, "===")
		if ok {
			fmt.Fprintf(&out, "%s\n", x.RatString())
		} else {
			fmt.Fprintln(&out, p.ErrorString("foo").Error())
		}
	}
	test(" 1234/4321")
	test(" 12344444444444444444")
	test(" 1.2")
	test(" 1.")
	test(" 1/")
	test(" 1/0")
	test(" ")
	expected := `
===
1234/4321
===
12344444444444444444
===
6/5
===
1
===
foo: illegal rational number
	|  1/
	|  ^^
===
foo: illegal rational number
	|  1/0
	|  ^^^
===
foo: expect rational number (e.g. '3/5', '3.14')
	|  
	|  ^
	`
	assert.Equal(t, dedent(expected), out.String())
}

func Test_Parser_String(t *testing.T) {
	var out strings.Builder
	test := func(input string) {
		p := parser.NewString(input)
		p.Space()
		x, ok := p.String(" ,")
		fmt.Fprintln(&out, "===")
		if ok {
			fmt.Fprintf(&out, "%s\n", x)
		} else {
			fmt.Fprintln(&out, p.ErrorString("foo").Error())
		}
	}
	test(" a-b,")
	test(" [ a'\\[\\]]")
	test(" \" a'\\[\\]\"")
	test(" [ a'\\[\\]")
	test(" \xa9ab")
	test(" [\xa9ab\\\xa9]")
	test(" ,bar")
	expected := `
===
a-b
===
 a'[]
===
 a'[]
===
foo: expect ']', string unclosed
	|  [ a'\[\]
	|  ^^^^^^^^
===
` + "\xa9" + `ab
===
` + "\xa9" + `ab` + "\xa9" + `
===
foo: expect string
	|  ,bar
	|  ^
`
	assert.Equal(t, dedent(expected), out.String())
}

func Test_Parser_Bool(t *testing.T) {
	var out strings.Builder
	test := func(input string) {
		p := parser.NewString(input)
		p.Space()
		x, ok := p.Bool()
		fmt.Fprintln(&out, "===")
		if ok {
			fmt.Fprintf(&out, "%v\n", x)
		} else {
			fmt.Fprintln(&out, p.ErrorString("foo").Error())
		}
	}
	test(" true")
	test(" false")
	test(" tru")
	expected := `
===
true
===
false
===
foo: expect 'true' or 'false'
	|  tru
	|  ^
`
	assert.Equal(t, dedent(expected), out.String())
}

func Test_Parser_Byte(t *testing.T) {
	var out strings.Builder
	test := func(input string, b byte) {
		p := parser.NewString(input)
		p.Space()
		ok := p.Byte(b)
		fmt.Fprintln(&out, "===")
		if !ok {
			fmt.Fprintln(&out, p.ErrorString("foo").Error())
		}
	}
	test(" =", '=')
	test(" a", '=')
	test(" ", '=')
	expected := `
===
===
foo: expect '='
	|  a
	|  ^
===
foo: expect '='
	|  
	|  ^
`
	assert.Equal(t, dedent(expected), out.String())
}

func Test_Parser_AnyByte(t *testing.T) {
	var out strings.Builder
	test := func(input string, charset string) {
		p := parser.NewString(input)
		p.Space()
		b, ok := p.AnyByte(charset)
		fmt.Fprintln(&out, "===")
		if ok {
			fmt.Fprintf(&out, "%c\n", rune(b))
		} else {
			fmt.Fprintln(&out, p.ErrorString("foo").Error())
		}
	}
	test(" =a", "=")
	test(" =a", "=,")
	test(" a", "=,")
	expected := `
===
=
===
=
===
foo: expect any of '=', ','
	|  a
	|  ^
`
	assert.Equal(t, dedent(expected), out.String())
}

func Test_Parser_Term(t *testing.T) {
	var out strings.Builder
	test := func(input string, term string) {
		p := parser.NewString(input)
		p.Space()
		ok := p.Term(term)
		fmt.Fprintln(&out, "===")
		if !ok {
			fmt.Fprintln(&out, p.ErrorString("foo").Error())
		}
	}
	test(" oid", "oid")
	test(" '", "'")
	test(" oi", "oid'")
	expected := `
===
===
===
foo: expect "oid'"
	|  oi
	|  ^
`
	assert.Equal(t, dedent(expected), out.String())
}
