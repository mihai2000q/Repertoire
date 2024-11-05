package router

import (
	"repertoire/server/api/handler"
	"repertoire/server/api/server"
)

type AlbumRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.AlbumHandler
}

func (a AlbumRouter) RegisterRoutes() {
	api := a.requestHandler.PrivateRouter.Group("/albums")
	{
		api.GET("/:id", a.handler.Get)
		api.GET("", a.handler.GetAll)
		api.POST("", a.handler.Create)
		api.POST("/add-song", a.handler.AddSong)
		api.PUT("", a.handler.Update)
		api.PUT("/move-song", a.handler.MoveSong)
		api.DELETE("/song/:songID/from/:id", a.handler.RemoveSong)
		api.DELETE("/:id", a.handler.Delete)
	}

	imagesApi := api.Group("/images")
	{
		imagesApi.PUT("", a.handler.SaveImage)
		imagesApi.DELETE("/:id", a.handler.DeleteImage)
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
