package domain

import (
	"go.uber.org/fx"
	"repertoire/domain/service"
)

var services = fx.Options(
	fx.Provide(service.NewAuthService),
	fx.Provide(service.NewUserService),
)

var Module = fx.Options(
	services,
)
