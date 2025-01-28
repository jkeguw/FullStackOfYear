package config

import (
	"github.com/spf13/viper"
)

func InitViper() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// 环境变量配置
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CPC")

	return viper.ReadInConfig()
}
