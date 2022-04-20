package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(loglevel string) (*zap.Logger, error) {
	cfg := zap.Config{
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	if loglevel == "error" {
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}
	if loglevel == "info" {
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	if loglevel == "debug" {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		cfg.DisableCaller = false
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
