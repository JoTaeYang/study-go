package main

import (
	"errors"

	"github.com/JoTaeYang/study-go/pkg/service"
)

var (
	configDefaultPath = "./"
	configDefaultName = "config.yaml"
	cfg               = service.Config{}
)

func InitConfig() error {
	_, err := service.ReadConfig(&cfg, configDefaultPath+configDefaultName)
	if err != nil {
		//TODO ::session manager file load
		return errors.New("service config error")
	}

	service.InitConfig(&cfg)
	return nil
}
