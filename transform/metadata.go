package transform

import (
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/schema"
	cfgf "github.com/hsfzxjy/imbed/transform/configfactory"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type _constraint[C any, P ParamFor[C]] struct{}

func (_constraint[C, P]) buildApplier(params, config any) (Applier, error) {
	return params.(P).BuildTransform(config.(C))
}

type constraint interface {
	buildApplier(params, config any) (Applier, error)
}

type metadata[C any, P ParamFor[C]] struct {
	*Registry
	name    string
	aliases []string

	paramsSchema schema.Schema[P]
	configSchema schema.Schema[C]

	category Category
}

func (m *metadata[C, P]) Name() string {
	return m.name
}

func (m *metadata[C, P]) Category() Category {
	return m.category
}

func (m *metadata[C, P]) scanFrom(scanner schema.Scanner, copt cfgf.Opt) (View, error) {
	params, err := m.paramsSchema.ScanFrom(scanner)
	if err != nil {
		return nil, err
	}
	return &view[C, P]{
		md:     m,
		params: params,
		cfgOpt: copt,
	}, nil
}

func (m *metadata[C, P]) decodeMsg(reader *fastbuf.R, cfgOid ref.OID) (View, error) {
	params, err := m.paramsSchema.DecodeMsg(reader)
	if err != nil {
		return nil, err
	}
	return &view[C, P]{
		md:     m,
		params: params,
		cfgOpt: cfgf.OID(cfgOid),
	}, nil
}
