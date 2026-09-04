package internal

import (
	"repertoire/server/internal/env"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(env.NewEnv),
)
