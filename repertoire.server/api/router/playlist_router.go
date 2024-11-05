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
		api.POST("", p.handler.Create)
		api.POST("/add-song", p.handler.AddSong)
		api.PUT("", p.handler.Update)
		api.DELETE("/:id", p.handler.Delete)
		api.DELETE("/song/:songID/from/:id", p.handler.RemoveSong)
	}

	imagesApi := api.Group("/images")
	{
		imagesApi.PUT("", p.handler.SaveImage)
		imagesApi.DELETE("/:id", p.handler.DeleteImage)
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
