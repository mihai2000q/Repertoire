package album

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	albumData "repertoire/server/test/integration/test/data/album"
	"repertoire/server/test/integration/test/utils"
	"testing"
	"time"
)

func TestCreateAlbum_WhenSuccessful_ShouldCreateAlbum(t *testing.T) {
	tests := []struct {
		name    string
		request requests.CreateAlbumRequest
	}{
		{
			"Minimal",
			requests.CreateAlbumRequest{
				Title:       "New Album",
				ReleaseDate: &[]time.Time{time.Now()}[0],
			},
		},
		{
			"With Artist",
			requests.CreateAlbumRequest{
				Title:    "New Album with Artist",
				ArtistID: &[]uuid.UUID{albumData.Artists[0].ID}[0],
			},
		},
		{
			"With New Artist",
			requests.CreateAlbumRequest{
				Title:      "New Album with new Artist",
				ArtistName: &[]string{"New Artist Name"}[0],
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

			user := albumData.Users[0]

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().
				WithUser(user).
				POST(w, "/api/albums", tt.request)

			// then
			var response struct{ ID uuid.UUID }
			_ = json.Unmarshal(w.Body.Bytes(), &response)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.NotEmpty(t, response)
			assertCreatedAlbum(t, tt.request, response.ID, user.ID)
		})
	}
}

func assertCreatedAlbum(t *testing.T, request requests.CreateAlbumRequest, albumID uuid.UUID, userID uuid.UUID) {
	db := utils.GetDatabase(t)

	var album model.Album
	db.Preload("Artist").Find(&album, albumID)

	assert.Equal(t, request.Title, album.Title)
	assertion.Time(t, request.ReleaseDate, album.ReleaseDate)
	assert.Nil(t, album.ImageURL)
	assert.Equal(t, userID, album.UserID)

	if request.ArtistID != nil {
		assert.Equal(t, request.ArtistID, album.ArtistID)
	}

	if request.ArtistName != nil {
		assert.NotEmpty(t, album.Artist.ID)
		assert.Equal(t, *request.ArtistName, album.Artist.Name)
		assert.Equal(t, userID, album.Artist.UserID)
	}
}
