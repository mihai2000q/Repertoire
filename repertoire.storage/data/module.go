package data

import (
	"repertoire/storage/data/logger"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var loggers = fx.Options(
	fx.Provide(logger.NewLogger),
	fx.Provide(logger.NewFxLogger),
	fx.Provide(logger.NewGinLogger),
	fx.WithLogger(func(logger *logger.FxLogger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: logger.Logger.Logger}
	}),
)

var Module = fx.Options(
	loggers,
)
