package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewConfig(opts ...ConfigOption) zap.Config {
	config := zap.NewDevelopmentConfig()

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	config.EncoderConfig.LineEnding = "\n"
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	config.DisableCaller = true

	for _, opt := range opts {
		opt(&config)
	}

	return config
}
