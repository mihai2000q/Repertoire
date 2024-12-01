package artist

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	artistData "repertoire/server/test/integration/test/data/artist"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateArtist_WhenArtistIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	request := requests.UpdateArtistRequest{
		ID:   uuid.New(),
		Name: "New Name",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/artists", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateArtist_WhenSuccessful_ShouldUpdateArtist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artist := artistData.Artists[0]

	request := requests.UpdateArtistRequest{
		ID:   artist.ID,
		Name: "New Name",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/artists", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase()
	db.Find(&artist, artist.ID)

	assertUpdatedArtist(t, request, artist)
}

func assertUpdatedArtist(t *testing.T, request requests.UpdateArtistRequest, artist model.Artist) {
	assert.Equal(t, request.ID, artist.ID)
	assert.Equal(t, request.Name, artist.Name)
}
