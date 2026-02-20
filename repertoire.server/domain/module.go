package domain

import (
	"repertoire/server/domain/message"
	"repertoire/server/domain/processor"
	"repertoire/server/domain/provider"
	"repertoire/server/domain/service"
	"repertoire/server/domain/usecase"

	"go.uber.org/fx"
)

var processors = fx.Options(
	fx.Provide(processor.NewProgressProcessor),
	fx.Provide(processor.NewSongProcessor),
)

var providers = fx.Options(
	fx.Provide(provider.NewCurrentUserProvider),
	fx.Provide(provider.NewStorageFilePathProvider),
)

var services = fx.Options(
	fx.Provide(service.NewAlbumService),
	fx.Provide(service.NewArtistService),
	fx.Provide(service.NewPlaylistService),
	fx.Provide(service.NewSearchService),
	fx.Provide(service.NewSongArrangementService),
	fx.Provide(service.NewSongSectionService),
	fx.Provide(service.NewSongService),
	fx.Provide(service.NewUserDataService),
	fx.Provide(service.NewUserService),
)

var Module = fx.Options(
	processors,
	providers,
	message.Module,
	usecase.Module,
	services,
)
