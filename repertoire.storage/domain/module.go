package domain

import (
	"go.uber.org/fx"
	"repertoire/storage/domain/service"
)

var Module = fx.Options(
	fx.Provide(service.NewJwtService),
)
