package main

import (
	"go.uber.org/fx"
	"repertoire/api"
	"repertoire/config"
	"repertoire/data"
	"repertoire/domain"
)

func main() {
	fx.New(
		fx.Provide(config.NewEnv),
		data.Module,
		domain.Module,
		api.Module,
	).Run()
}
