package revparse

import (
	"bytes"
	"strings"

	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/transform"
	"github.com/tinylib/msgp/msgp"
)

func Solve(models []*db.AssetModel, registry *transform.Registry) (string, error) {
	var builder strings.Builder
	first := true
	for i, model := range models {
		if i == 0 {
			builder.WriteString("db.oid@")
			builder.WriteString(model.OID.FmtHumanize())
			builder.WriteString(" ")
		}
		buf := bytes.NewBuffer(model.TransSeqRaw)
		r := msgp.NewReader(buf)
		for buf.Len() > 0 || r.Buffered() > 0 {
			t, err := registry.DecodeMsg(r)
			if err != nil {
				return "", err
			}
			if !first {
				builder.WriteString(", ")
			}
			first = false
			builder.WriteString(t.Metadata().Name())
			builder.WriteByte('@')
			v := NewVisitor(&builder)
			builder.WriteString(t.ConfigHash().FmtHumanize())
			err = t.Visit(&v)
			if err != nil {
				return "", err
			}
		}
	}
	return builder.String(), nil
}
