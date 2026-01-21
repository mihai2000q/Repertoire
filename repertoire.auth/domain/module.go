package domain

import (
	"repertoire/auth/domain/service"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(service.NewCentrifugoService),
	fx.Provide(service.NewMainService),
	fx.Provide(service.NewStorageService),
)
