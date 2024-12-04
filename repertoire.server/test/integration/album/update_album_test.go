package album

import (
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

func TestUpdateAlbum_WhenAlbumIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	request := requests.UpdateAlbumRequest{
		ID:    uuid.New(),
		Title: "New Title",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/albums", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateAlbum_WhenSuccessful_ShouldUpdateAlbum(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	album := albumData.Albums[0]

	request := requests.UpdateAlbumRequest{
		ID:          album.ID,
		Title:       "New Title",
		ReleaseDate: &[]time.Time{time.Now()}[0],
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/albums", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase()
	db.Find(&album, album.ID)

	assertUpdatedAlbum(t, request, album)
}

func assertUpdatedAlbum(t *testing.T, request requests.UpdateAlbumRequest, album model.Album) {
	assert.Equal(t, request.ID, album.ID)
	assert.Equal(t, request.Title, album.Title)
	assertion.Time(t, request.ReleaseDate, album.ReleaseDate)
}