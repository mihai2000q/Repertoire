package routes

import (
	"context"
	"repertoire/server/api/router"

	"go.uber.org/fx"
)

type Routes []Route

type Route interface {
	RegisterRoutes()
}

func NewRoutes(
	lc fx.Lifecycle,
	albumRouter router.AlbumRouter,
	artistRouter router.ArtistRouter,
	authRouter router.AuthRouter,
	playlistRouter router.PlaylistRouter,
	songRouter router.SongRouter,
	userRouter router.UserRouter,
) *Routes {
	routes := &Routes{
		albumRouter,
		artistRouter,
		authRouter,
		playlistRouter,
		songRouter,
		userRouter,
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
