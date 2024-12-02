package section

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"slices"
	"testing"
)

func TestDeleteSongSection_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/sections/"+uuid.New().String()+"/from/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteSongSection_WhenSectionIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/sections/"+uuid.New().String()+"/from/"+song.ID.String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteSongSection_WhenSuccessful_ShouldDeleteSection(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	section := songData.Songs[0].Sections[1]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/sections/"+section.ID.String()+"/from/"+section.SongID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase()

	var sections []model.SongSection
	db.Order("\"order\"").Find(&sections, &model.SongSection{SongID: section.SongID})

	assert.True(t,
		slices.IndexFunc(sections, func(t model.SongSection) bool {
			return t.ID == section.ID
		}) == -1,
		"Song Section has not been deleted",
	)

	for i := range sections {
		assert.Equal(t, uint(i), sections[i].Order)
	}
}
