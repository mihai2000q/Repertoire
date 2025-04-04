package logger

import (
	"go.uber.org/zap"
)

type FxLogger struct {
	*Logger
}

func NewFxLogger(logger *Logger) *FxLogger {
	logger = &Logger{logger.WithOptions(zap.WithCaller(false))}
	return &FxLogger{logger}
}
