package artist

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	artistData "repertoire/server/test/integration/test/data/artist"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteArtist_WhenSuccessful_ShouldDeleteArtist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artist := artistData.Artists[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/artists/"+artist.ID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase()

	var deletedArtist model.Artist
	db.Find(&deletedArtist, artist.ID)

	assert.Empty(t, deletedArtist)
}
