package router

import (
	"repertoire/api/handler"
	"repertoire/api/server"
)

type AlbumRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.AlbumHandler
}

func (a AlbumRouter) RegisterRoutes() {
	api := a.requestHandler.PrivateRouter.Group("/albums")
	{
		api.GET("/:id", a.handler.Get)
		api.GET("/", a.handler.GetAll)
		api.POST("/", a.handler.Create)
		api.PUT("/", a.handler.Update)
		api.DELETE("/:id", a.handler.Delete)
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
