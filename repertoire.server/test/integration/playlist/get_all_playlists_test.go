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

	"github.com/stretchr/testify/assert"
)

func TestGetAllPlaylists_WhenSuccessful_ShouldReturnPlaylists(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/playlists")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var responsePlaylists []model.EnhancedPlaylist
	_ = json.Unmarshal(w.Body.Bytes(), &responsePlaylists)

	db := utils.GetDatabase(t)

	var playlists []model.Playlist
	db.Preload("Songs").Find(&playlists)

	for i := range responsePlaylists {
		assertion.ResponseEnhancedPlaylist(t, playlists[i], responsePlaylists[i])
	}
}
