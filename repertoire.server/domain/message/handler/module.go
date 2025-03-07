package handler

import (
	"go.uber.org/fx"
	"repertoire/server/domain/message/handler/album"
	"repertoire/server/domain/message/handler/artist"
	"repertoire/server/domain/message/handler/search"
	"repertoire/server/domain/message/handler/song"
)

var albumHandlers = fx.Options(
	fx.Provide(album.NewAlbumCreatedHandler),
	fx.Provide(album.NewAlbumDeletedHandler),
)

var artistHandlers = fx.Options(
	fx.Provide(artist.NewArtistCreatedHandler),
)

var songHandlers = fx.Options(
	fx.Provide(song.NewSongDeletedHandler),
)

var searchHandlers = fx.Options(
	fx.Provide(search.NewAddToSearchEngineHandler),
	fx.Provide(search.NewDeleteFromSearchEngineHandler),
)

var Module = fx.Options(
	albumHandlers,
	artistHandlers,
	songHandlers,
	searchHandlers,
)
