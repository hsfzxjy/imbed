package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type span struct{ start, end int }

func (p *Parser) fmtFunc(str string, start, end uint32) string {
	var b strings.Builder
	b.Grow(len(str) + len(p.buf)*2 + 10)
	b.WriteString(str)
	b.WriteString("\n\t| ")
	b.WriteString(p.buf)
	b.WriteString("\n\t| ")
	var i uint32
	for i = 0; i < start; i++ {
		b.WriteRune(' ')
	}
	if end <= start {
		end = start + 1
	}
	for ; i < end; i++ {
		b.WriteRune('^')
	}
	return b.String()
}

type parseError struct {
	tag parseErrorTag

	byteBased   byte
	stringBased string
	raw         error
}

type parseErrorTag uint8

const (
	parseErrorNil parseErrorTag = iota
	parseErrorByte
	parseErrorAnyByte
	parseErrorTerm
	parseErrorNum
	parseErrorMsg
	parseErrorUnclosedString
)

func (p *parseError) Get() error {
	switch p.tag {
	case parseErrorByte:
		return byteError(p.byteBased)
	case parseErrorAnyByte:
		return anyByteError(p.stringBased)
	case parseErrorTerm:
		return termError(p.stringBased)
	case parseErrorNum:
		return numError{err: p.raw}
	case parseErrorMsg:
		return errorMsg(p.stringBased)
	case parseErrorUnclosedString:
		return unclosedStringError(p.byteBased)
	default:
		return nil
	}
}

func (p *parseError) SetByteError(err byteError) {
	p.tag = parseErrorByte
	p.byteBased = byte(err)
}

func (p *parseError) SetAnyByteError(err anyByteError) {
	p.tag = parseErrorAnyByte
	p.stringBased = string(err)
}

func (p *parseError) SetTermError(err termError) {
	p.tag = parseErrorTerm
	p.stringBased = string(err)
}

func (p *parseError) SetNumError(err error) {
	p.tag = parseErrorNum
	p.raw = err
}

func (p *parseError) SetMsgError(msg string) {
	p.tag = parseErrorMsg
	p.stringBased = string(msg)
}

func (p *parseError) SetUnclosedStringError(ch byte) {
	p.tag = parseErrorUnclosedString
	p.byteBased = ch
}

func (p *parseError) Clear() {
	p.tag = parseErrorNil
}

type errorMsg string

func (e errorMsg) Error() string {
	return string(e)
}

type unclosedStringError byte

func (e unclosedStringError) Error() string {
	return "expect " + strconv.QuoteRune(rune(e)) + " to close string"
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

func (p *Parser) Error(err error) error {
	if err == nil {
		return p.lastErr.Get()
	}
	return errors.Join(err, p.lastErr.Get())
}

func (p *Parser) ErrorString(msg string) error {
	err := p.lastErr.Get()
	if msg == "" {
		return err
	}
	if err == nil {
		return errors.New(msg)
	}
	return fmt.Errorf("%s: %w", msg, err)
}
