package router

import (
	"repertoire/server/api/handler"
	"repertoire/server/api/server"
)

type SongPartRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.SongPartHandler
}

func (s SongPartRouter) RegisterRoutes() {
	api := s.requestHandler.PrivateRouter.Group("/songs/parts")
	{
		api.POST("", s.handler.Create)
		api.POST("bulk-rehearsals", s.handler.BulkRehearsals)
		api.PUT("", s.handler.Update)
		api.PUT("/all", s.handler.UpdateAll)
		api.PUT("/move-in-song", s.handler.MoveInSong)
		api.PUT("/bulk-delete", s.handler.BulkDelete)
		api.DELETE("/:id/from/:songID", s.handler.Delete)
		api.DELETE("/:id/from/:songID/and/:sectionID", s.handler.Delete)
	}
}

func NewSongPartRouter(
	requestHandler *server.RequestHandler,
	handler *handler.SongPartHandler,
) SongPartRouter {
	return SongPartRouter{
		handler:        handler,
		requestHandler: requestHandler,
	}
}
