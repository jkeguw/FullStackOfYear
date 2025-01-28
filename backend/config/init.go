package config

import (
	"fmt"
	"go.uber.org/zap"
)

var (
	Cfg    *Config
	Logger *zap.Logger
)

func Init() error {
	var err error

	if err := InitViper(); err != nil {
		return fmt.Errorf("init viper error: %v", err)
	}

	// load config
	Cfg, err = LoadConfig()
	if err != nil {
		return fmt.Errorf("load config error: %v", err)
	}

	// init logger
	Logger = InitLogger()

	return nil
}
