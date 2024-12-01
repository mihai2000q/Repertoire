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

	"github.com/stretchr/testify/assert"
)

func TestGetAllArtists_WhenSuccessful_ShouldReturnArtists(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/artists")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var responseArtists []model.Artist
	_ = json.Unmarshal(w.Body.Bytes(), &responseArtists)

	db := utils.GetDatabase()

	var artists []model.Artist
	db.Find(&artists)

	for i := range responseArtists {
		assertion.ResponseArtist(t, artists[i], responseArtists[i])
	}
}
