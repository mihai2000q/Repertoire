package main

import (
	"repertoire/server/api"
	"repertoire/server/data"
	"repertoire/server/domain"
	"repertoire/server/internal"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		internal.Module,
		data.Module,
		domain.Module,
		api.Module,
	).Run()
}
