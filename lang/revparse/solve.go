package revparse

import (
	"strings"

	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/transform"
)

func Solve(models []*db.AssetModel, registry transform.Registry) (string, error) {
	var builder strings.Builder
	first := true
	for _, model := range models {
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
