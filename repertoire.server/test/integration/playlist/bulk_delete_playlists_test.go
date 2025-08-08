package playlist

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	playlistData "repertoire/server/test/integration/test/data/playlist"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBulkDeletePlaylists_WhenPlaylistsAreNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.BulkDeletePlaylistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/playlists/bulk-delete", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBulkDeletePlaylists_WhenSuccessful_ShouldDeletePlaylists(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	messages := utils.SubscribeToTopic(topics.PlaylistsDeletedTopic)

	request := requests.BulkDeletePlaylistsRequest{
		IDs: []uuid.UUID{
			playlistData.Playlists[0].ID,
			playlistData.Playlists[1].ID,
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/playlists/bulk-delete", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var deletedPlaylists []model.Playlist
	db.Find(&deletedPlaylists, request.IDs)
	assert.Empty(t, deletedPlaylists)

	assertion.AssertMessage(t, messages, func(payloadPlaylists []model.Playlist) {
		assert.Len(t, payloadPlaylists, len(request.IDs))
		for i := range payloadPlaylists {
			assert.Contains(t, request.IDs, payloadPlaylists[i].ID)
		}
	})
}
