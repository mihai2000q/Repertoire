package logger

import (
	"bytes"

	"go.uber.org/zap"
)

type GinLogger struct {
	*Logger
}

func NewGinLogger(logger *Logger) *GinLogger {
	logger = &Logger{logger.WithOptions(zap.WithCaller(false))}
	return &GinLogger{logger}
}

func (l *GinLogger) Write(p []byte) (n int, err error) {
	l.Info(string(bytes.TrimRight(p, "\n")))
	return len(p), nil
}
