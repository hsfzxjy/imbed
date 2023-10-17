package schema

import (
	"bytes"

	"github.com/tinylib/msgp/msgp"
)

func EncodeBytes[S any](schema Schema[S], value S) ([]byte, error) {
	var buf bytes.Buffer
	var w = msgp.NewWriter(&buf)
	err := schema.EncodeMsg(w, value)
	if err != nil {
		return nil, err
	}
	if err = w.Flush(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
