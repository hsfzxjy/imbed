package parser

import (
	"errors"
	"strconv"
	"strings"
)

type span struct{ start, end int }

func (p *Parser) ClearLastErr() {
	p.lastErr.error = nil
	p.lastErr.span = span{-1, -1}
}

func (p *Parser) setLastErr(err error) {
	p.lastErr.error = err
}

func (p *Parser) setLastErrString(msg string) {
	p.lastErr.error = errors.New(msg)
}

func (p *Parser) setLastErrSpan(span span) {
	p.lastErr.span = span
}

func (p *Parser) Error(err error) error {
	var perr *parserError
	if errors.As(err, &perr) {
		return err
	}
	sp := p.lastErr.span
	if sp.start < 0 {
		sp = span{p.i, p.i + 1}
	}
	if sp.start == sp.end {
		sp.end++
	}
	return &parserError{p.buf, err, p.lastErr.error, sp}
}

func (p *Parser) ErrorString(msg string) error {
	return p.Error(errors.New(msg))
}

type parserError struct {
	buf        string
	err1, err2 error
	span
}

func (e *parserError) Error() string {
	var b strings.Builder
	var err1String, err2String string
	if e.err1 != nil {
		err1String = e.err1.Error()
	}
	if e.err2 != nil {
		err2String = e.err2.Error()
	}
	b.Grow(len(err1String) + len(err2String) + len(e.buf)*2 + 10)
	b.WriteString(err1String)
	if err2String != "" {
		if err1String != "" {
			b.WriteString(": ")
		}
		b.WriteString(err2String)
	}
	b.WriteString("\n\t| ")
	b.WriteString(e.buf)
	b.WriteString("\n\t| ")
	var i int
	for i = 0; i < e.start; i++ {
		b.WriteRune(' ')
	}
	for ; i < e.end; i++ {
		b.WriteRune('^')
	}
	return b.String()
}

func (e *parserError) Unwrap() []error {
	var errs = make([]error, 0, 2)
	if e.err1 != nil {
		errs = append(errs, e.err1)
	}
	if e.err2 != nil {
		errs = append(errs, e.err2)
	}
	return errs
}

type byteError byte

func (e byteError) Error() string {
	return "expect " + strconv.QuoteRune(rune(e))
}

type anyByteError string

func (e anyByteError) Error() string {
	var b strings.Builder
	if len(e) == 1 {
		return byteError(e[0]).Error()
	}
	b.WriteString("expect any of ")
	for i, ch := range e {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(strconv.QuoteRune(ch))
	}
	return b.String()
}

type termError string

func (e termError) Error() string {
	return "expect " + strconv.Quote(string(e))
}

type numError struct{ err error }

func (e numError) Error() string {
	if nerr, ok := e.err.(*strconv.NumError); ok {
		return "expect 64-bit integer: " + nerr.Err.Error()
	} else {
		return e.err.Error()
	}
}
