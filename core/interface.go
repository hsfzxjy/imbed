package core

import (
	"net/url"

	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/schema"
)

type ConfigProvider interface {
	ProvideWorkspaceConfig(key string) (schema.Scanner, error)
	ProvideStockConfig(ndl.Needle) ([]byte, error)
	ProvideConfigByOID(ref.OID) ([]byte, error)
}

type App interface {
	DBDir() string
	TmpDir() string
	WorkDir() string
	ConfigFilePath() string
	FilePath(filename string) string
	Mode() Mode
	BuildMode() BuildMode

	ProvideWorkspaceConfig(key string) (schema.Scanner, error)
	ProxyFunc() func(reqURL *url.URL) (*url.URL, error)
}

type Mode int

const (
	ModeNone Mode = iota
	ModeReadonly
	ModeReadWrite
)

type BuildMode int

const (
	BuildUseCache BuildMode = iota
	BuildRedeploy
	BuildAll
)
