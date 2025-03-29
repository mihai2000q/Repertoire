package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"repertoire/server/internal"
	"sync"
	"time"
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

	finalCore := consoleCore

	// if the app is running normally, not in integration testing env
	// then create and append the file encoder
	if os.Getenv("INTEGRATION_TESTING_ENVIRONMENT_FILE_PATH") == "" {
		lumberjackLogger := &lumberjack.Logger{
			Filename: getLogFile(env.LogOutput),
			MaxSize:  100,
			MaxAge:   31,
			Compress: true,
		}
		fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(lumberjackLogger), levelsMap[env.LogLevel])

		go dailyRotation(lumberjackLogger, env)

		finalCore = zapcore.NewTee(finalCore, fileCore)
	}

	return &Logger{zap.New(
		finalCore,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)}
}

func getLogFile(logOutput string) string {
	if err := os.MkdirAll(logOutput, os.ModePerm); err != nil {
		panic(err)
	}
	return filepath.Join(logOutput, time.Now().Format("2006-01-02")+".log")
}

var mu sync.Mutex

func dailyRotation(lj *lumberjack.Logger, env internal.Env) {
	for {
		now := time.Now()
		nextMidnight := now.Truncate(24 * time.Hour).Add(24 * time.Hour)
		time.Sleep(time.Until(nextMidnight))

		mu.Lock()
		lj.Filename = getLogFile(env.LogOutput)
		_ = lj.Rotate()
		mu.Unlock()
	}
}
