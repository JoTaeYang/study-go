package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/JoTaeYang/study-go/pkg/db/ddb"
	"gopkg.in/yaml.v2"
)

// Config ...
type Config struct {
	Auth struct {
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_access_key"`
	} `yaml:"auth,omitempty"`

	DynamoDB struct {
		Region string `yaml:"region"`
		Table  string `yaml:"table"`
	} `yaml:"ddb,omitempty"`
}

func ReadConfig(conf interface{}, path string) (interface{}, error) {
	filename, _ := filepath.Abs(path)
	f, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = yaml.Unmarshal(f, conf)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return conf, nil
}

func InitConfig(cfg *Config) error {
	var err error
	if cfg.DynamoDB.Table != "" {
		ddbConfig := ddb.Config{
			Region:    cfg.DynamoDB.Region,
			AccessKey: cfg.Auth.AccessKey,
			SecretKey: cfg.Auth.SecretKey,
			Table:     cfg.DynamoDB.Table,
		}

		if cfg.DynamoDB.Region == "local" {
			err = ddb.LocalInit(ddbConfig)
		} else {
			err = ddb.Init(ddbConfig)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
