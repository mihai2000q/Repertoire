package router

import (
	"repertoire/api/handler"
	"repertoire/api/server"
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
		api.PUT("", p.handler.Update)
		api.DELETE("/:id", p.handler.Delete)
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
