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
		guitarTuningsApi.POST("", s.handler.CreateGuitarTuning)
		guitarTuningsApi.PUT("/move", s.handler.MoveGuitarTuning)
		guitarTuningsApi.DELETE("/:id", s.handler.DeleteGuitarTuning)
	}

	sectionsApi := api.Group("/sections")
	{
		sectionsApi.POST("", s.handler.CreateSection)
		sectionsApi.PUT("", s.handler.UpdateSection)
		sectionsApi.PUT("/move", s.handler.MoveSection)
		sectionsApi.DELETE("/:id/from/:songID", s.handler.DeleteSection)
	}

	sectionTypesApi := sectionsApi.Group("/types")
	{
		sectionTypesApi.GET("", s.handler.GetSectionTypes)
		sectionTypesApi.POST("", s.handler.CreateSectionType)
		sectionTypesApi.PUT("/move", s.handler.MoveSectionType)
		sectionTypesApi.DELETE("/:id", s.handler.DeleteSectionType)
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
