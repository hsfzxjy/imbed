package transform

import (
	"fmt"

	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/schema"
	cfgf "github.com/hsfzxjy/imbed/transform/configfactory"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type Registry struct {
	metadataTable map[string]*metadata
	composerTable map[Category]Composer
}

func NewRegistry() *Registry {
	return &Registry{
		map[string]*metadata{},
		map[Category]Composer{},
	}
}

func (r *Registry) RegisterComposer(cat Category, composer Composer) *Registry {
	if _, ok := r.composerTable[cat]; ok {
		panic("Composer with category " + cat + "is already registered")
	}
	r.composerTable[cat] = composer
	return r
}

func (r *Registry) ScanFrom(name string, scanner schema.Scanner, copt cfgf.Opt) (*view, error) {
	m, ok := r.metadataTable[name]
	if !ok {
		return nil, fmt.Errorf("no transform named %q", name)
	}
	params, err := m.paramsSchema.ScanFromAny(scanner)
	if err != nil {
		return nil, err
	}
	return &view{
		md:         m,
		params:     params,
		cfgFactory: copt(name, m.configSchema),
	}, nil
}

func (r *Registry) DecodeMsg(stepData db.Step) (*view, error) {
	var reader = fastbuf.R{Buf: stepData.Data}
	name, err := reader.ReadString()
	if err != nil {
		return nil, err
	}
	md, ok := r.Lookup(name)
	if !ok {
		return nil, fmt.Errorf("no transform named %q", name)
	}
	params, err := md.paramsSchema.DecodeMsgAny(&reader)
	if err != nil {
		return nil, err
	}
	return &view{
		md:         md,
		params:     params,
		cfgFactory: cfgf.OID(stepData.ConfigOID)(name, md.configSchema),
	}, nil
}

func (r *Registry) Lookup(name string) (*metadata, bool) {
	m, ok := r.metadataTable[name]
	return m, ok
}
