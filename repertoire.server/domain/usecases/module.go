package usecases

import (
	"go.uber.org/fx"
	"repertoire/domain/usecases/auth"
)

var authUseCases = fx.Options(
	fx.Provide(auth.NewRefresh),
	fx.Provide(auth.NewSignIn),
	fx.Provide(auth.NewSignUp),
)

var Module = fx.Options(
	authUseCases,
)
