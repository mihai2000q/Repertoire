package data

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"repertoire/server/data/cache"
	"repertoire/server/data/database"
	"repertoire/server/data/http"
	"repertoire/server/data/logger"
	"repertoire/server/data/message"
	"repertoire/server/data/repository"
	"repertoire/server/data/search"
	"repertoire/server/data/service"
)

var loggers = fx.Options(
	fx.Provide(logger.NewLogger),
	fx.Provide(logger.NewFxLogger),
	fx.Provide(logger.NewGinLogger),
	fx.Provide(logger.NewGormLogger),
	fx.Provide(logger.NewRestyLogger),
	fx.Provide(logger.NewWatermillLogger),
	fx.WithLogger(func(logger *logger.FxLogger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: logger.Logger.Logger}
	}),
)

var httpClients = fx.Options(
	fx.Provide(http.NewRestyClient),
	fx.Provide(http.NewStorageClient),
)

var repositories = fx.Options(
	fx.Provide(repository.NewAlbumRepository),
	fx.Provide(repository.NewArtistRepository),
	fx.Provide(repository.NewPlaylistRepository),
	fx.Provide(repository.NewSongRepository),
	fx.Provide(repository.NewUserDataRepository),
	fx.Provide(repository.NewUserRepository),
)

var services = fx.Options(
	fx.Provide(service.NewBCryptService),
	fx.Provide(service.NewJwtService),
	fx.Provide(service.NewMeiliTaskTrackerService),
	fx.Provide(service.NewMessagePublisherService),
	fx.Provide(service.NewSearchEngineService),
	fx.Provide(service.NewStorageService),
)

var Module = fx.Options(
	fx.Provide(cache.NewCache),
	loggers,
	fx.Provide(database.NewClient),
	httpClients,
	fx.Provide(message.NewPublisher),
	repositories,
	fx.Provide(search.NewMeiliClient),
	services,
)
