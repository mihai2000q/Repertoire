package section

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestUpdateSongSection_WhenSectionIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateSongSectionRequest{
		ID:     uuid.New(),
		Name:   "New Chorus Name",
		TypeID: songData.Users[0].SongSectionTypes[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateSongSection_WhenSuccessful_ShouldUpdateSection(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateSongSectionRequest{
		ID:         songData.Songs[0].Sections[0].ID,
		Name:       "New Chorus Name",
		Confidence: 99,
		Rehearsals: 23,
		TypeID:     songData.Users[0].SongSectionTypes[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var section model.SongSection
	db.Find(&section, &model.SongSection{ID: request.ID})

	assertUpdatedSongSection(t, section, request)
}

func assertUpdatedSongSection(
	t *testing.T,
	songSection model.SongSection,
	request requests.UpdateSongSectionRequest,
) {
	assert.Equal(t, request.Name, songSection.Name)
	assert.Equal(t, request.Confidence, songSection.Confidence)
	assert.Equal(t, request.Rehearsals, songSection.Rehearsals)
	assert.Equal(t, request.TypeID, songSection.SongSectionTypeID)
}
