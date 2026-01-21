package domain

import (
	"repertoire/storage/domain/service"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(service.NewJwtService),
)
