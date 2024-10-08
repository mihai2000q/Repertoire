package data

import (
	"go.uber.org/fx"
	"repertoire/data/database"
	"repertoire/data/repository"
)

var repositories = fx.Options(
	fx.Provide(repository.NewUserRepository),
)

var Module = fx.Options(
	fx.Provide(database.NewClient),
	repositories,
	fx.Invoke(func(database.Client) {}),
)
