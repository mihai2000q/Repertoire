package section

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
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

func TestUpdateSongSection_WhenRehearsalsAreDecreasing_ShouldReturnBadRequestError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	section := songData.Songs[0].Sections[0]

	request := requests.UpdateSongSectionRequest{
		ID:         section.ID,
		Name:       "New Chorus Name",
		Rehearsals: section.Rehearsals - 1,
		TypeID:     songData.Users[0].SongSectionTypes[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections", request)

	// then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSongSection_WhenSuccessful_ShouldUpdateSection(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateSongSectionRequest{
		ID:     songData.Songs[0].Sections[1].ID,
		Name:   "New Chorus Name",
		TypeID: songData.Users[0].SongSectionTypes[0].ID,
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

func TestUpdateSongSection_WhenSuccessfulWithRehearsals_ShouldUpdateSectionAddHistoryAndChangeScore(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	section := songData.Songs[0].Sections[0]
	request := requests.UpdateSongSectionRequest{
		ID:         section.ID,
		Name:       "New Chorus Name",
		Rehearsals: 15,
		Confidence: section.Confidence,
		TypeID:     songData.Users[0].SongSectionTypes[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var newSection model.SongSection
	db.Preload("History", func(*gorm.DB) *gorm.DB {
		return db.Order("created_at")
	}).Find(&newSection, &model.SongSection{ID: request.ID})

	assertUpdatedSongSection(t, newSection, request)

	assert.NotEqual(t, section.Rehearsals, newSection.Rehearsals)
	assert.NotEqual(t, section.RehearsalsScore, newSection.RehearsalsScore)
	assert.NotEqual(t, section.Progress, newSection.Progress)

	assert.NotEmpty(t, newSection.History[len(newSection.History)-1].ID)
	assert.Equal(t, section.Rehearsals, newSection.History[len(newSection.History)-1].From)
	assert.Equal(t, request.Rehearsals, newSection.History[len(newSection.History)-1].To)
	assert.Equal(t, model.RehearsalsProperty, newSection.History[len(newSection.History)-1].Property)
}

func TestUpdateSongSection_WhenSuccessfulWithConfidence_ShouldUpdateSectionAddHistoryAndChangeScore(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	section := songData.Songs[0].Sections[0]
	request := requests.UpdateSongSectionRequest{
		ID:         section.ID,
		Name:       "New Chorus Name",
		Rehearsals: section.Rehearsals,
		Confidence: 25,
		TypeID:     songData.Users[0].SongSectionTypes[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var newSection model.SongSection
	db.Preload("History", func(*gorm.DB) *gorm.DB {
		return db.Order("created_at")
	}).Find(&newSection, &model.SongSection{ID: request.ID})

	assertUpdatedSongSection(t, newSection, request)

	assert.NotEqual(t, section.Confidence, newSection.Confidence)
	assert.NotEqual(t, section.ConfidenceScore, newSection.ConfidenceScore)
	assert.NotEqual(t, section.Progress, newSection.Progress)

	assert.NotEmpty(t, newSection.History[len(newSection.History)-1].ID)
	assert.Equal(t, section.Confidence, newSection.History[len(newSection.History)-1].From)
	assert.Equal(t, request.Confidence, newSection.History[len(newSection.History)-1].To)
	assert.Equal(t, model.ConfidenceProperty, newSection.History[len(newSection.History)-1].Property)
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
