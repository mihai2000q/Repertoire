package domain

import (
	"go.uber.org/fx"
	"repertoire/domain/provider"
	"repertoire/domain/service"
)

var providers = fx.Options(
	fx.Provide(provider.NewCurrentUserProvider),
)

var services = fx.Options(
	fx.Provide(service.NewAuthService),
	fx.Provide(service.NewUserService),
)

var Module = fx.Options(
	providers,
	services,
)
