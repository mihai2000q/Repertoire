package domain

import (
	"repertoire/server/domain/provider"
	"repertoire/server/domain/service"
	"repertoire/server/domain/usecase"

	"go.uber.org/fx"
)

var providers = fx.Options(
	fx.Provide(provider.NewCurrentUserProvider),
	fx.Provide(provider.NewStorageFilePathProvider),
)

var services = fx.Options(
	fx.Provide(service.NewAlbumService),
	fx.Provide(service.NewArtistService),
	fx.Provide(service.NewAuthService),
	fx.Provide(service.NewPlaylistService),
	fx.Provide(service.NewSongService),
	fx.Provide(service.NewUserService),
)

var Module = fx.Options(
	providers,
	usecase.Module,
	services,
)
