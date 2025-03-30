package domain

import (
	"go.uber.org/fx"
	"repertoire/auth/domain/service"
)

var Module = fx.Options(
	fx.Provide(service.NewCentrifugoService),
	fx.Provide(service.NewMainService),
	fx.Provide(service.NewStorageService),
)
