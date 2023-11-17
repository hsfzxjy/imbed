package transform

import (
	"fmt"

	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/schema"
	cfgf "github.com/hsfzxjy/imbed/transform/configfactory"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type Registry struct {
	metadataTable map[string]Metadata
	composerTable map[Category]Composer
}

func NewRegistry() *Registry {
	return &Registry{
		make(map[string]Metadata, 32),
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

func (r *Registry) ScanFrom(name string, scanner schema.Scanner, copt cfgf.Opt) (View, error) {
	m, ok := r.metadataTable[name]
	if !ok {
		return nil, fmt.Errorf("no transform named %q", name)
	}
	return m.scanFrom(scanner, copt)
}

func (r *Registry) DecodeMsg(stepData db.Step) (View, error) {
	var reader = fastbuf.R{Buf: stepData.Data}
	name, err := reader.ReadString()
	if err != nil {
		return nil, err
	}
	md, ok := r.metadataTable[name]
	if !ok {
		return nil, fmt.Errorf("no transform named %q", name)
	}
	return md.decodeMsg(&reader, stepData.ConfigOID)
}
