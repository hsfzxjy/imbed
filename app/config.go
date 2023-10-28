package app

import (
	"os"

	toml "github.com/pelletier/go-toml/v2"
)

const DB_DIR = ".imbed"
const TMP_DIR = "tmp"
const FILES_DIR = "files"
const CONFIG_FILENAME = "imbed.toml"

func loadConfigFile(filepath string) (map[string]any, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	decoder := toml.NewDecoder(f)
	var ret = map[string]any{}
	err = decoder.Decode(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
