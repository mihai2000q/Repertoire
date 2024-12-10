package artist

import (
	"encoding/json"
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

func TestCreateArtist_WhenSuccessful_ShouldCreateArtist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	user := artistData.Users[0]

	request := requests.CreateArtistRequest{
		Name: "New Artist",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithUser(user).
		POST(w, "/api/artists", request)

	// then
	var response struct{ ID uuid.UUID }
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, response)
	assertCreatedArtist(t, request, response.ID, user.ID)
}

func assertCreatedArtist(t *testing.T, request requests.CreateArtistRequest, artistID uuid.UUID, userID uuid.UUID) {
	db := utils.GetDatabase(t)

	var artist model.Artist
	db.Find(&artist, artistID)

	assert.Equal(t, request.Name, artist.Name)
	assert.Nil(t, artist.ImageURL)
	assert.Equal(t, userID, artist.UserID)
}
