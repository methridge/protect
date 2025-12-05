package logger

import (
	"go.uber.org/zap"
)

var globalLogger *zap.SugaredLogger

// New creates and returns a new logger instance
func New() *zap.SugaredLogger {
	if globalLogger != nil {
		return globalLogger
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.FatalLevel + 1) // Disable by default

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	globalLogger = logger.Sugar()
	return globalLogger
}

// Get returns the global logger instance
func Get() *zap.SugaredLogger {
	if globalLogger == nil {
		return New()
	}
	return globalLogger
}

// SetLevel sets the logging level
func SetLevel(level string) error {
	var zapLevel zap.AtomicLevel

	switch level {
	case "none":
		// Disable all logging
		zapLevel = zap.NewAtomicLevelAt(zap.FatalLevel + 1)
	case "debug":
		zapLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		zapLevel = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		// Default to none
		zapLevel = zap.NewAtomicLevelAt(zap.FatalLevel + 1)
	}

	config := zap.NewProductionConfig()
	config.Level = zapLevel

	logger, err := config.Build()
	if err != nil {
		return err
	}

	if globalLogger != nil {
		globalLogger.Sync()
	}

	globalLogger = logger.Sugar()
	return nil
}
