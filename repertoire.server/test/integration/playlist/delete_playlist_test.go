package playlist

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	playlistData "repertoire/server/test/integration/test/data/playlist"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestDeletePlaylist_WhenPlaylistIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/playlists/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeletePlaylist_WhenSuccessful_ShouldDeletePlaylist(t *testing.T) {
	tests := []struct {
		name     string
		playlist model.Playlist
	}{
		{
			"Without Files",
			playlistData.Playlists[0],
		},
		{
			"With Image",
			playlistData.Playlists[1],
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().DELETE(w, "/api/playlists/"+test.playlist.ID.String())

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			db := utils.GetDatabase()

			var deletedPlaylist model.Playlist
			db.Find(&deletedPlaylist, test.playlist.ID)

			assert.Empty(t, deletedPlaylist)
		})
	}
}
