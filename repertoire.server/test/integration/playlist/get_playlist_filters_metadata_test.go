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

func TestGetPlaylistFiltersMetadata_WhenSuccessful_ShouldReturnPlaylistFiltersMetadata(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/playlists/filters-metadata")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var metadata model.PlaylistFiltersMetadata
	_ = json.Unmarshal(w.Body.Bytes(), &metadata)

	var playlists []model.Playlist
	db := utils.GetDatabase(t)
	db.Preload("Songs").Find(&playlists)

	assertion.PlaylistFiltersMetadata(t, metadata, playlists)
}
