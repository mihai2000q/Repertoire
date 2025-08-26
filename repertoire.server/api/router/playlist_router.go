package router

import (
	"repertoire/server/api/handler"
	"repertoire/server/api/server"
)

type PlaylistRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.PlaylistHandler
}

func (p PlaylistRouter) RegisterRoutes() {
	api := p.requestHandler.PrivateRouter.Group("/playlists")
	{
		api.GET("/:id", p.handler.Get)
		api.GET("", p.handler.GetAll)
		api.GET("/filters-metadata", p.handler.GetFiltersMetadata)
		api.POST("", p.handler.Create)
		api.POST("/add-albums", p.handler.AddAlbums)
		api.POST("/add-artists", p.handler.AddArtists)
		api.POST("/perfect-rehearsals", p.handler.AddPerfectRehearsals)
		api.PUT("", p.handler.Update)
		api.PUT("/bulk-delete", p.handler.BulkDelete)
		api.DELETE("/:id", p.handler.Delete)
	}

	imagesApi := api.Group("/images")
	{
		imagesApi.PUT("", p.handler.SaveImage)
		imagesApi.DELETE("/:id", p.handler.DeleteImage)
	}

	songsApi := api.Group("/songs")
	{
		songsApi.GET("/:id", p.handler.GetSongs)
		songsApi.POST("/add", p.handler.AddSongs)
		songsApi.POST("/shuffle", p.handler.Shuffle)
		songsApi.PUT("/move", p.handler.MoveSong)
		songsApi.PUT("/remove", p.handler.RemoveSongs)
	}
}

func NewPlaylistRouter(
	requestHandler *server.RequestHandler,
	handler *handler.PlaylistHandler,
) PlaylistRouter {
	return PlaylistRouter{
		handler:        handler,
		requestHandler: requestHandler,
	}
}
