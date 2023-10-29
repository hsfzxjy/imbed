package parser

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type state struct {
	i       int
	lastErr struct {
		error
		span
	}
}

type Parser struct {
	buf string
	state
}

func New(input []string) *Parser {
	buf := strings.Join(input, " ")
	return &Parser{buf: buf}
}
func NewString(input string) *Parser {
	return &Parser{buf: input}
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

func (p *Parser) span(off int) span {
	return span{p.i, p.i + off}
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
	p.ClearLastErr()
	if s := p.current(); s != "" && s[0] == b {
		p.advance(1)
		ok = true
	} else {
		p.setLastErr(byteError(b))
	}
	return
}

func (p *Parser) Term(term string) (ok bool) {
	if p == nil {
		return
	}
	p.ClearLastErr()
	if s := p.current(); s != "" && strings.HasPrefix(s, term) {
		p.advance(len(term))
		ok = true
	} else {
		p.setLastErr(termError(term))
	}
	return
}

func (p *Parser) AnyByte(charset string) (matched byte, ok bool) {
	if p == nil {
		return
	}
	p.ClearLastErr()
	if s := p.current(); s != "" {
		if b := s[0]; strings.IndexByte(charset, b) >= 0 {
			p.advance(1)
			matched, ok = b, true
			return
		}
	}
	p.setLastErr(anyByteError(charset))
	return
}

func (p *Parser) Space() {
	if p == nil {
		return
	}
	buf := p.buf
	i := p.i
	p.ClearLastErr()
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
	p.ClearLastErr()
	s := p.current()
	var off int
LOOP:
	for off = 0; off < len(s); off++ {
		switch b := s[off]; {
		case b == '-':
		case '0' <= b && b <= '9':
		default:
			break LOOP
		}
	}
	p.setLastErrSpan(p.span(off))
	if off == 0 {
		p.setLastErrString("expect 64-bit integer (e.g. '42')")
		return
	}
	if x, err := strconv.ParseInt(s[:off], 10, 64); err == nil {
		value, ok = x, true
		p.advance(off)
	} else {
		p.setLastErr(numError{err})
	}
	return
}

func (p *Parser) Rat() (value *big.Rat, ok bool) {
	if p == nil {
		return
	}
	p.ClearLastErr()
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
	if off < len(s) && (s[off] == '.' || s[off] == '/') {
		off++
	}
	for ; off < len(s); off++ {
		if b := s[off]; !('0' <= b && b <= '9') {
			break
		}
	}
	p.setLastErrSpan(p.span(off))
	if off == 0 {
		p.setLastErrString("expect rational number (e.g. '3/5', '3.14')")
		return
	}
	if rat, good := new(big.Rat).SetString(s[:off]); good {
		value, ok = rat, true
		p.advance(off)
	} else {
		p.setLastErrString("illegal rational number")
	}
	return
}

func (p *Parser) Ident() (value string, ok bool) {
	if p == nil {
		return
	}
	p.ClearLastErr()
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

	p.setLastErrSpan(p.span(off))

	if off == 0 {
		p.setLastErrString("expect identifier (e.g. 'foo.bar')")
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
	p.ClearLastErr()
	s := p.current()
	if len(s) == 0 {
		p.setLastErrString("expect string")
		return
	}
	switch s[0] {
	case '[':
		return p.quotedString(']')
	case '"':
		return p.quotedString('"')
	case '\'':
		return p.quotedString('\'')
	default:
		return p.unquotedString(stopRunes)
	}
}

func (p *Parser) quotedString(closing rune) (value string, ok bool) {
	s := p.current()
	var (
		off      int = 1
		escaping bool
		closed   bool
		b        strings.Builder
	)
LOOP:
	for off < len(s) {
		r, size := utf8.DecodeRuneInString(s[off:])
		if r == utf8.RuneError {
			if size == 1 {
				b.WriteByte(s[off])
				off += size
			}
			escaping = false
			continue
		}
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
			escaping = false
		} else {
			switch r {
			case '\\':
				escaping = true
				off += size
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
	p.setLastErrSpan(p.span(off))
	if !closed {
		p.setLastErr(fmt.Errorf("expect %q, string unclosed", closing))
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
			off += size
		case unicode.IsSpace(r):
			break LOOP
		case strings.ContainsRune(stopRunes, r):
			break LOOP
		default:
			off += size
		}
	}
	p.setLastErrSpan(p.span(off))
	if off == 0 {
		p.setLastErrString("expect string")
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
	p.ClearLastErr()
	s := p.current()
	if len(s) >= 4 && s[:4] == "true" {
		value, ok = true, true
		p.setLastErrSpan(p.span(4))
		p.advance(4)
	} else if len(s) >= 5 && s[:5] == "false" {
		value, ok = false, true
		p.setLastErrSpan(p.span(5))
		p.advance(5)
	} else {
		p.setLastErrString("expect 'true' or 'false'")
	}
	return
}

func (p *Parser) Rest() string {
	p.ClearLastErr()
	result := p.current()
	p.setLastErrSpan(p.span(len(p.buf)))
	p.i = len(p.buf)
	return result
}

func (p *Parser) EOF() bool {
	return p == nil || p.i == len(p.buf)
}

func (p *Parser) Snapshot() any {
	return p.state
}

func (p *Parser) Reset(snapshot any) {
	p.state = snapshot.(state)
}
