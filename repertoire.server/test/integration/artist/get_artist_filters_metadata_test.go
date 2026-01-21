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

func TestGetArtistFiltersMetadata_WhenSuccessful_ShouldReturnArtistFiltersMetadata(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/artists/filters-metadata")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var metadata model.ArtistFiltersMetadata
	_ = json.Unmarshal(w.Body.Bytes(), &metadata)

	var artists []model.Artist
	db := utils.GetDatabase(t)
	db.Preload("BandMembers").
		Preload("Albums").
		Preload("Songs").
		Find(&artists)

	assertion.ArtistFiltersMetadata(t, metadata, artists)
}
