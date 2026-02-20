package router

import (
	"repertoire/server/api/handler"
	"repertoire/server/api/server"
)

type SongArrangementRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.SongArrangementHandler
}

func (s SongArrangementRouter) RegisterRoutes() {
	api := s.requestHandler.PrivateRouter.Group("/songs/arrangements")
	{
		api.GET("", s.handler.GetAll)
		api.POST("", s.handler.Create)
		api.PUT("", s.handler.Update)
		api.PUT("/default", s.handler.UpdateDefault)
		api.PUT("/move", s.handler.Move)
		api.DELETE("/:id/from/:songID", s.handler.Delete)
	}
}

func NewSongArrangementRouter(
	requestHandler *server.RequestHandler,
	handler *handler.SongArrangementHandler,
) SongArrangementRouter {
	return SongArrangementRouter{
		handler:        handler,
		requestHandler: requestHandler,
	}
}
