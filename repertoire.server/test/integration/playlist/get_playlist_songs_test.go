package playlist

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	playlistData "repertoire/server/test/integration/test/data/playlist"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestGetPlaylistSongs_WhenSuccessful_ShouldReturnPlaylistWithSongsOrderedByTrackNo(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	playlist := playlistData.Playlists[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/playlists/songs/"+playlist.ID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var response wrapper.WithTotalCount[model.Song]
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	var playlistSongs []model.PlaylistSong
	db := utils.GetDatabase(t)
	db.Preload("Song").
		Preload("Song.Artist").
		Preload("Song.Album").
		Order("song_track_no").
		Find(&playlistSongs, model.PlaylistSong{PlaylistID: playlist.ID})

	for i := range playlistSongs {
		assertion.ResponsePlaylistSong(t, playlistSongs[i], response.Models[i])
	}
}

func TestGetPlaylist_WhenRequestHasSongsOrderBy_ShouldReturnPlaylistAndSongsOrderedByQueryParam(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	playlist := playlistData.Playlists[0]
	songsOrder := "\"Song\".title desc"

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/playlists/songs/"+playlist.ID.String()+"?orderBy="+songsOrder)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var response wrapper.WithTotalCount[model.Song]
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	var playlistSongs []model.PlaylistSong
	db := utils.GetDatabase(t)
	db.Joins("Song").
		Joins("Song.Artist").
		Joins("Song.Album").
		Order(songsOrder).
		Find(&playlistSongs, model.PlaylistSong{PlaylistID: playlist.ID})

	for i := range playlistSongs {
		playlistSongs[i].Song.ToFullImageURL()
		assertion.ResponsePlaylistSong(t, playlistSongs[i], response.Models[i])
	}
}
