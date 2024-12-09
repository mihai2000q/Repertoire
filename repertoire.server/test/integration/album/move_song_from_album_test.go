package album

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	albumData "repertoire/server/test/integration/test/data/album"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestMoveSongFromAlbum_WhenAlbumIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	request := requests.MoveSongFromAlbumRequest{
		ID:         uuid.New(),
		SongID:     uuid.New(),
		OverSongID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/albums/move-song", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveSongFromAlbum_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	album := albumData.Albums[0]
	request := requests.MoveSongFromAlbumRequest{
		ID:         album.ID,
		SongID:     uuid.New(),
		OverSongID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/albums/move-song", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveSongFromAlbum_WhenOverSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	album := albumData.Albums[0]
	request := requests.MoveSongFromAlbumRequest{
		ID:         album.ID,
		SongID:     album.Songs[0].ID,
		OverSongID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/albums/move-song", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveSongFromAlbum_WhenSuccessful_ShouldMoveSongs(t *testing.T) {
	tests := []struct {
		name      string
		album     model.Album
		index     int
		overIndex int
	}{
		{
			"From upper position to lower",
			albumData.Albums[0],
			2,
			0,
		},
		{
			"From lower position to upper",
			albumData.Albums[0],
			0,
			2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

			request := requests.MoveSongFromAlbumRequest{
				ID:         test.album.ID,
				SongID:     test.album.Songs[test.index].ID,
				OverSongID: test.album.Songs[test.overIndex].ID,
			}

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().PUT(w, "/api/albums/move-song", request)

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			var album model.Album
			db := utils.GetDatabase()
			db.Preload("Songs", func(db *gorm.DB) *gorm.DB {
				return db.Order("songs.album_track_no")
			}).Find(&album, request.ID)

			assertMovedSongs(t, request, album, test.index, test.overIndex)
		})
	}
}

func assertMovedSongs(t *testing.T, request requests.MoveSongFromAlbumRequest, album model.Album, index int, overIndex int) {
	assert.Equal(t, request.ID, album.ID)

	if index < overIndex {
		assert.Equal(t, album.Songs[overIndex-1].ID, request.OverSongID)
	} else if index > overIndex {
		assert.Equal(t, album.Songs[overIndex+1].ID, request.OverSongID)
	}

	assert.Equal(t, album.Songs[overIndex].ID, request.SongID)
	for i, song := range album.Songs {
		assert.Equal(t, uint(i)+1, *song.AlbumTrackNo)
	}
}
