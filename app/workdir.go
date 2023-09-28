package app

import (
	"os"
	"path"
)

func sanitizeWorkDir(workDir string) (string, error) {
	if workDir == "" {
		return os.Getwd()
	}
	return workDir, nil
}

func isDir(p string) bool {
	info, err := os.Stat(p)
	return err == nil && info.IsDir()
}

func findWorkspace(workDir string) (string, bool) {
	for {
		tryDbDir := path.Join(workDir, DB_DIR)
		if isDir(tryDbDir) {
			return workDir, true
		}
		nextWorkDir := path.Dir(workDir)
		if workDir == nextWorkDir {
			return "", false
		}
		workDir = nextWorkDir
	}
}
