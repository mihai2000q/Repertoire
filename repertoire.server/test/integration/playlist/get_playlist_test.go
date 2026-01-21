package playlist

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	playlistData "repertoire/server/test/integration/test/data/playlist"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

func TestGetPlaylist_WhenSuccessful_ShouldReturnPlaylist(t *testing.T) {
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
	db.Find(&playlist, playlist.ID)

	assertion.ResponsePlaylist(t, playlist, responsePlaylist)
}
