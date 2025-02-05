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
	"slices"
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
	tests := []struct {
		name    string
		request requests.AddSongsToAlbumRequest
	}{
		{
			"when song has different artist",
			requests.AddSongsToAlbumRequest{
				ID: albumData.Albums[0].ID,
				SongIDs: []uuid.UUID{
					albumData.Songs[3].ID, // same artist
					albumData.Songs[4].ID,
				},
			},
		},
		{
			"when song has artist, but album doesn't",
			requests.AddSongsToAlbumRequest{
				ID: albumData.Albums[1].ID,
				SongIDs: []uuid.UUID{
					albumData.Songs[0].ID, // no artist
					albumData.Songs[3].ID,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().POST(w, "/api/albums/add-songs", test.request)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestAddSongsToAlbum_WhenSuccessful_ShouldHaveSongsOnAlbum(t *testing.T) {
	tests := []struct {
		name    string
		request requests.AddSongsToAlbumRequest
	}{
		{
			"Normal Add",
			requests.AddSongsToAlbumRequest{
				ID: albumData.Albums[1].ID,
				SongIDs: []uuid.UUID{
					albumData.Songs[0].ID,
					albumData.Songs[1].ID,
				},
			},
		},
		{
			"When the album has an artist and release date, new songs will inherit",
			requests.AddSongsToAlbumRequest{
				ID: albumData.Albums[0].ID,
				SongIDs: []uuid.UUID{
					albumData.Songs[0].ID,
					albumData.Songs[1].ID,
					albumData.Songs[2].ID,
					albumData.Songs[3].ID,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().POST(w, "/api/albums/add-songs", test.request)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
			assertSongsAddedToAlbum(t, test.request)
		})
	}
}

func assertSongsAddedToAlbum(t *testing.T, request requests.AddSongsToAlbumRequest) {
	db := utils.GetDatabase(t)

	var album model.Album
	db.Preload("Songs", func(db *gorm.DB) *gorm.DB {
		return db.Order("songs.album_track_no")
	}).Find(&album, request.ID)

	assert.GreaterOrEqual(t, len(album.Songs), len(request.SongIDs))

	sizeDiff := len(album.Songs) - len(request.SongIDs)
	for i := 0; i < len(album.Songs)-sizeDiff; i++ {
		assert.Equal(t, request.SongIDs[i], album.Songs[i+sizeDiff].ID)
	}

	for i, song := range album.Songs {
		assert.Equal(t, uint(i+1), *song.AlbumTrackNo)
		assert.Equal(t, album.ID, *song.AlbumID)
		assert.Equal(t, album.ArtistID, song.ArtistID)
		if album.ReleaseDate != nil && slices.Contains(request.SongIDs, song.ID) {
			assert.NotNil(t, song.ReleaseDate)
		}
	}
}
