package db

import (
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/core/ref/lazy"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type StepTpl struct {
	Config ConfigTpl

	Params lazy.MustDataSHAObject
}

type stepListTpl struct {
	list []*StepTpl
	lazy.SHAF
	supportsRemove bool
}

func (t *stepListTpl) computeSHA() ref.Sha256 {
	var szHash fastbuf.Size
	szHash.Reserve(2 * len(t.list) * ref.Sha256{}.Sizeof())
	w := szHash.Build()
	for _, t := range t.list {
		w.
			AppendRaw(t.Config.O.MustSHA().Raw()).
			AppendRaw(t.Params.MustSHA().Raw())
	}
	return ref.Sha256HashSum(w.Result())
}

func (t *stepListTpl) create(tx *Tx) ([]*ConfigModel, []byte, error) {
	var szData fastbuf.Size
	for _, t := range t.list {
		szData.
			Reserve(ref.OID{}.Sizeof()).
			ReserveBytes(t.Params.MustData())
	}
	wData := szData.Build()
	cfgModels := make([]*ConfigModel, len(t.list))
	for i, t := range t.list {
		cfgModel, err := t.Config.create(tx)
		if err != nil {
			return nil, nil, err
		}
		cfgModels[i] = cfgModel
		wData.
			AppendRaw(cfgModel.OID.Raw()).
			WriteBytes(t.Params.MustData())
	}
	return cfgModels, wData.Result(), nil
}

func (t *stepListTpl) SupportsRemove() bool {
	return t.supportsRemove
}

type StepListTpl interface {
	MustSHA() ref.Sha256
	SupportsRemove() bool
	create(tx *Tx) ([]*ConfigModel, []byte, error)
}

func NewStepListTpl(list []*StepTpl, SupportsRemove bool) StepListTpl {
	t := &stepListTpl{list: list}
	t.SHAF.Fn = t.computeSHA
	return t
}
