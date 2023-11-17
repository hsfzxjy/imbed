package transform

import (
	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/db/configq"
	"github.com/hsfzxjy/imbed/schema"
	cfgf "github.com/hsfzxjy/imbed/transform/configfactory"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type ParamFor[C any] interface {
	BuildTransform(cfg C) (Applier, error)
}

type Applier interface {
	asset.Applier
	EncodeMsg(w *fastbuf.W)
}

type Composer interface {
	Compose(appliers []Applier) (Applier, error)
}

type ConfigProvider interface {
	configq.Provider
	core.WorkspaceConfigProvider
}

type Metadata interface {
	scanFrom(scanner schema.Scanner, copt cfgf.Opt) (View, error)
	decodeMsg(r *fastbuf.R, cfgOid ref.OID) (View, error)
}

type View interface {
	Name() string
	ConfigHash(ctx db.Context) ref.Sha256
	VisitParams(v schema.Visitor) error
	Build(cp ConfigProvider) (*Transform, error)
}
