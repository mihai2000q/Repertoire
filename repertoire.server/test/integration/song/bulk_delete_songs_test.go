package song

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBulkDeleteSongs_WhenSongsAreNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.BulkDeleteSongsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/bulk-delete", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBulkDeleteSongs_WhenSuccessful_ShouldDeleteSongs(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	messages := utils.SubscribeToTopic(topics.SongsDeletedTopic)

	songs := []model.Song{
		songData.Songs[0],
		songData.Songs[1],
		songData.Songs[4],
		songData.Songs[6],
		songData.Songs[9],
	}
	request := requests.BulkDeleteSongsRequest{
		IDs: []uuid.UUID{
			songData.Songs[0].ID,
			songData.Songs[1].ID,
			songData.Songs[4].ID,
			songData.Songs[6].ID,
			songData.Songs[9].ID,
		},
	}

	db := utils.GetDatabase(t)
	db.Preload("Playlists").Find(&songs, request.IDs)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/bulk-delete", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db = db.Session(&gorm.Session{NewDB: true})

	var deletedSongs []model.Song
	db.Find(&deletedSongs, request.IDs)
	assert.Empty(t, deletedSongs)

	assertOrderedAlbums(t, songs, request, db)
	assertOrderedPlaylists(t, songs, request, db)

	assertion.AssertMessage(t, messages, func(payloadSongs []model.Song) {
		assert.Len(t, payloadSongs, len(songs))
		for i := range payloadSongs {
			assert.Contains(t, request.IDs, payloadSongs[i].ID)
		}
	})
}

func assertOrderedAlbums(t *testing.T, songs []model.Song, request requests.BulkDeleteSongsRequest, db *gorm.DB) {
	var albumIDs []uuid.UUID
	for _, song := range songs {
		if song.AlbumID != nil && !slices.Contains(albumIDs, *song.AlbumID) {
			albumIDs = append(albumIDs, *song.AlbumID)
		}
	}

	var albums []model.Album
	db.
		Preload("Songs", func(db *gorm.DB) *gorm.DB {
			return db.Order("album_track_no")
		}).
		Find(&albums, albumIDs)

	for _, album := range albums {
		for i, s := range album.Songs {
			assert.NotContains(t, request.IDs, s.ID)
			assert.Equal(t, uint(i+1), *s.AlbumTrackNo)
		}
	}
}

func assertOrderedPlaylists(t *testing.T, songs []model.Song, request requests.BulkDeleteSongsRequest, db *gorm.DB) {
	var playlistsIds []uuid.UUID
	for _, song := range songs {
		for _, playlist := range song.Playlists {
			if !slices.Contains(playlistsIds, playlist.ID) {
				playlistsIds = append(playlistsIds, playlist.ID)
			}
		}
	}

	var playlists []model.Playlist
	db.
		Preload("PlaylistSongs", func(db *gorm.DB) *gorm.DB {
			return db.Order("song_track_no")
		}).
		Find(&playlists, playlistsIds)

	for _, playlist := range playlists {
		for i, playlistSong := range playlist.PlaylistSongs {
			assert.NotContains(t, request.IDs, playlistSong.SongID)
			assert.Equal(t, uint(i+1), playlistSong.SongTrackNo)
		}
	}
}
