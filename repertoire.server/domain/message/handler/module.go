package handler

import (
	"go.uber.org/fx"
	"repertoire/server/domain/message/handler/album"
	"repertoire/server/domain/message/handler/artist"
	"repertoire/server/domain/message/handler/search"
var albumHandlers = fx.Options(
	fx.Provide(album.NewAlbumCreatedHandler),
	fx.Provide(album.NewAlbumDeletedHandler),
)

var artistHandlers = fx.Options(
	fx.Provide(artist.NewArtistCreatedHandler),
)

var searchHandlers = fx.Options(
	fx.Provide(search.NewAddToSearchEngineHandler),
	fx.Provide(search.NewDeleteFromSearchEngineHandler),
)

var Module = fx.Options(
	albumHandlers,
	artistHandlers,
	searchHandlers,
)
