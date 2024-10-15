package router

import (
	"repertoire/api/handler"
	"repertoire/api/server"
)

type AlbumRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.AlbumHandler
}

func (s AlbumRouter) RegisterRoutes() {
	api := s.requestHandler.PrivateRouter.Group("/albums")
	{
		api.GET("/:id", s.handler.Get)
		api.GET("/", s.handler.GetAll)
		api.POST("/", s.handler.Create)
		api.PUT("/", s.handler.Update)
		api.DELETE("/:id", s.handler.Delete)
	}
}

func NewAlbumRouter(
	requestHandler *server.RequestHandler,
	handler *handler.AlbumHandler,
) AlbumRouter {
	return AlbumRouter{
		handler:        handler,
		requestHandler: requestHandler,
	}
}
