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
	publicApi := u.requestHandler.PublicRouter.Group("/user")
	{
		publicApi.POST("/sign-up", u.handler.SignUp)
	}

	api := u.requestHandler.PrivateRouter.Group("/users")
	{
		api.GET("/current", u.handler.GetCurrentUser)
		api.GET("/:id", u.handler.Get)
		api.PUT("", u.handler.Update)
		api.DELETE("", u.handler.Delete)
	}

	picturesApi := api.Group("/pictures")
	{
		picturesApi.PUT("", u.handler.SaveProfilePicture)
		picturesApi.DELETE("", u.handler.DeleteProfilePicture)
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
