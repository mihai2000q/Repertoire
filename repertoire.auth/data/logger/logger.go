package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"repertoire/auth/internal"
)

var levelsMap = map[string]zapcore.Level{
	"DEBUG": zapcore.DebugLevel,
	"INFO":  zapcore.InfoLevel,
	"WARN":  zapcore.WarnLevel,
	"ERROR": zapcore.ErrorLevel,
}

type Logger struct {
	*zap.Logger
}

func NewLogger(env internal.Env) *Logger {
	// console encoder
	consoleConfiguration := zap.NewDevelopmentEncoderConfig()
	consoleConfiguration.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleConfiguration),
		zapcore.AddSync(os.Stdout),
		levelsMap[env.LogLevel],
	)
	return &Logger{zap.New(
		consoleCore,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)}
}

func NewLoggerMock() *Logger {
	return &Logger{Logger: zap.Must(zap.NewDevelopment())}
}
