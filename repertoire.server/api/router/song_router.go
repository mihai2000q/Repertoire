package router

import (
	"repertoire/server/api/handler"
	"repertoire/server/api/server"
)

type SongRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.SongHandler
}

func (s SongRouter) RegisterRoutes() {
	api := s.requestHandler.PrivateRouter.Group("/songs")
	{
		api.GET("/:id", s.handler.Get)
		api.GET("", s.handler.GetAll)
		api.GET("/guitar-tunings", s.handler.GetGuitarTunings)
		api.POST("", s.handler.Create)
		api.PUT("", s.handler.Update)
		api.DELETE("/:id", s.handler.Delete)
	}

	sectionsApi := api.Group("/sections")
	{
		sectionsApi.GET("/types", s.handler.GetSectionTypes)
		sectionsApi.POST("", s.handler.CreateSection)
		sectionsApi.PUT("", s.handler.UpdateSection)
		sectionsApi.PUT("/move", s.handler.MoveSection)
		sectionsApi.DELETE("/:id/from/:songID", s.handler.DeleteSection)
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
