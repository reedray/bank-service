package logger

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
)

type Logger interface {
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
	Fatal(message string, args ...interface{})
}

type loggerImpl struct {
	logger *zap.Logger
}

func NewLogger(level string) (Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	switch strings.ToLower(level) {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("can`t create logger %w", err)
	}
	return &loggerImpl{logger: logger}, nil
}

func (l *loggerImpl) Fatal(message string, args ...interface{}) {
	l.logger.Fatal(message)
}

func (l *loggerImpl) Debug(message string, args ...interface{}) {

	l.logger.Debug(message)
}

func (l *loggerImpl) Info(message string, args ...interface{}) {
	l.logger.Info(message)
}

func (l *loggerImpl) Warn(message string, args ...interface{}) {
	l.logger.Warn(message)
}

func (l *loggerImpl) Error(message string, args ...interface{}) {

	l.logger.Error(message)

}
