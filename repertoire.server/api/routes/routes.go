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
	playlistRouter router.PlaylistRouter,
	searchRouter router.SearchRouter,
	songRouter router.SongRouter,
	userDataRouter router.UserDataRouter,
	userRouter router.UserRouter,
) *Routes {
	routes := &Routes{
		albumRouter,
		artistRouter,
		playlistRouter,
		searchRouter,
		songRouter,
		userDataRouter,
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
