package usecases

import (
	"go.uber.org/fx"
	"repertoire/domain/usecases/auth"
	"repertoire/domain/usecases/user"
)

var authUseCases = fx.Options(
	fx.Provide(auth.NewRefresh),
	fx.Provide(auth.NewSignIn),
	fx.Provide(auth.NewSignUp),
)

var userUseCases = fx.Options(
	fx.Provide(user.NewGetUser),
)

var Module = fx.Options(
	authUseCases,
	userUseCases,
)
