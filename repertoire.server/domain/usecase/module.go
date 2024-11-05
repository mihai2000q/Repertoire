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
	fx.Provide(album.NewAddSongToAlbum),
	fx.Provide(album.NewCreateAlbum),
	fx.Provide(album.NewDeleteAlbum),
	fx.Provide(album.NewGetAlbum),
	fx.Provide(album.NewGetAllAlbums),
	fx.Provide(album.NewMoveSongFromAlbum),
	fx.Provide(album.NewRemoveSongFromAlbum),
	fx.Provide(album.NewSaveImageToAlbum),
	fx.Provide(album.NewUpdateAlbum),
)

var artistUseCases = fx.Options(
	fx.Provide(artist.NewCreateArtist),
	fx.Provide(artist.NewDeleteArtist),
	fx.Provide(artist.NewGetArtist),
	fx.Provide(artist.NewGetAllArtists),
	fx.Provide(artist.NewUpdateArtist),
)

var authUseCases = fx.Options(
	fx.Provide(auth.NewRefresh),
	fx.Provide(auth.NewSignIn),
	fx.Provide(auth.NewSignUp),
)

var playlistUseCases = fx.Options(
	fx.Provide(playlist.NewAddSongToPlaylist),
	fx.Provide(playlist.NewCreatePlaylist),
	fx.Provide(playlist.NewDeletePlaylist),
	fx.Provide(playlist.NewGetPlaylist),
	fx.Provide(playlist.NewGetAllPlaylists),
	fx.Provide(playlist.NewRemoveSongFromPlaylist),
	fx.Provide(playlist.NewUpdatePlaylist),
)

var songUseCases = fx.Options(
	fx.Provide(song.NewCreateSong),
	fx.Provide(song.NewDeleteSong),
	fx.Provide(song.NewDeleteImageFromSong),
	fx.Provide(song.NewGetAllSongs),
	fx.Provide(song.NewGetGuitarTunings),
	fx.Provide(song.NewGetSong),
	fx.Provide(song.NewSaveImageToSong),
	fx.Provide(song.NewUpdateSong),

	fx.Provide(section.NewCreateSongSection),
	fx.Provide(section.NewDeleteSongSection),
	fx.Provide(section.NewGetSongSectionTypes),
	fx.Provide(section.NewMoveSongSection),
	fx.Provide(section.NewUpdateSongSection),
)

var userUseCases = fx.Options(
	fx.Provide(user.NewGetUser),
	fx.Provide(user.NewUpdateUser),
)

var Module = fx.Options(
	albumUseCases,
	artistUseCases,
	authUseCases,
	playlistUseCases,
	songUseCases,
	userUseCases,
)
