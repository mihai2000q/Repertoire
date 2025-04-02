package logger

import "go.uber.org/zap"

type RestyLogger struct {
	*Logger
}

func (r RestyLogger) Errorf(format string, v ...interface{}) {
	r.Logger.Sugar().Errorf("[HTTP] "+format, v...)
}

func (r RestyLogger) Warnf(format string, v ...interface{}) {
	r.Logger.Sugar().Warnf("[HTTP] "+format, v...)
}

func (r RestyLogger) Debugf(format string, v ...interface{}) {
	r.Logger.Sugar().Debugf("[HTTP] "+format, v...)
}

func NewRestyLogger(logger *Logger) *RestyLogger {
	logger = &Logger{logger.WithOptions(zap.WithCaller(false))}
	return &RestyLogger{logger}
}
