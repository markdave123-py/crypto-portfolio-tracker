package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetupLogger() (*zap.Logger, error) {
	logLevel := zap.InfoLevel
	if os.Getenv("ENVIRONMENT") == "local" {
		logLevel = zap.DebugLevel
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(logLevel), // Set log level to debug
		Development:      true,                           // Enables DPanic level logging
		Encoding:         "json",                         // or "console" for human-readable logs
		OutputPaths:      []string{"stdout"},             // Specify where to output logs
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "level",
			TimeKey:       "time",
			CallerKey:     "caller",
			FunctionKey:   zapcore.OmitKey,
			StacktraceKey: "stacktrace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.CapitalLevelEncoder, // "debug", "info", "warn" etc.
			EncodeTime:    zapcore.ISO8601TimeEncoder,  // Use ISO8601 format for time
			EncodeCaller:  zapcore.ShortCallerEncoder,  // Short file name + line number
		},
	}

	baseLogger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("could not build logger: %w", err)
	}

	return baseLogger, nil
}
