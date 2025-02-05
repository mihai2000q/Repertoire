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
		api.POST("", s.handler.Create)
		api.POST("/perfect-rehearsal", s.handler.AddPerfectRehearsal)
		api.PUT("", s.handler.Update)
		api.DELETE("/:id", s.handler.Delete)
	}

	imagesApi := api.Group("/images")
	{
		imagesApi.PUT("", s.handler.SaveImage)
		imagesApi.DELETE("/:id", s.handler.DeleteImage)
	}

	guitarTuningsApi := api.Group("/guitar-tunings")
	{
		guitarTuningsApi.GET("", s.handler.GetGuitarTunings)
	}

	sectionsApi := api.Group("/sections")
	{
		sectionsApi.POST("", s.handler.CreateSection)
		sectionsApi.PUT("", s.handler.UpdateSection)
		sectionsApi.PUT("/occurrences", s.handler.UpdateSectionsOccurrences)
		sectionsApi.PUT("/move", s.handler.MoveSection)
		sectionsApi.DELETE("/:id/from/:songID", s.handler.DeleteSection)
	}

	sectionTypesApi := sectionsApi.Group("/types")
	{
		sectionTypesApi.GET("", s.handler.GetSectionTypes)
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
