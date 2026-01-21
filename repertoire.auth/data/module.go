package data

import (
	"repertoire/auth/data/database"
	"repertoire/auth/data/logger"
	"repertoire/auth/data/repository"
	"repertoire/auth/data/service"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var loggers = fx.Options(
	fx.Provide(logger.NewLogger),
	fx.Provide(logger.NewFxLogger),
	fx.Provide(logger.NewGinLogger),
	fx.Provide(logger.NewGormLogger),
	fx.WithLogger(func(logger *logger.FxLogger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: logger.Logger.Logger}
	}),
)

var services = fx.Options(
	fx.Provide(service.NewBCryptService),
	fx.Provide(service.NewJwtService),
)

var Module = fx.Options(
	loggers,
	fx.Provide(repository.NewUserRepository),
	services,
	fx.Provide(database.NewClient),
)
