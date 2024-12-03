package artist

import (
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	artistData "repertoire/server/test/integration/test/data/artist"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteArtist_WhenArtistIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/artists/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteArtist_WhenSuccessful_ShouldDeleteArtist(t *testing.T) {
	tests := []struct {
		name   string
		artist model.Artist
	}{
		{
			"Without Files",
			artistData.Artists[1],
		},
		{
			"With Image",
			artistData.Artists[0],
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().DELETE(w, "/api/artists/"+test.artist.ID.String())

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			db := utils.GetDatabase()

			var deletedArtist model.Artist
			db.Find(&deletedArtist, test.artist.ID)

			assert.Empty(t, deletedArtist)
		})
	}
}
