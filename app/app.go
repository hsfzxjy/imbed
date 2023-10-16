package app

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"sync"

	"golang.org/x/net/http/httpproxy"

	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/schema"
	schemascanner "github.com/hsfzxjy/imbed/schema/scanner"
	"github.com/spf13/pflag"
)

type App struct {
	mode    core.Mode
	workDir string
	dbDir   string
	cfgTree map[string]any

	dbs    db.Service
	dbOnce sync.Once

	proxyConfig *httpproxy.Config
	proxyOnce   sync.Once
}

func ParseAndRun(cmdArgs []string, specs Commands) error {
	var err error
	if len(cmdArgs) == 1 {
		return errors.New("no subcommand")
	}
	cmd := cmdArgs[1]
	for _, spec := range specs.specs {
		if spec.Name != cmd {
			continue
		}

		var workDir string
		global := pflag.NewFlagSet("Global Options", pflag.ContinueOnError)
		global.StringVarP(&workDir, "work-dir", "d", "", "Specify the working directory")

		flagSet := pflag.NewFlagSet(spec.Name, pflag.ContinueOnError)
		flagSet.AddFlagSet(global)
		flagSet.AddFlagSet(spec.FlagSet)

		if err = flagSet.Parse(cmdArgs[2:]); err != nil {
			return err
		}
		spec.FlagSet = flagSet

		if workDir, err = sanitizeWorkDir(workDir); err != nil {
			return err
		}

		var cfgTree map[string]any
		if spec.Mode != core.ModeNone {
			var ok bool
			workDir, ok = findWorkspace(workDir)
			if !ok {
				return errors.New("no workspace")
			}
			cfgTree, err = loadConfigFile(path.Join(workDir, CONFIG_FILENAME))
			if err != nil {
				return err
			}
		}

		app := &App{
			mode:    spec.Mode,
			workDir: workDir,
			dbDir:   path.Join(workDir, DB_DIR),
			cfgTree: cfgTree,
		}
		return spec.Runner(app, spec)
	}

	return fmt.Errorf("no command %s", cmd)
}

func (s *App) DBDir() string {
	return s.dbDir
}

func (s *App) WorkDir() string {
	return s.workDir
}

func (s *App) FilePath(objectName string) string {
	return path.Join(s.dbDir, FILES_DIR, objectName)
}

func (s *App) ConfigFilePath() string {
	return path.Join(s.workDir, CONFIG_FILENAME)
}

func (s *App) Mode() core.Mode {
	return s.mode
}

func (s *App) ProvideWorkspaceConfig(key string) (schema.Scanner, error) {
	var cfg any
	if s.cfgTree != nil {
		cfg = s.cfgTree[key]
	}
	return schemascanner.Any(cfg), nil
}

func (s *App) BuildMode() core.BuildMode {
	panic("TODO")
}

func (s *App) DB() db.Service {
	s.dbOnce.Do(func() {
		dbs, err := db.Open(s)
		if err != nil {
			panic(err)
		}
		s.dbs = dbs
	})
	return s.dbs
}

func (s *App) ProxyFunc() func(reqURL *url.URL) (*url.URL, error) {
	s.proxyOnce.Do(func() {
		s.proxyConfig = httpproxy.FromEnvironment()
	})
	return s.proxyConfig.ProxyFunc()
}
