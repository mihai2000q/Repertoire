package domain

import (
	"repertoire/domain/provider"
	"repertoire/domain/service"
	"repertoire/domain/usecase"

	"go.uber.org/fx"
)

var providers = fx.Options(
	fx.Provide(provider.NewCurrentUserProvider),
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
