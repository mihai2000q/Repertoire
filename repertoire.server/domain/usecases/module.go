package usecases

import (
	"repertoire/domain/usecases/album"
	"repertoire/domain/usecases/auth"
	"repertoire/domain/usecases/song"
	"repertoire/domain/usecases/user"

	"go.uber.org/fx"
)

var albumUseCases = fx.Options(
	fx.Provide(album.NewGetAlbum),
	fx.Provide(album.NewGetAllAlbums),
	fx.Provide(album.NewCreateAlbum),
	fx.Provide(album.NewUpdateAlbum),
	fx.Provide(album.NewDeleteAlbum),
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
	albumUseCases,
	authUseCases,
	songUseCases,
	userUseCases,
)
