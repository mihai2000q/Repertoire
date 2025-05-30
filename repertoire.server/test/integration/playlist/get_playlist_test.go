package playlist

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	playlistData "repertoire/server/test/integration/test/data/playlist"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestGetPlaylist_WhenUserIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/playlists/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetPlaylist_WhenSuccessful_ShouldReturnPlaylistWithSongsOrderedByTrackNo(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	playlist := playlistData.Playlists[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/playlists/"+playlist.ID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var responsePlaylist model.Playlist
	_ = json.Unmarshal(w.Body.Bytes(), &responsePlaylist)

	db := utils.GetDatabase(t)
	db.Preload("PlaylistSongs", func(db *gorm.DB) *gorm.DB {
		return db.
			Preload("Song").
			Preload("Song.Artist").
			Preload("Song.Album").
			Order("song_track_no")
	}).Find(&playlist, playlist.ID)

	assertion.ResponsePlaylist(t, playlist, responsePlaylist, true)
}

func TestGetPlaylist_WhenRequestHasSongsOrderBy_ShouldReturnPlaylistAndSongsOrderedByQueryParam(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	playlist := playlistData.Playlists[0]
	songsOrder := "title desc"

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/playlists/"+playlist.ID.String()+"?songsOrderBy="+songsOrder)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var responsePlaylist model.Playlist
	_ = json.Unmarshal(w.Body.Bytes(), &responsePlaylist)

	db := utils.GetDatabase(t)
	db.Preload("PlaylistSongs", func(db *gorm.DB) *gorm.DB {
		return db.
			Joins("JOIN songs ON songs.id = playlist_songs.song_id").
			Preload("Song").
			Preload("Song.Artist").
			Preload("Song.Album").
			Order(songsOrder)
	}).Find(&playlist, playlist.ID)

	assertion.ResponsePlaylist(t, playlist, responsePlaylist, true)
}
