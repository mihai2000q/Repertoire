package logger

import (
	"go.uber.org/zap"
	"repertoire/server/data/logger"
)

func NewLoggerMock() *logger.Logger {
	return &logger.Logger{Logger: zap.Must(zap.NewDevelopment())}
}
