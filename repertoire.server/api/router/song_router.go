package router

import (
	"repertoire/api/handler"
	"repertoire/api/server"
)

type SongRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.SongHandler
}

func (s SongRouter) RegisterRoutes() {
	api := s.requestHandler.PrivateRouter.Group("/songs")
	{
		api.GET("/:id", s.handler.Get)
		api.GET("/", s.handler.GetAll)
		api.GET("/guitar-tunings", s.handler.GetGuitarTunings)
		api.GET("/section-types", s.handler.GetSectionTypes)
		api.POST("/", s.handler.Create)
		api.PUT("/", s.handler.Update)
		api.DELETE("/:id", s.handler.Delete)
	}
}

func NewSongRouter(
	requestHandler *server.RequestHandler,
	handler *handler.SongHandler,
) SongRouter {
	return SongRouter{
		handler:        handler,
		requestHandler: requestHandler,
	}
}
