package router

import (
	"repertoire/api/handler"
	"repertoire/api/server"
)

type ArtistRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.ArtistHandler
}

func (s ArtistRouter) RegisterRoutes() {
	api := s.requestHandler.PrivateRouter.Group("/artists")
	{
		api.GET("/:id", s.handler.Get)
		api.GET("/", s.handler.GetAll)
		api.POST("/", s.handler.Create)
		api.PUT("/", s.handler.Update)
		api.DELETE("/:id", s.handler.Delete)
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
