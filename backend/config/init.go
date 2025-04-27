package config

import (
	"fmt"
	"go.uber.org/zap"
)

var (
	Cfg    *Config
	Logger *zap.Logger
)

// InitConfig initializes the configuration from the given path
func InitConfig(configPath string) error {
	var err error

	// load config
	Cfg, err = LoadConfigFromPath(configPath)
	if err != nil {
		return fmt.Errorf("load config error: %v", err)
	}

	return nil
}

// GetConfig returns the current configuration
func GetConfig() *Config {
	return Cfg
}

// InitLogger initializes the logger
func InitLogger() error {
	Logger = CreateLogger()
	return nil
}