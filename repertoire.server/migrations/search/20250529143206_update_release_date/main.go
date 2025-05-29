package main

import (
	"repertoire/server/data/database"
	"repertoire/server/data/logger"
	"repertoire/server/data/search"
	"repertoire/server/internal"
	"repertoire/server/internal/migration/utils"
	"repertoire/server/model"
)

var uid = "20250529143206"
var name = "update_release_date"

func main() {
	env := internal.NewEnv()
	log := logger.NewLogger(env)
	meiliClient := search.NewMeiliClient(env)
	dbClient := database.NewClient(logger.NewGormLogger(log), env)

	if utils.HasMigrationAlreadyBeenApplied(meiliClient, uid) {
		return
	}

	log.Info("Importing albums...")
	addAlbums(dbClient, meiliClient)
	log.Info("Albums imported!")

	log.Info("Importing songs...")
	addSongs(dbClient, meiliClient)
	log.Info("Songs imported!")

	utils.SaveMigrationStatus(meiliClient, uid, name)
}

func addAlbums(dbClient database.Client, meiliClient search.MeiliClient) {
	var albums []model.Album
	err := dbClient.Joins("Artist").Find(&albums).Error
	if err != nil {
		panic(err)
	}

	if len(albums) == 0 {
		return
	}

	var meiliAlbums []model.AlbumSearch
	for _, album := range albums {
		meiliAlbums = append(meiliAlbums, album.ToSearch())
	}
	_, err = meiliClient.Index("search").AddDocuments(meiliAlbums)
	if err != nil {
		panic(err)
	}
}

func addSongs(dbClient database.Client, meiliClient search.MeiliClient) {
	var songs []model.Song
	err := dbClient.Joins("Album").Joins("Artist").Find(&songs).Error
	if err != nil {
		panic(err)
	}

	if len(songs) == 0 {
		return
	}

	var meiliSongs []model.SongSearch
	for _, song := range songs {
		meiliSongs = append(meiliSongs, song.ToSearch())
	}
	_, err = meiliClient.Index("search").AddDocuments(meiliSongs)
	if err != nil {
		panic(err)
	}
}
