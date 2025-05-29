package main

import (
	"github.com/meilisearch/meilisearch-go"
	"repertoire/server/data/database"
	"repertoire/server/data/logger"
	"repertoire/server/data/search"
	"repertoire/server/internal"
	"repertoire/server/internal/migration/utils"
	"repertoire/server/model"
)

var uid = "20250226172615"
var name = "initial_create"

func main() {
	env := internal.NewEnv()
	log := logger.NewLogger(env)
	dbClient := database.NewClient(logger.NewGormLogger(log), env)
	meiliClient := search.NewMeiliClient(env)

	if utils.HasMigrationAlreadyBeenApplied(meiliClient, uid) {
		return
	}

	_, err := meiliClient.CreateIndex(&meilisearch.IndexConfig{
		Uid:        "search",
		PrimaryKey: "id",
	})
	if err != nil {
		panic(err)
	}

	_, err = meiliClient.Index("search").UpdateFilterableAttributes(&[]string{
		"type", "userId", "album", "album.id", "artist", "artist.id",
	})
	if err != nil {
		panic(err)
	}

	_, err = meiliClient.Index("search").UpdateSortableAttributes(&[]string{
		"title", "name", "updatedAt", "createdAt", "album", "album.title", "artist", "artist.name",
	})
	if err != nil {
		panic(err)
	}

	log.Info("Importing artists...")
	addArtists(dbClient, meiliClient)
	log.Info("Artists added!")

	log.Info("Importing albums...")
	addAlbums(dbClient, meiliClient)
	log.Info("Albums imported!")

	log.Info("Importing songs...")
	addSongs(dbClient, meiliClient)
	log.Info("Songs imported!")

	log.Info("Importing playlists...")
	addPlaylists(dbClient, meiliClient)
	log.Info("Playlists imported!")

	utils.SaveMigrationStatus(meiliClient, uid, name)
}

func addArtists(dbClient database.Client, meiliClient search.MeiliClient) {
	var artists []model.Artist
	err := dbClient.Find(&artists).Error
	if err != nil {
		panic(err)
	}

	if len(artists) == 0 {
		return
	}

	var meiliArtists []model.ArtistSearch
	for _, artist := range artists {
		meiliArtists = append(meiliArtists, artist.ToSearch())
	}
	_, err = meiliClient.Index("search").AddDocuments(meiliArtists)
	if err != nil {
		panic(err)
	}
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

func addPlaylists(dbClient database.Client, meiliClient search.MeiliClient) {
	var playlists []model.Playlist
	err := dbClient.Find(&playlists).Error
	if err != nil {
		panic(err)
	}

	if len(playlists) == 0 {
		return
	}

	var meiliPlaylists []model.PlaylistSearch
	for _, playlist := range playlists {
		meiliPlaylists = append(meiliPlaylists, playlist.ToSearch())
	}
	_, err = meiliClient.Index("search").AddDocuments(meiliPlaylists)
	if err != nil {
		panic(err)
	}
}
