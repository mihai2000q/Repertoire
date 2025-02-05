package data

import (
	"go.uber.org/fx"
	"repertoire/server/data/database"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
)

var repositories = fx.Options(
	fx.Provide(repository.NewAlbumRepository),
	fx.Provide(repository.NewArtistRepository),
	fx.Provide(repository.NewPlaylistRepository),
	fx.Provide(repository.NewSongRepository),
	fx.Provide(repository.NewUserRepository),
)

var services = fx.Options(
	fx.Provide(service.NewBCryptService),
	fx.Provide(service.NewJwtService),
	fx.Provide(service.NewStorageService),
)

var Module = fx.Options(
	fx.Provide(database.NewClient),
	repositories,
	services,
)
