package core

import (
	"net/url"

	"github.com/hsfzxjy/imbed/schema"
)

type WorkspaceConfigProvider interface {
	WorkspaceConfigScanner(key string) (schema.Scanner, error)
}

type App interface {
	DBDir() string
	TmpDir() string
	WorkDir() string
	ConfigFilePath() string
	FilePath(filename string) string
	Mode() Mode
	BuildMode() BuildMode

	WorkspaceConfigProvider
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
