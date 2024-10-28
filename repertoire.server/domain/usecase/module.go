package usecase

import (
	"repertoire/domain/usecase/album"
	"repertoire/domain/usecase/artist"
	"repertoire/domain/usecase/auth"
	"repertoire/domain/usecase/playlist"
	"repertoire/domain/usecase/song"
	"repertoire/domain/usecase/song/section"
	"repertoire/domain/usecase/user"

	"go.uber.org/fx"
)

var albumUseCases = fx.Options(
	fx.Provide(album.NewGetAlbum),
	fx.Provide(album.NewGetAllAlbums),
	fx.Provide(album.NewCreateAlbum),
	fx.Provide(album.NewUpdateAlbum),
	fx.Provide(album.NewDeleteAlbum),
)

var artistUseCases = fx.Options(
	fx.Provide(artist.NewGetArtist),
	fx.Provide(artist.NewGetAllArtists),
	fx.Provide(artist.NewCreateArtist),
	fx.Provide(artist.NewUpdateArtist),
	fx.Provide(artist.NewDeleteArtist),
)

var authUseCases = fx.Options(
	fx.Provide(auth.NewRefresh),
	fx.Provide(auth.NewSignIn),
	fx.Provide(auth.NewSignUp),
)

var playlistUseCases = fx.Options(
	fx.Provide(playlist.NewGetPlaylist),
	fx.Provide(playlist.NewGetAllPlaylists),
	fx.Provide(playlist.NewCreatePlaylist),
	fx.Provide(playlist.NewUpdatePlaylist),
	fx.Provide(playlist.NewDeletePlaylist),
)

var songUseCases = fx.Options(
	fx.Provide(song.NewGetSong),
	fx.Provide(song.NewGetAllSongs),
	fx.Provide(song.NewGetGuitarTunings),
	fx.Provide(section.NewCreateSongSection),
	fx.Provide(section.NewGetSongSectionTypes),
	fx.Provide(song.NewCreateSong),
	fx.Provide(song.NewUpdateSong),
	fx.Provide(song.NewDeleteSong),
)

var userUseCases = fx.Options(
	fx.Provide(user.NewGetUser),
)

var Module = fx.Options(
	albumUseCases,
	artistUseCases,
	authUseCases,
	playlistUseCases,
	songUseCases,
	userUseCases,
)
