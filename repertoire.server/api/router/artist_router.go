package router

import (
	"repertoire/server/api/handler"
	"repertoire/server/api/server"
)

type ArtistRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.ArtistHandler
}

func (a ArtistRouter) RegisterRoutes() {
	api := a.requestHandler.PrivateRouter.Group("/artists")
	{
		api.GET("/:id", a.handler.Get)
		api.GET("", a.handler.GetAll)
		api.POST("", a.handler.Create)
		api.PUT("", a.handler.Update)
		api.DELETE("/:id", a.handler.Delete)
	}
}

func NewArtistRouter(
	requestHandler *server.RequestHandler,
	handler *handler.ArtistHandler,
) ArtistRouter {
	return ArtistRouter{
		handler:        handler,
		requestHandler: requestHandler,
	}
}
