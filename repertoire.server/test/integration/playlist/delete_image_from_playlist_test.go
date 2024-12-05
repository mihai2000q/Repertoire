package playlist

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/test/integration/test/core"
	playlistData "repertoire/server/test/integration/test/data/playlist"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestDeleteImageFromPlaylist_WhenPlaylistIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/playlists/images/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteImageFromPlaylist_WhenPlaylistHasNoImage_ShouldReturnBadRequestError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	playlist := playlistData.Playlists[1]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/playlists/images/"+playlist.ID.String())

	// then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteImageFromPlaylist_WhenSuccessful_ShouldUpdatePlaylistAndDeleteImage(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	playlist := playlistData.Playlists[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/playlists/images/"+playlist.ID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase()
	db.Find(&playlist, playlist.ID)

	assert.Nil(t, playlist.ImageURL)
}
