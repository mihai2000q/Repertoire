package playlist

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	playlistData "repertoire/server/test/integration/test/data/playlist"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestAddSongsToPlaylist_WhenSuccessful_ShouldHaveSongsOnPlaylist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddSongsToPlaylistRequest{
		ID: playlistData.Playlists[1].ID,
		SongIDs: []uuid.UUID{
			playlistData.Songs[1].ID,
			playlistData.Songs[2].ID,
			playlistData.Songs[3].ID,
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/add-songs", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)
	assertSongsAddedToPlaylist(t, request)
}

func assertSongsAddedToPlaylist(t *testing.T, request requests.AddSongsToPlaylistRequest) {
	db := utils.GetDatabase(t)

	var playlist model.Playlist
	db.Preload("PlaylistSongs", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Song").Order("song_track_no")
	}).Find(&playlist, request.ID)

	assert.GreaterOrEqual(t, len(playlist.Songs), len(request.SongIDs))

	sizeDiff := len(playlist.Songs) - len(request.SongIDs)
	for i := 0; i < len(playlist.Songs)-sizeDiff; i++ {
		assert.Equal(t, request.SongIDs[i], playlist.Songs[i+sizeDiff].ID)
	}

	for i, song := range playlist.Songs {
		assert.Equal(t, uint(i+1), song.PlaylistTrackNo)
	}
}
