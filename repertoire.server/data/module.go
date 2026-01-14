package data

import (
	"repertoire/server/data/cache"
	"repertoire/server/data/database"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/http"
	"repertoire/server/data/http/client"
	"repertoire/server/data/logger"
	"repertoire/server/data/message"
	"repertoire/server/data/realtime"
	"repertoire/server/data/repository"
	"repertoire/server/data/search"
	"repertoire/server/data/service"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var caches = fx.Options(
	fx.Provide(cache.NewMeiliCache),
	fx.Provide(cache.NewStorageCache),
	fx.Provide(cache.NewCentrifugoCache),
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

var databaseOptions = fx.Options(
	fx.Provide(database.NewClient),
	fx.Provide(transaction.NewManager),
)

var httpClients = fx.Options(
	fx.Provide(http.NewRestyClient),
	fx.Provide(client.NewAuthClient),
	fx.Provide(client.NewStorageClient),
)

var repositories = fx.Options(
	fx.Provide(repository.NewAlbumRepository),
	fx.Provide(repository.NewArtistRepository),
	fx.Provide(repository.NewPlaylistRepository),
	fx.Provide(repository.NewSongRepository),
	fx.Provide(repository.NewSongSectionRepository),
	fx.Provide(repository.NewUserDataRepository),
	fx.Provide(repository.NewUserRepository),
)

var services = fx.Options(
	fx.Provide(service.NewAuthService),
	fx.Provide(service.NewBCryptService),
	fx.Provide(service.NewJwtService),
	fx.Provide(service.NewSearchTaskTrackerService),
	fx.Provide(service.NewMessagePublisherService),
	fx.Provide(service.NewRealTimeService),
	fx.Provide(service.NewSearchEngineService),
	fx.Provide(service.NewStorageService),
)

var Module = fx.Options(
	caches,
	loggers,
	databaseOptions,
	httpClients,
	fx.Provide(message.NewPublisher),
	fx.Provide(realtime.NewCentrifugoClient),
	repositories,
	fx.Provide(search.NewMeiliClient),
	services,
)
