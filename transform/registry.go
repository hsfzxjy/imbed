package transform

import (
	"fmt"

	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/schema"
	cfgf "github.com/hsfzxjy/imbed/transform/configfactory"
	"github.com/tinylib/msgp/msgp"
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

func (r *Registry) ScanFrom(name string, scanner schema.Scanner) (*Data, error) {
	m, ok := r.metadataTable[name]
	if !ok {
		return nil, fmt.Errorf("no transform named %q", name)
	}
	params, err := m.paramsSchema.ScanFromAny(scanner)
	if err != nil {
		return nil, err
	}
	return &Data{m, params}, nil
}

func (r *Registry) DecodeMsg(reader *msgp.Reader) (*Builder, error) {
	cfgHash, err := ref.FromReader[ref.Sha256Hash](reader)
	if err != nil {
		return nil, err
	}
	name, err := reader.ReadString()
	if err != nil {
		return nil, err
	}
	metadata, ok := r.Lookup(name)
	if !ok {
		return nil, fmt.Errorf("no transform named %q", name)
	}
	params, err := metadata.paramsSchema.DecodeMsgAny(reader)
	if err != nil {
		return nil, err
	}
	return &Builder{
		&Data{metadata, params},
		metadata.ConfigFactory(cfgf.Hash(cfgHash)),
	}, nil
}

func (r *Registry) Lookup(name string) (*metadata, bool) {
	m, ok := r.metadataTable[name]
	return m, ok
}
