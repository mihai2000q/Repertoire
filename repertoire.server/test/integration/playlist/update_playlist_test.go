package playlist

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	playlistData "repertoire/server/test/integration/test/data/playlist"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestUpdatePlaylist_WhenPlaylistIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.UpdatePlaylistRequest{
		ID:    uuid.New(),
		Title: "New Title",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/playlists", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdatePlaylist_WhenSuccessful_ShouldUpdatePlaylist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	playlist := playlistData.Playlists[0]

	request := requests.UpdatePlaylistRequest{
		ID:          playlist.ID,
		Title:       "New Title",
		Description: "New description",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/playlists", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase()
	db.Find(&playlist, playlist.ID)

	assertUpdatedPlaylist(t, request, playlist)
}

func assertUpdatedPlaylist(t *testing.T, request requests.UpdatePlaylistRequest, playlist model.Playlist) {
	assert.Equal(t, request.ID, playlist.ID)
	assert.Equal(t, request.Title, playlist.Title)
	assert.Equal(t, request.Description, playlist.Description)
}
