package main

import (
	"go.uber.org/fx"
	"repertoire/storage/api"
	"repertoire/storage/data"
	"repertoire/storage/domain"
	"repertoire/storage/internal"
)

func main() {
	fx.New(
		internal.Module,
		data.Module,
		domain.Module,
		api.Module,
	).Run()
}
