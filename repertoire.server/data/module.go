package data

import (
	"go.uber.org/fx"
	"repertoire/data/database"
	"repertoire/data/repository"
	"repertoire/data/service"
)

var repositories = fx.Options(
	fx.Provide(repository.NewUserRepository),
)

var services = fx.Options(
	fx.Provide(service.NewJwtService),
)

var Module = fx.Options(
	fx.Provide(database.NewClient),
	repositories,
	services,
	fx.Invoke(func(database.Client) {}),
)
