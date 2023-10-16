package revparse

import (
	"strings"

	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/transform"
)

func Solve(models []*db.AssetModel, registry transform.Registry) (string, error) {
	var builder strings.Builder
	first := true
	for i, model := range models {
		if i == 0 {
			builder.WriteString("db.oid@")
			builder.WriteString(model.OID.FmtHumanize())
			builder.WriteString(" ")
		}
		transforms, err := registry.DecodeParams(model.TransSeqRaw)
		if err != nil {
			return "", err
		}
		for _, t := range transforms {
			if !first {
				builder.WriteString(", ")
			}
			first = false
			builder.WriteString(t.Metadata().Name())
			builder.WriteByte('@')
			v := NewVisitor(&builder)
			builder.WriteString(t.ConfigHash().FmtHumanize())
			err := t.VisitParams(&v)
			if err != nil {
				return "", err
			}
		}
	}
	return builder.String(), nil
}
