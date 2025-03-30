package routes

import (
	"context"
	"repertoire/auth/api/router"

	"go.uber.org/fx"
)

type Routes []Route

type Route interface {
	RegisterRoutes()
}

func NewRoutes(
	lc fx.Lifecycle,
	mainRouter router.MainRouter,
	centrifugoRouter router.CentrifugoRouter,
) *Routes {
	routes := &Routes{
		mainRouter,
		centrifugoRouter,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			routes.setup()
			return nil
		},
	})

	return routes
}

func (r Routes) setup() {
	for _, route := range r {
		route.RegisterRoutes()
	}
}
