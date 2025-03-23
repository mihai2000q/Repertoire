package logger

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormlogger "gorm.io/gorm/logger"
	"time"
)

var gormLevelsMap = map[zapcore.Level]gormlogger.LogLevel{
	zapcore.DebugLevel: gormlogger.Info,
	zapcore.InfoLevel:  gormlogger.Info,
	zapcore.WarnLevel:  gormlogger.Warn,
	zapcore.ErrorLevel: gormlogger.Error,
}

type GormLogger struct {
	*Logger
	gormlogger.Config
}

func NewGormLogger(logger *Logger) *GormLogger {
	logger = &Logger{logger.WithOptions(
		zap.WithCaller(false),
	)}

	return &GormLogger{
		logger,
		gormlogger.Config{
			LogLevel: gormLevelsMap[logger.Level()],
		},
	}
}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	l.LogLevel = level
	return l
}

func (l *GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.Logger.Sugar().Infof("[GORM] "+str, args...)
	}
}

func (l *GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.Logger.Sugar().Warnf("[GORM] "+str, args...)
	}
}

func (l *GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.Logger.Sugar().Errorf("[GORM] "+str, args...)
	}
}

// Trace prints trace messages
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil {
		l.Logger.Error("[GORM] SQL query failed",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed),
			zap.Error(err),
		)
	} else {
		l.Logger.Debug("[GORM] SQL query executed",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed),
		)
	}
}
