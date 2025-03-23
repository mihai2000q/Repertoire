package main

import (
	"github.com/meilisearch/meilisearch-go"
	"repertoire/server/data/database"
	"repertoire/server/data/search"
	"repertoire/server/internal"
	"repertoire/server/internal/migration/utils"
	"repertoire/server/model"
	"time"
)

var uid = "20250226172615"
var name = "initial_create"

func main() {
	env := internal.NewEnv()
	dbClient := database.NewClient(env)
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

	addArtists(dbClient, meiliClient)
	addAlbums(dbClient, meiliClient)
	addSongs(dbClient, meiliClient)
	addPlaylists(dbClient, meiliClient)

	utils.SaveMigrationStatus(meiliClient, uid, name)
}

func addArtists(dbClient database.Client, meiliClient search.Client) {
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
		artist.UpdatedAt = artist.UpdatedAt.UTC()
		artist.CreatedAt = artist.CreatedAt.UTC()
		meiliArtists = append(meiliArtists, artist.ToSearch())
	}
	_, err = meiliClient.Index("search").AddDocuments(meiliArtists)
	if err != nil {
		panic(err)
	}
}

func addAlbums(dbClient database.Client, meiliClient search.Client) {
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
		if album.ReleaseDate != nil {
			album.ReleaseDate = &[]time.Time{album.ReleaseDate.UTC()}[0]
		}
		album.UpdatedAt = album.UpdatedAt.UTC()
		album.CreatedAt = album.CreatedAt.UTC()
		meiliAlbums = append(meiliAlbums, album.ToSearch())
	}
	_, err = meiliClient.Index("search").AddDocuments(meiliAlbums)
	if err != nil {
		panic(err)
	}
}

func addSongs(dbClient database.Client, meiliClient search.Client) {
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
		if song.ReleaseDate != nil {
			song.ReleaseDate = &[]time.Time{song.ReleaseDate.UTC()}[0]
		}
		song.UpdatedAt = song.UpdatedAt.UTC()
		song.CreatedAt = song.CreatedAt.UTC()
		meiliSongs = append(meiliSongs, song.ToSearch())
	}
	_, err = meiliClient.Index("search").AddDocuments(meiliSongs)
	if err != nil {
		panic(err)
	}
}

func addPlaylists(dbClient database.Client, meiliClient search.Client) {
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
		playlist.UpdatedAt = playlist.UpdatedAt.UTC()
		playlist.CreatedAt = playlist.CreatedAt.UTC()
		meiliPlaylists = append(meiliPlaylists, playlist.ToSearch())
	}
	_, err = meiliClient.Index("search").AddDocuments(meiliPlaylists)
	if err != nil {
		panic(err)
	}
}
