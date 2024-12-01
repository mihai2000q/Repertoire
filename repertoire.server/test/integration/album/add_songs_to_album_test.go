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

func TestAddSongsToAlbum_WhenAlbumIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	request := requests.AddSongsToAlbumRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/albums/add-songs", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAddSongsToAlbum_WhenSongAlreadyHasAnAlbum_ShouldReturnBadRequestError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	request := requests.AddSongsToAlbumRequest{
		ID: albumData.Albums[1].ID,
		SongIDs: []uuid.UUID{
			albumData.Songs[0].ID,
			albumData.Albums[0].Songs[0].ID,
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/albums/add-songs", request)

	// then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddSongsToAlbum_WhenSongHasDifferentArtist_ShouldReturnBadRequestError(t *testing.T) {
	t.Skip("Test has proven to show limitations of this use case, so it will have to be rewritten")
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	request := requests.AddSongsToAlbumRequest{
		ID: albumData.Albums[0].ID,
		SongIDs: []uuid.UUID{
			albumData.Songs[0].ID,
			albumData.Songs[4].ID,
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/albums/add-songs", request)

	// then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddSongsToAlbum_WhenSuccessful_ShouldHaveSongsOnAlbum(t *testing.T) {
	t.Skip("Test has proven to show limitations of this use case, so it will have to be rewritten")

	tests := []struct {
		name    string
		request requests.AddSongsToAlbumRequest
	}{
		{
			"When the album has an artist, new songs will inherit",
			requests.AddSongsToAlbumRequest{
				ID: albumData.Albums[0].ID,
				SongIDs: []uuid.UUID{
					albumData.Songs[0].ID,
					albumData.Songs[1].ID,
					albumData.Songs[2].ID,
				},
			},
		},
		{
			"When they don't have an artist",
			requests.AddSongsToAlbumRequest{
				ID: albumData.Albums[1].ID,
				SongIDs: []uuid.UUID{
					albumData.Songs[0].ID,
					albumData.Songs[1].ID,
					albumData.Songs[2].ID,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().POST(w, "/api/albums/add-songs", tt.request)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
			assertSongsAddedToAlbum(t, tt.request)
		})
	}
}

func assertSongsAddedToAlbum(t *testing.T, request requests.AddSongsToAlbumRequest) {
	db := utils.GetDatabase()

	var album model.Album
	db.
		Preload("Songs", func(db *gorm.DB) *gorm.DB {
			return db.Order("songs.album_track_no")
		}).
		Find(&album, request.ID)

	assert.GreaterOrEqual(t, len(album.Songs), len(request.SongIDs))

	sizeDiff := len(album.Songs) - len(request.SongIDs)
	for i := sizeDiff; i < len(album.Songs); i++ {
		assert.Equal(t, request.SongIDs[i], album.Songs[i].ID)
	}

	for i, song := range album.Songs {
		assert.Equal(t, uint(i+1), *song.AlbumTrackNo)
		assert.Equal(t, album.ID, song.AlbumID)
		assert.Equal(t, album.ArtistID, song.ArtistID)
	}
}
