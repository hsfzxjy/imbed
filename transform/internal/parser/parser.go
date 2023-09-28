package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	buf  []string
	i, j int
}

func New(input []string) *Parser {
	return &Parser{input, 0, 0}
}
func NewString(input string) *Parser {
	return &Parser{[]string{input}, 0, 0}
}

func (p *Parser) advance(off int) {
	p.j += off
	for p.i < len(p.buf) {
		s := p.buf[p.i]
		if p.j >= len(s) {
			p.j = 0
			p.i++
		} else {
			return
		}
	}
}

func (p *Parser) current() string {
	if p.i >= len(p.buf) {
		return ""
	}
	return p.buf[p.i][p.j:]
}

func (p *Parser) Byte(b byte) (ok bool) {
	if p == nil {
		return
	}
	if s := p.current(); s != "" && s[0] == b {
		p.advance(1)
		ok = true
	}
	return
}

func (p *Parser) AnyByte(charset string) (matched byte, ok bool) {
	if p == nil {
		return
	}
	if s := p.current(); s != "" {
		if b := s[0]; strings.IndexByte(charset, b) >= 0 {
			p.advance(1)
			matched, ok = b, true
		}

	}
	return
}

func (p *Parser) Space() {
	if p == nil {
		return
	}
	buf := p.buf
	i, j := p.i, p.j
LOOP:
	for i < len(buf) {
		s := buf[i]
		for j < len(s) {
			switch s[j] {
			case ' ', '\t':
				j++
			default:
				break LOOP
			}
		}
		i++
	}
	p.i, p.j = i, j
}

func (p *Parser) Int64() (value int64, ok bool) {
	if p == nil {
		return
	}
	s := p.current()
	var end int
LOOP:
	for end = 0; end < len(s); end++ {
		switch b := s[end]; {
		case b == '-':
		case '0' <= b && b <= '9':
		default:
			break LOOP
		}
	}
	if end == 0 {
		return
	}
	if x, err := strconv.ParseInt(s[:end], 10, 64); err == nil {
		value, ok = x, true
		p.advance(end)
	}
	return
}

func (p *Parser) Float64() (value float64, ok bool) {
	if p == nil {
		return
	}
	s := p.current()
	var end int
LOOP:
	for end = 0; end < len(s); end++ {
		switch b := s[end]; {
		case b == '-' || b == '.':
		case '0' <= b && b <= '9':
		default:
			break LOOP
		}
	}
	if end == 0 {
		return
	}
	if x, err := strconv.ParseFloat(s[:end], 64); err == nil {
		value, ok = x, true
		p.advance(end)
	}
	return
}

func (p *Parser) Ident() (value string, ok bool) {
	if p == nil {
		return
	}
	s := p.current()
	var end int
LOOP:
	for end = 0; end < len(s); end++ {
		switch b := s[end]; {
		case b == '_' || b == '.':
		case 'a' <= b && b <= 'z':
		case 'A' <= b && b <= 'Z':
		case '0' <= b && b <= '9':
		default:
			break LOOP
		}
	}

	if end == 0 {
		return
	}

	value, ok = s[:end], true
	p.advance(end)
	return
}

func (p *Parser) String() (value string, ok bool) {
	if p == nil {
		return
	}
	s := p.current()
	var end int
LOOP:
	for end = 0; end < len(s); end++ {
		switch s[end] {
		case '_', '.', '-':
			continue LOOP
		}
		switch b := s[end]; {
		case 'a' <= b && b <= 'z':
		case 'A' <= b && b <= 'Z':
		case '0' <= b && b <= '9':
		default:
			break LOOP
		}
	}

	if end == 0 {
		return
	}

	value, ok = s[:end], true
	p.advance(end)
	return
}

func (p *Parser) Bool() (value, ok bool) {
	if p == nil {
		return
	}
	s := p.current()
	if len(s) >= 4 && s[:4] == "true" {
		value, ok = true, true
		p.advance(4)
	} else if len(s) >= 5 && s[:5] == "false" {
		value, ok = false, true
		p.advance(5)
	}
	return
}

func (p *Parser) EOF() bool {
	return p == nil || p.i == len(p.buf)
}

func (p *Parser) Error(err error) error {
	if e, ok := err.(*parserError); ok {
		return e
	}
	return &parserError{p, err}
}

type parserError struct {
	*Parser
	error
}

func (e *parserError) Error() string {
	p := e.Parser
	if p == nil {
		return e.error.Error()
	}
	text := strings.Join(p.buf, " ")
	var pos int
	for k := 0; k < p.i; k++ {
		pos += len(p.buf[k]) + 1
	}
	if p.j > 0 {
		pos += p.j + 1
	}
	return fmt.Sprintf("%s\n\t%s\n\t% *s", e.error.Error(), text, pos, "^")
}

func (e *parserError) Unwrap() error { return e.error }
