package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CreateLogger creates and returns a new zap logger
func CreateLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, _ := config.Build()
	return logger
}