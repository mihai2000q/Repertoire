package main

import (
	"repertoire/server/data/search"
	"repertoire/server/internal"
	"repertoire/server/internal/migration/utils"
)

var uid = "20250419104445"
var name = "add_id_filter"

func main() {
	env := internal.NewEnv()
	meiliClient := search.NewMeiliClient(env)

	if utils.HasMigrationAlreadyBeenApplied(meiliClient, uid) {
		return
	}

	_, err := meiliClient.Index("search").UpdateFilterableAttributes(&[]interface{}{
		"id", "type", "userId", "album", "album.id", "artist", "artist.id",
	})
	if err != nil {
		panic(err)
	}

	utils.SaveMigrationStatus(meiliClient, uid, name)
}
