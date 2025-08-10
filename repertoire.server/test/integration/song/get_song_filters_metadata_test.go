package song

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSongFiltersMetadata_WhenSuccessful_ShouldReturnSongFiltersMetadata(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/songs/filters-metadata")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var metadata model.SongFiltersMetadata
	_ = json.Unmarshal(w.Body.Bytes(), &metadata)

	var songs []model.Song
	db := utils.GetDatabase(t)
	db.Preload("Sections").Preload("Sections.SongSectionType").Find(&songs)

	assertion.SongFiltersMetadata(t, metadata, songs)
}
