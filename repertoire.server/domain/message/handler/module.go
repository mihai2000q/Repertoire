package handler

import (
	"repertoire/server/domain/message/handler/album"
	"repertoire/server/domain/message/handler/artist"
	"repertoire/server/domain/message/handler/playlist"
	"repertoire/server/domain/message/handler/search"
	"repertoire/server/domain/message/handler/song"
	"repertoire/server/domain/message/handler/storage"
	"repertoire/server/domain/message/handler/user"

	"go.uber.org/fx"
)

var albumHandlers = fx.Options(
	fx.Provide(album.NewAlbumCreatedHandler),
	fx.Provide(album.NewAlbumsDeletedHandler),
	fx.Provide(album.NewAlbumsUpdatedHandler),
)

var artistHandlers = fx.Options(
	fx.Provide(artist.NewArtistCreatedHandler),
	fx.Provide(artist.NewArtistsDeletedHandler),
	fx.Provide(artist.NewArtistUpdatedHandler),
)

var playlistHandlers = fx.Options(
	fx.Provide(playlist.NewPlaylistCreatedHandler),
	fx.Provide(playlist.NewPlaylistsDeletedHandler),
	fx.Provide(playlist.NewPlaylistUpdatedHandler),
)

var songHandlers = fx.Options(
	fx.Provide(song.NewSongCreatedHandler),
	fx.Provide(song.NewSongsDeletedHandler),
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

var storageHandlers = fx.Options(
	fx.Provide(storage.NewDeleteDirectoriesStorageHandler),
)

var Module = fx.Options(
	albumHandlers,
	artistHandlers,
	playlistHandlers,
	songHandlers,
	userHandlers,
	searchHandlers,
	storageHandlers,
)
