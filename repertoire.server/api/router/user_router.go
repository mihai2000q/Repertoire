package router

import (
	"repertoire/server/api/handler"
	"repertoire/server/api/server"
)

type UserRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.UserHandler
}

func (u UserRouter) RegisterRoutes() {
	api := u.requestHandler.PrivateRouter.Group("/users")
	{
		api.GET("/current", u.handler.GetCurrentUser)
		api.GET("/:id", u.handler.Get)
		api.PUT("", u.handler.Update)
		api.PUT("/pictures", u.handler.SaveProfilePicture)
		api.DELETE("/pictures", u.handler.DeleteProfilePicture)
	}
}

func NewUserRouter(
	requestHandler *server.RequestHandler,
	handler *handler.UserHandler,
) UserRouter {
	return UserRouter{
		handler:        handler,
		requestHandler: requestHandler,
	}
}
