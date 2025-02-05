package router

import (
	"repertoire/server/api/handler"
	"repertoire/server/api/server"
)

type UserDataRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.UserDataHandler
}

func (u UserDataRouter) RegisterRoutes() {
	api := u.requestHandler.PrivateRouter.Group("/user-data")

	bandMemberRolesApi := api.Group("/band-member-roles")
	{
		bandMemberRolesApi.POST("", u.handler.CreateBandMemberRole)
		bandMemberRolesApi.PUT("/move", u.handler.MoveBandMemberRole)
		bandMemberRolesApi.DELETE("/:id", u.handler.DeleteBandMemberRole)
	}
	
	guitarTuningsApi := api.Group("/guitar-tunings")
	{
		guitarTuningsApi.POST("", u.handler.CreateGuitarTuning)
		guitarTuningsApi.PUT("/move", u.handler.MoveGuitarTuning)
		guitarTuningsApi.DELETE("/:id", u.handler.DeleteGuitarTuning)
	}

	songSectionTypesApi := api.Group("/song-section-types")
	{
		songSectionTypesApi.POST("", u.handler.CreateSectionType)
		songSectionTypesApi.PUT("/move", u.handler.MoveSectionType)
		songSectionTypesApi.DELETE("/:id", u.handler.DeleteSectionType)
	}
}

func NewUserDataRouter(
	requestHandler *server.RequestHandler,
	handler *handler.UserDataHandler,
) UserDataRouter {
	return UserDataRouter{
		handler:        handler,
		requestHandler: requestHandler,
	}
}
