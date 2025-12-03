package logger

import (
	"strings"

	"go.uber.org/zap"
)

type ConfigOption func(c *zap.Config)

func WithLevel(level string) ConfigOption {
	return func(c *zap.Config) {
		switch strings.ToUpper(level) {
		case "INFO":
			c.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		case "WARN", "WARNING":
			c.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
		case "ERROR":
			c.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
		case "PANIC":
			c.Level = zap.NewAtomicLevelAt(zap.PanicLevel)
		case "FATAL":
			c.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
		default:
			c.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		}
	}
}

func WithDebugLevel() ConfigOption {
	return func(c *zap.Config) {
		c.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
}

func WithOutputStdout() ConfigOption {
	return func(c *zap.Config) {
		c.OutputPaths = []string{"stdout"}
		c.ErrorOutputPaths = []string{"stderr"}
	}
}

func WithEncodingJSON() ConfigOption {
	return func(c *zap.Config) {
		c.Encoding = "json"
	}
}
