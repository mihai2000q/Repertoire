package main

import (
	"fmt"
	"github.com/meilisearch/meilisearch-go"
	"repertoire/server/data/database"
	"repertoire/server/data/search"
	"repertoire/server/internal"
	"repertoire/server/internal/enums"
	"repertoire/server/model"
	"time"
)

func main() {
	env := internal.NewEnv()
	dbClient := database.NewClient(env)
	meiliClient := search.NewMeiliClient(env)

	_, err := meiliClient.CreateIndex(&meilisearch.IndexConfig{
		Uid:        "search",
		PrimaryKey: "id",
	})
	if err != nil {
		panic(err)
	}

	_, err = meiliClient.Index("search").UpdateFilterableAttributes(&[]string{"type", "userId"})
	if err != nil {
		panic(err)
	}

	addArtists(dbClient, meiliClient)
	addAlbums(dbClient, meiliClient)
	addSongs(dbClient, meiliClient)
	addPlaylists(dbClient, meiliClient)

	fmt.Println(time.Now().Format("YYYY/MM/DD hh/mm/ss") + " OK 20250226172615_initial_create imported!")
}

func addArtists(dbClient database.Client, meiliClient meilisearch.ServiceManager) {
	var artists []model.Artist
	err := dbClient.DB.Find(&artists).Error
	if err != nil {
		panic(err)
	}

	if len(artists) == 0 {
		return
	}

	var meiliArtists []model.ArtistSearch
	for _, artist := range artists {
		meiliArtist := model.ArtistSearch{
			ImageUrl:  artist.ImageURL.StripURL(),
			Name:      artist.Name,
			UpdatedAt: artist.UpdatedAt,
			SearchBase: model.SearchBase{
				ID:     "artist-" + artist.ID.String(),
				Type:   enums.Artist,
				UserID: artist.UserID,
			},
		}

		meiliArtists = append(meiliArtists, meiliArtist)
	}
	_, err = meiliClient.Index("search").AddDocuments(meiliArtists)
	if err != nil {
		panic(err)
	}
}

func addAlbums(dbClient database.Client, meiliClient meilisearch.ServiceManager) {
	var albums []model.Album
	err := dbClient.DB.Joins("Artist").Find(&albums).Error
	if err != nil {
		panic(err)
	}

	if len(albums) == 0 {
		return
	}

	var meiliAlbums []model.AlbumSearch
	for _, album := range albums {
		meiliAlbum := model.AlbumSearch{
			ImageUrl:  album.ImageURL.StripURL(),
			Title:     album.Title,
			UpdatedAt: album.UpdatedAt,
			SearchBase: model.SearchBase{
				ID:     "album-" + album.ID.String(),
				Type:   enums.Album,
				UserID: album.UserID,
			},
		}

		if album.Artist != nil {
			meiliAlbum.Artist = &model.AlbumArtistSearch{
				ID:       album.Artist.ID,
				Name:     album.Artist.Name,
				ImageUrl: album.Artist.ImageURL.StripURL(),
			}
		}

		meiliAlbums = append(meiliAlbums, meiliAlbum)
	}
	_, err = meiliClient.Index("search").AddDocuments(meiliAlbums)
	if err != nil {
		panic(err)
	}
}

func addSongs(dbClient database.Client, meiliClient meilisearch.ServiceManager) {
	var songs []model.Song
	err := dbClient.DB.Joins("Album").Joins("Artist").Find(&songs).Error
	if err != nil {
		panic(err)
	}

	if len(songs) == 0 {
		return
	}

	var meiliSongs []model.SongSearch
	for _, song := range songs {
		meiliSong := model.SongSearch{
			ImageUrl:  song.ImageURL.StripURL(),
			Title:     song.Title,
			UpdatedAt: song.UpdatedAt,
			SearchBase: model.SearchBase{
				ID:     "song-" + song.ID.String(),
				Type:   enums.Song,
				UserID: song.UserID,
			},
		}

		if song.Artist != nil {
			meiliSong.Artist = &model.SongArtistSearch{
				ID:       song.Artist.ID,
				Name:     song.Artist.Name,
				ImageUrl: song.Artist.ImageURL.StripURL(),
			}
		}

		if song.Album != nil {
			meiliSong.Album = &model.SongAlbumSearch{
				ID:       song.Album.ID,
				Title:    song.Album.Title,
				ImageUrl: song.Album.ImageURL.StripURL(),
			}
		}

		meiliSongs = append(meiliSongs, meiliSong)
	}
	_, err = meiliClient.Index("search").AddDocuments(meiliSongs)
	if err != nil {
		panic(err)
	}
}

func addPlaylists(dbClient database.Client, meiliClient meilisearch.ServiceManager) {
	var playlists []model.Playlist
	err := dbClient.DB.Find(&playlists).Error
	if err != nil {
		panic(err)
	}

	if len(playlists) == 0 {
		return
	}

	var meiliPlaylists []model.PlaylistSearch
	for _, playlist := range playlists {
		meiliPlaylist := model.PlaylistSearch{
			ImageUrl:  playlist.ImageURL.StripURL(),
			Title:     playlist.Title,
			UpdatedAt: playlist.UpdatedAt,
			SearchBase: model.SearchBase{
				ID:     "playlist-" + playlist.ID.String(),
				Type:   enums.Playlist,
				UserID: playlist.UserID,
			},
		}

		meiliPlaylists = append(meiliPlaylists, meiliPlaylist)
	}
	_, err = meiliClient.Index("search").AddDocuments(meiliPlaylists)
	if err != nil {
		panic(err)
	}
}
