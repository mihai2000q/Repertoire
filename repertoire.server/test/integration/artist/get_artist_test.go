package artist

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	artistData "repertoire/server/test/integration/test/data/artist"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetArtist_WhenUserIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/artists/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetArtist_WhenSuccessful_ShouldReturnArtist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artist := artistData.Artists[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/artists/"+artist.ID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var responseArtist model.Artist
	_ = json.Unmarshal(w.Body.Bytes(), &responseArtist)

	db := utils.GetDatabase(t)
	db.Find(&artist, artist.ID)

	assertion.ResponseArtist(t, artist, responseArtist, true)
}
