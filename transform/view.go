package transform

import (
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/core/ref/lazy"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/schema"
	cfgf "github.com/hsfzxjy/imbed/transform/configfactory"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type view[C any, P ParamFor[C]] struct {
	md     *metadata[C, P]
	params P
	cfgOpt cfgf.Opt
}

func (r *view[C, P]) Name() string {
	return r.md.name
}

func (r *view[C, P]) ConfigHash(ctx db.Context) ref.Sha256 {
	return r.cfgOpt.ConfigHash(ctx)
}

func (r *view[C, P]) VisitParams(v schema.Visitor) error {
	return r.md.paramsSchema.Visit(v, r.params)
}

func (r *view[C, P]) buildConfig(cp ConfigProvider) (C, db.ConfigTpl, error) {
	var cfg C
	var tpl db.ConfigTpl
	if r.cfgOpt.FromDB() {
		cfgModel, err := r.cfgOpt.QueryModel(cp)
		if err != nil {
			return cfg, tpl, err
		}
		if cfgModel != nil {
			var reader = fastbuf.R{Buf: cfgModel.Data}
			cfg, err = r.md.configSchema.DecodeMsg(&reader)
			if err != nil {
				return cfg, tpl, err
			}
		}
		tpl = db.ConfigTpl{O: cfgModel}
	} else {
		scanner, err := cp.WorkspaceConfigScanner(r.md.name)
		if err != nil {
			return cfg, tpl, err
		}
		cfg, err = r.md.configSchema.ScanFrom(scanner)
		if err != nil {
			return cfg, tpl, err
		}
		tpl = db.ConfigTpl{O: lazy.DataFunc(func() []byte {
			var w fastbuf.W
			r.md.configSchema.EncodeMsg(&w, cfg)
			return w.Result()
		})}
	}
	return cfg, tpl, nil
}

func (r *view[C, P]) Build(cp ConfigProvider) (*Transform, error) {
	cfg, cfgTpl, err := r.buildConfig(cp)
	if err != nil {
		return nil, err
	}
	applier, err := r.params.BuildTransform(cfg)
	if err != nil {
		return nil, err
	}
	return &Transform{
		Name:     r.md.name,
		Applier:  applier,
		Category: r.md.category,
		model: db.StepTpl{
			Config: cfgTpl,
			Params: lazy.DataFuncSHAFunc(func() []byte {
				var w fastbuf.W
				w.WriteString(r.md.name)
				r.md.paramsSchema.EncodeMsg(&w, r.params)
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
