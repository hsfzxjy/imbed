package schema

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type schemaError struct {
	// in reversed order
	path []string
	err  error
	op   string
}

type schemaErrorWrapper struct{ e *schemaError }

func (w schemaErrorWrapper) Error() string {
	e := w.e
	var b strings.Builder
	if e.op != "" {
		b.WriteString(e.op)
		b.WriteString(": ")
	}

	if len(e.path) > 0 {
		b.WriteRune('[')
		for i := len(e.path) - 1; i > 0; i-- {
			b.WriteString(e.path[i])
			b.WriteString(" -> ")
		}
		b.WriteString(e.path[0])
		b.WriteString("]: ")
	}

	b.WriteString(e.err.Error())
	return b.String()
}

func (w schemaErrorWrapper) Unwrap() error {
	return w.e.err
}

func (e *schemaError) Unwrap() error { return e.err }
func (e *schemaError) AppendPath(part string) *schemaError {
	if e == nil {
		return nil
	}
	e.path = append(e.path, part)
	return e
}
func (e *schemaError) AppendIndex(i int) *schemaError {
	if e == nil {
		return nil
	}
	e.path = append(e.path, "["+strconv.Itoa(i)+"]")
	return e
}
func (e *schemaError) SetOp(op string) *schemaError {
	if e == nil {
		return nil
	}
	e.op = op
	return e
}
func (e *schemaError) AsError() error {
	if e == nil {
		return nil
	}
	return schemaErrorWrapper{e}
}
func newError(err error) *schemaError {
	switch e := err.(type) {
	case nil:
		return nil
	case schemaErrorWrapper:
		return e.e
	}
	return &schemaError{
		err: err,
	}
}

var ErrRequired = errors.New("value is required")
var ErrUnexpectedField = errors.New("unexpected field")
var ErrValidation = errors.New("validation error")

// func requiredErr(typ string, err error) error {
// 	return fmt.Errorf("%w (type=%s): %w", ErrRequired, typ, err)
// }

// func required(typ string) error {
// 	return fmt.Errorf("%w (type=%s)", ErrRequired, typ)
// }

func unexpectedField(badFieldName string) error {
	return fmt.Errorf("%w %q", ErrUnexpectedField, badFieldName)
}

func validation(err error) error {
	return fmt.Errorf("%w: %w", ErrValidation, err)
}

var ErrDecodeMsgBadSig = errors.New("bad type signature")

func badSig(expected, actual sig) error {
	return fmt.Errorf("%w: expected=%x, actual=%x", ErrDecodeMsgBadSig, expected, actual)
}
