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

func TestMoveSongFromPlaylist_WhenPlaylistIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.MoveSongFromPlaylistRequest{
		ID:         uuid.New(),
		SongID:     uuid.New(),
		OverSongID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/playlists/move-song", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveSongFromPlaylist_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	playlist := playlistData.Playlists[0]
	request := requests.MoveSongFromPlaylistRequest{
		ID:         playlist.ID,
		SongID:     uuid.New(),
		OverSongID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/playlists/move-song", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveSongFromPlaylist_WhenOverSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	playlist := playlistData.Playlists[0]
	request := requests.MoveSongFromPlaylistRequest{
		ID:         playlist.ID,
		SongID:     playlistData.Songs[0].ID,
		OverSongID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/playlists/move-song", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveSongFromPlaylist_WhenSuccessful_ShouldMoveSongs(t *testing.T) {
	tests := []struct {
		name      string
		playlist  model.Playlist
		index     int
		overIndex int
	}{
		{
			"From upper position to lower",
			playlistData.Playlists[0],
			2,
			0,
		},
		{
			"From lower position to upper",
			playlistData.Playlists[0],
			0,
			2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

			request := requests.MoveSongFromPlaylistRequest{
				ID:         test.playlist.ID,
				SongID:     playlistData.Songs[test.index].ID,
				OverSongID: playlistData.Songs[test.overIndex].ID,
			}

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().PUT(w, "/api/playlists/move-song", request)

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			var playlist model.Playlist
			db := utils.GetDatabase()
			db.Preload("PlaylistSongs", func(db *gorm.DB) *gorm.DB {
				return db.Preload("Song").Order("song_track_no")
			}).Find(&playlist, request.ID)

			assertMovedSongs(t, request, playlist, test.index, test.overIndex)
		})
	}
}

func assertMovedSongs(
	t *testing.T,
	request requests.MoveSongFromPlaylistRequest,
	playlist model.Playlist,
	index int,
	overIndex int,
) {
	assert.Equal(t, request.ID, playlist.ID)

	if index < overIndex {
		assert.Equal(t, playlist.Songs[overIndex-1].ID, request.OverSongID)
	} else if index > overIndex {
		assert.Equal(t, playlist.Songs[overIndex+1].ID, request.OverSongID)
	}

	assert.Equal(t, playlist.Songs[overIndex].ID, request.SongID)
	for i, song := range playlist.Songs {
		assert.Equal(t, uint(i)+1, song.PlaylistTrackNo)
	}
}
