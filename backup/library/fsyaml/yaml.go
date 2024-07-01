package fsyaml

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Env struct {
		Port string `yaml:"port"`
	} `yaml:"env,omitempty"`
}

func ReadEnv(cfg *Config, path string) error {
	filename, _ := filepath.Abs(path)
	f, err := os.ReadFile(filename)
	if nil != err {
		return err
	}

	err = yaml.Unmarshal(f, cfg)
	if nil != err {
		return err
	}

	return nil
}

func Init(cfg *Config, path string) error {
	var err error
	err = ReadEnv(cfg, path)
	if nil != err {
		return err
	}
	return nil
}
