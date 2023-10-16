package parser

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Parser struct {
	buf string
	i   int
}

func New(input []string) *Parser {
	buf := strings.Join(input, " ")
	return &Parser{buf, 0}
}
func NewString(input string) *Parser {
	return &Parser{input, 0}
}

func (p *Parser) advance(off int) {
	p.i = min(p.i+off, len(p.buf))
}

func (p *Parser) current() string {
	if p.i >= len(p.buf) {
		return ""
	}
	return p.buf[p.i:]
}

func (p *Parser) PeekByte() byte {
	if p.EOF() {
		return 0
	}
	return p.buf[p.i]
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

func (p *Parser) Term(term string) (ok bool) {
	if p == nil {
		return
	}
	if s := p.current(); s != "" && strings.HasPrefix(s, term) {
		p.advance(len(term))
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
	i := p.i
LOOP:
	for i < len(buf) {
		r, size := utf8.DecodeRuneInString(buf[i:])
		switch {
		case r == utf8.RuneError:
			break LOOP
		case !unicode.IsSpace(r):
			break LOOP
		}
		i += size
	}
	p.i = i
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

func (p *Parser) Rat() (value *big.Rat, ok bool) {
	if p == nil {
		return
	}
	s := p.current()
	var off int
	if off < len(s) && s[off] == '-' {
		off++
	}
	for ; off < len(s); off++ {
		if b := s[off]; !('0' <= b && b <= '9') {
			break
		}
	}
	if off < len(s) && s[off] == '.' || s[off] == '/' {
		off++
	}
	for ; off < len(s); off++ {
		if b := s[off]; !('0' <= b && b <= '9') {
			break
		}
	}
	if off == 0 {
		return
	}
	if rat, good := new(big.Rat).SetString(s[:off]); good {
		value, ok = rat, true
		p.advance(off)
	}
	return
}

func (p *Parser) Ident() (value string, ok bool) {
	if p == nil {
		return
	}
	s := p.current()
	var off int
LOOP:
	for off = 0; off < len(s); off++ {
		switch b := s[off]; {
		case b == '_' || b == '.':
		case 'a' <= b && b <= 'z':
		case 'A' <= b && b <= 'Z':
		case '0' <= b && b <= '9':
		default:
			break LOOP
		}
	}

	if off == 0 {
		return
	}

	value, ok = s[:off], true
	p.advance(off)
	return
}

// String matches unquoted string, quoted string and bracket string at
// the beginning.
// Unquoted string stops when encountering white spaces or any rune
// in stopRunes.
func (p *Parser) String(stopRunes string) (value string, ok bool) {
	if p == nil {
		return
	}
	s := p.current()
	if len(s) == 0 {
		return
	}
	switch s[0] {
	case '[':
		p.advance(1)
		return p.quotedString(']')
	case '"':
		p.advance(1)
		return p.quotedString('"')
	case '\'':
		p.advance(1)
		return p.quotedString('\'')
	default:
		return p.unquotedString(stopRunes)
	}
}

func (p *Parser) quotedString(closing rune) (value string, ok bool) {
	s := p.current()
	var (
		off      int
		escaping bool
		closed   bool
		b        strings.Builder
	)
LOOP:
	for off < len(s) {
		r, size := utf8.DecodeRuneInString(s[off:])
		if escaping {
			switch r {
			case 'a':
				r = '\a'
			case 'b':
				r = '\b'
			case 'f':
				r = '\f'
			case 'n':
				r = '\n'
			case 'r':
				r = '\r'
			case 't':
				r = '\t'
			case 'v':
				r = '\v'
			}
			b.WriteRune(r)
			off += size
		} else {
			switch r {
			case '\\':
				escaping = true
			case closing:
				off += size
				closed = true
				break LOOP
			default:
				b.WriteRune(r)
				off += size
			}
		}
	}
	if !closed {
		return
	}
	value, ok = b.String(), true
	p.advance(off)
	return
}

func (p *Parser) unquotedString(stopRunes string) (value string, ok bool) {
	s := p.current()
	var off int
LOOP:
	for off < len(s) {
		r, size := utf8.DecodeRuneInString(s[off:])
		switch {
		case r == utf8.RuneError:
			break LOOP
		case unicode.IsSpace(r):
			break LOOP
		case strings.ContainsRune(stopRunes, r):
			break LOOP
		default:
			off += size
		}
	}

	if off == 0 {
		return
	}
	value, ok = s[:off], true
	p.advance(off)
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

func (p *Parser) Rest() string {
	result := p.current()
	p.i = len(p.buf)
	return result
}

func (p *Parser) EOF() bool {
	return p == nil || p.i == len(p.buf)
}

func (p *Parser) Error(err error) error {
	var perr *parserError
	if errors.As(err, &perr) {
		return err
	}
	return &parserError{p, err}
}

func (p *Parser) Expect(expected string) error {
	return p.Error(fmt.Errorf("expect %s", expected))
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
	return fmt.Sprintf("%s\n\t%s\n\t% *s", e.error.Error(), p.buf, p.i+1, "^")
}

func (e *parserError) Unwrap() error { return e.error }
