package playlist

import (
	"encoding/json"
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

func TestCreatePlaylist_WhenSuccessful_ShouldCreatePlaylist(t *testing.T) {
	tests := []struct {
		name    string
		request requests.CreatePlaylistRequest
	}{
		{
			"Minimal",
			requests.CreatePlaylistRequest{
				Title: "New Playlist",
			},
		},
		{
			"Maximal",
			requests.CreatePlaylistRequest{
				Title:       "New Playlist",
				Description: "Description of the playlist",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

			user := playlistData.Users[0]

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().
				WithUser(user).
				POST(w, "/api/playlists", tt.request)

			// then
			var response struct{ ID uuid.UUID }
			_ = json.Unmarshal(w.Body.Bytes(), &response)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.NotEmpty(t, response)
			assertCreatedPlaylist(t, tt.request, response.ID, user.ID)
		})
	}
}

func assertCreatedPlaylist(t *testing.T, request requests.CreatePlaylistRequest, playlistID uuid.UUID, userID uuid.UUID) {
	db := utils.GetDatabase()

	var playlist model.Playlist
	db.Find(&playlist, playlistID)

	assert.Equal(t, request.Title, playlist.Title)
	assert.Equal(t, request.Description, playlist.Description)
	assert.Nil(t, playlist.ImageURL)
	assert.Equal(t, userID, playlist.UserID)
}
