package router

import (
	"repertoire/server/api/handler"
	"repertoire/server/api/server"
)

type SongSectionRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.SongHandler
}

func (s SongSectionRouter) RegisterRoutes() {
	api := s.requestHandler.PrivateRouter.Group("/songs/sections")
	{
		api.POST("", s.handler.CreateSection)
		api.PUT("", s.handler.UpdateSection)
		api.PUT("/occurrences", s.handler.UpdateSectionsOccurrences)
		api.PUT("/partial-occurrences", s.handler.UpdateSectionsPartialOccurrences)
		api.PUT("/all", s.handler.UpdateAllSections)
		api.PUT("/move", s.handler.MoveSection)
		api.PUT("/bulk", s.handler.BulkDeleteSections)
		api.DELETE("/:id/from/:songID", s.handler.DeleteSection)
	}

	api.Group("/types").GET("", s.handler.GetSectionTypes)
}

func NewSongSectionRouter(
	requestHandler *server.RequestHandler,
	handler *handler.SongHandler,
) SongSectionRouter {
	return SongSectionRouter{
		handler:        handler,
		requestHandler: requestHandler,
	}
}
