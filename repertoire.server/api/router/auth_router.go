package router

import (
	"repertoire/server/api/handler"
	"repertoire/server/api/server"
)

type AuthRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.AuthHandler
}

func (u AuthRouter) RegisterRoutes() {
	api := u.requestHandler.PublicRouter.Group("/auth")
	{
		api.PUT("/refresh", u.handler.Refresh)
		api.PUT("/sign-in", u.handler.SignIn)
		api.POST("/sign-up", u.handler.SignUp)
	}

	centrifugoApi := u.requestHandler.PrivateRouter.Group("/auth/centrifugo")
	{
		centrifugoApi.GET("/centrifugo", u.handler.GetCentrifugoToken)
	}
}

func NewAuthRouter(
	requestHandler *server.RequestHandler,
	handler *handler.AuthHandler,
) AuthRouter {
	return AuthRouter{
		handler:        handler,
		requestHandler: requestHandler,
	}
}
