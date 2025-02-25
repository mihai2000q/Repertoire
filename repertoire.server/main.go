package main

import (
	"go.uber.org/fx"
	"repertoire/server/api"
	"repertoire/server/data"
	"repertoire/server/domain"
	"repertoire/server/internal"
)

func main() {
	fx.New(
		internal.Module,
		data.Module,
		domain.Module,
		api.Module,
	).Run()
}
