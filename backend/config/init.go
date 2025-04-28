package config

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"strings"
)

var (
	Cfg    *Config
	Logger *zap.Logger
)

// InitConfig initializes the configuration
func InitConfig() error {
	// 检查是否有环境变量CONFIG_FILE
	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		// 默认配置文件路径
		configPath = "./config/config.yaml"
	}

	var err error

	// 根据配置文件结尾判断使用哪种加载方法
	if strings.HasSuffix(configPath, ".yaml") || strings.HasSuffix(configPath, ".yml") {
		// load config from yaml file
		Cfg, err = LoadConfigFromPath(configPath)
		if err != nil {
			return fmt.Errorf("load config error: %v", err)
		}
	} else {
		// 使用viper加载
		err = InitViper()
		if err != nil {
			return fmt.Errorf("init viper error: %v", err)
		}

		// 从viper加载配置
		Cfg = &Config{}
		// 从这里开始绑定viper配置到Cfg
		// ...
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
