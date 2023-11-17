package transform

import (
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/core/ref/lazy"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/schema"
	cfgf "github.com/hsfzxjy/imbed/transform/configfactory"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type view struct {
	md         *metadata
	params     any
	cfgFactory cfgf.Factory
}

func (r *view) Name() string {
	return r.md.name
}

func (r *view) ConfigHash(ctx db.Context) ref.Sha256 {
	return r.cfgFactory.ConfigHash(ctx)
}

func (r *view) VisitParams(v schema.Visitor) error {
	return r.md.paramsSchema.VisitAny(v, r.params)
}

func (r *view) Build(cp core.ConfigProvider) (*Transform, error) {
	cfg, err := r.cfgFactory.CreateConfig(cp)
	if err != nil {
		return nil, err
	}
	applier, err := r.md.buildApplier(r.params, cfg)
	if err != nil {
		return nil, err
	}
	return &Transform{
		Name:     r.md.name,
		Applier:  applier,
		Category: r.md.category,
		model: db.StepTpl{
			Config: db.ConfigTpl{
				O: lazy.DataFunc(func() []byte {
					var w fastbuf.W
					r.md.configSchema.EncodeMsgAny(&w, cfg)
					return w.Result()
				}),
			},
			Params: lazy.DataFuncSHAFunc(func() []byte {
				var w fastbuf.W
				w.WriteString(r.md.name)
				r.md.paramsSchema.EncodeMsgAny(&w, r.params)
				return w.Result()
			}, func() ref.Sha256 {
				var w fastbuf.W
				w.WriteString(r.md.name)
				applier.EncodeMsg(&w)
				return ref.Sha256Sum(w.Result())
			}),
		},
	}, nil
}
