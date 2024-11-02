package usecase

import (
	"repertoire/server/domain/usecase/album"
	"repertoire/server/domain/usecase/artist"
	"repertoire/server/domain/usecase/auth"
	"repertoire/server/domain/usecase/playlist"
	"repertoire/server/domain/usecase/song"
	"repertoire/server/domain/usecase/song/section"
	"repertoire/server/domain/usecase/user"

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
	fx.Provide(section.NewMoveSongSection),
	fx.Provide(section.NewUpdateSongSection),
	fx.Provide(section.NewDeleteSongSection),
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
