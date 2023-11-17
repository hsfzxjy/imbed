package revparse

import (
	"strings"

	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/transform"
)

func Solve(ctx db.Context, models []*db.AssetModel, registry *transform.Registry) (string, error) {
	var builder strings.Builder
	first := true
	for i, model := range models {
		if i == 0 {
			builder.WriteString("sha@")
			builder.WriteString(model.SHA.FmtHumanize())
			builder.WriteString(" ")
		}
		decoded := model.StepListData.Decode()
		for _, x := range decoded {
			view, err := registry.DecodeMsg(x)
			if err != nil {
				return "", err
			}
			if !first {
				builder.WriteString(", ")
			}
			first = false
			builder.WriteString(view.Name())
			builder.WriteByte('@')
			v := NewVisitor(&builder)
			builder.WriteString(view.ConfigHash(ctx).FmtHumanize())
			err = view.VisitParams(v)
			if err != nil {
				return "", err
			}
		}
	}
	return builder.String(), nil
}
