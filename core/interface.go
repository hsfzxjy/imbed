package core

import (
	"net/url"

	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/schema"
)

type ConfigProvider interface {
	ProvideWorkspaceConfig(key string) (schema.Reader, error)
	ProvideStockConfig(ndl.Needle) ([]byte, error)
}

type App interface {
	DBDir() string
	WorkDir() string
	ConfigFilePath() string
	FilePath(filename string) string
	Mode() Mode
	BuildMode() BuildMode

	ProvideWorkspaceConfig(key string) (schema.Reader, error)
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

type Iterator[T any] interface {
	HasNext() bool
	Next() T
}
