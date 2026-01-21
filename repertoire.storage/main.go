package main

import (
	"repertoire/storage/api"
	"repertoire/storage/data"
	"repertoire/storage/domain"
	"repertoire/storage/internal"

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
