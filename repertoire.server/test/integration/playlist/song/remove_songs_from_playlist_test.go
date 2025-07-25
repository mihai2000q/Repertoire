package song

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
	"slices"
	"testing"
)

func TestRemoveSongsFromPlaylist_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.RemoveSongsFromPlaylistRequest{
		ID: playlistData.Playlists[0].ID,
		PlaylistSongIDs: []uuid.UUID{
			playlistData.PlaylistsSongs[3].ID,
			playlistData.PlaylistsSongs[1].ID,
			uuid.New(),
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/playlists/songs/remove", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRemoveSongsFromPlaylist_WhenSuccessful_ShouldDeleteSongsFromPlaylist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	playlist := playlistData.Playlists[0]
	oldSongsLength := len(slices.DeleteFunc(playlistData.PlaylistsSongs, func(playlistSong model.PlaylistSong) bool {
		return playlistSong.PlaylistID != playlist.ID
	}))

	request := requests.RemoveSongsFromPlaylistRequest{
		ID: playlist.ID,
		PlaylistSongIDs: []uuid.UUID{
			playlistData.PlaylistsSongs[3].ID,
			playlistData.PlaylistsSongs[1].ID,
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/playlists/songs/remove", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)
	db.Preload("PlaylistSongs", func(db *gorm.DB) *gorm.DB {
		return db.Order("song_track_no")
	}).Find(&playlist, playlist.ID)
	assertRemoveSongsFromPlaylist(t, request, playlist, oldSongsLength)
}

func assertRemoveSongsFromPlaylist(
	t *testing.T,
	request requests.RemoveSongsFromPlaylistRequest,
	playlist model.Playlist,
	oldSongsLength int,
) {
	assert.Equal(t, playlist.ID, request.ID)

	assert.Len(t, playlist.PlaylistSongs, oldSongsLength-len(request.PlaylistSongIDs))
	for i, playlistSong := range playlist.PlaylistSongs {
		assert.NotContains(t, request.PlaylistSongIDs, playlistSong.ID)
		assert.Equal(t, uint(i)+1, playlistSong.SongTrackNo)
	}
}
