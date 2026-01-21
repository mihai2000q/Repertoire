package logger

import (
	"repertoire/server/data/logger"

	"go.uber.org/zap"
)

func NewLoggerMock() *logger.Logger {
	return &logger.Logger{Logger: zap.Must(zap.NewDevelopment())}
}
