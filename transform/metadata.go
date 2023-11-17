package transform

import (
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/schema"
	cfgf "github.com/hsfzxjy/imbed/transform/configfactory"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type metadata[C any, P ParamFor[C]] struct {
	*Registry
	name    string
	aliases []string

	paramsSchema schema.Schema[P]
	configSchema schema.Schema[C]

	category Category
}

func (m *metadata[C, P]) scanFrom(scanner schema.Scanner, copt cfgf.Opt) (View, error) {
	params, pos, err := m.paramsSchema.ScanFrom(scanner)
	if err != nil {
		return nil, pos.WrapError(err)
	}
	return &view[C, P]{
		md:        m,
		paramsPos: pos,
		params:    params,
		cfgOpt:    copt,
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
