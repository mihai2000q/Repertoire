package internal

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewEnv),
	fx.Provide(NewCache),
	fx.Provide(NewRestyClient),
)
