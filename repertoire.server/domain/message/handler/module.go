package handler

import (
	"go.uber.org/fx"
	"repertoire/server/domain/message/handler/album"
	"repertoire/server/domain/message/handler/artist"
	"repertoire/server/domain/message/handler/playlist"
	"repertoire/server/domain/message/handler/search"
	"repertoire/server/domain/message/handler/song"
	"repertoire/server/domain/message/handler/user"
)

var albumHandlers = fx.Options(
	fx.Provide(album.NewAlbumCreatedHandler),
	fx.Provide(album.NewAlbumDeletedHandler),
	fx.Provide(album.NewAlbumsUpdatedHandler),
)

var artistHandlers = fx.Options(
	fx.Provide(artist.NewArtistCreatedHandler),
	fx.Provide(artist.NewArtistDeletedHandler),
	fx.Provide(artist.NewArtistUpdatedHandler),
)

var playlistHandlers = fx.Options(
	fx.Provide(playlist.NewPlaylistCreatedHandler),
	fx.Provide(playlist.NewPlaylistDeletedHandler),
	fx.Provide(playlist.NewPlaylistUpdatedHandler),
)

var songHandlers = fx.Options(
	fx.Provide(song.NewSongCreatedHandler),
	fx.Provide(song.NewSongDeletedHandler),
	fx.Provide(song.NewSongsUpdatedHandler),
)

var userHandlers = fx.Options(
	fx.Provide(user.NewUserDeletedHandler),
)

var searchHandlers = fx.Options(
	fx.Provide(search.NewAddToSearchEngineHandler),
	fx.Provide(search.NewDeleteFromSearchEngineHandler),
	fx.Provide(search.NewUpdateFromSearchEngineHandler),
)

var Module = fx.Options(
	albumHandlers,
	artistHandlers,
	playlistHandlers,
	songHandlers,
	searchHandlers,
	userHandlers,
)
