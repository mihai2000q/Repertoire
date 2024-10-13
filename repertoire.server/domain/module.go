package domain

import (
	"repertoire/domain/provider"
	"repertoire/domain/service"

	"go.uber.org/fx"
)

var providers = fx.Options(
	fx.Provide(provider.NewCurrentUserProvider),
)

var services = fx.Options(
	fx.Provide(service.NewAuthService),
	fx.Provide(service.NewSongService),
	fx.Provide(service.NewUserService),
)

var Module = fx.Options(
	providers,
	services,
)
