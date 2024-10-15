package usecases

import (
	"go.uber.org/fx"
	"repertoire/domain/usecases/auth"
	"repertoire/domain/usecases/song"
	"repertoire/domain/usecases/user"
)

var authUseCases = fx.Options(
	fx.Provide(auth.NewRefresh),
	fx.Provide(auth.NewSignIn),
	fx.Provide(auth.NewSignUp),
)

var songUseCases = fx.Options(
	fx.Provide(song.NewGetSong),
	fx.Provide(song.NewGetAllSongs),
	fx.Provide(song.NewCreateSong),
	fx.Provide(song.NewUpdateSong),
	fx.Provide(song.NewDeleteSong),
)

var userUseCases = fx.Options(
	fx.Provide(user.NewGetUser),
)

var Module = fx.Options(
	authUseCases,
	songUseCases,
	userUseCases,
)
