package logger

import (
	"github.com/ThreeDotsLabs/watermill"
	"go.uber.org/zap"
)

type WatermillLogger struct {
	*Logger
}

func NewWatermillLogger(logger *Logger) *WatermillLogger {
	logger = &Logger{logger.WithOptions(
		zap.WithCaller(false),
	)}
	return &WatermillLogger{logger}
}

func (z *WatermillLogger) Error(msg string, err error, fields watermill.LogFields) {
	parsedFields := []zap.Field{zap.Error(err)}
	parsedFields = append(parsedFields, toZapFields(fields)...)
	z.Logger.Error("[Watermill] "+msg, parsedFields...)
}

func (z *WatermillLogger) Info(msg string, fields watermill.LogFields) {
	z.Logger.Info("[Watermill] "+msg, toZapFields(fields)...)
}

func (z *WatermillLogger) Debug(msg string, fields watermill.LogFields) {
	z.Logger.Debug("[Watermill] "+msg, toZapFields(fields)...)
}

func (z *WatermillLogger) Trace(msg string, fields watermill.LogFields) {
	z.Logger.Debug("[Watermill] "+msg, toZapFields(fields)...) // Use Debug for Trace since zap doesn't have a Trace level
}

func (z *WatermillLogger) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &WatermillLogger{
		&Logger{z.Logger.With(toZapFields(fields)...)},
	}
}

func toZapFields(fields watermill.LogFields) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	return zapFields
}
