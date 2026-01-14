package router

import (
	"repertoire/server/api/handler"
	"repertoire/server/api/server"
)

type SongSectionRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.SongSectionHandler
}

func (s SongSectionRouter) RegisterRoutes() {
	api := s.requestHandler.PrivateRouter.Group("/songs/sections")
	{
		api.POST("", s.handler.Create)
		api.POST("bulk-rehearsals", s.handler.BulkRehearsals)
		api.PUT("", s.handler.Update)
		api.PUT("/occurrences", s.handler.UpdateOccurrences)
		api.PUT("/partial-occurrences", s.handler.UpdatePartialOccurrences)
		api.PUT("/all", s.handler.UpdateAll)
		api.PUT("/move", s.handler.Move)
		api.PUT("/bulk", s.handler.BulkDelete)
		api.DELETE("/:id/from/:songID", s.handler.Delete)
	}

	api.Group("/types").GET("", s.handler.GetTypes)
}

func NewSongSectionRouter(
	requestHandler *server.RequestHandler,
	handler *handler.SongSectionHandler,
) SongSectionRouter {
	return SongSectionRouter{
		handler:        handler,
		requestHandler: requestHandler,
	}
}
