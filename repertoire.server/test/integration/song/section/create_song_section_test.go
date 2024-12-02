package section

import (
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

func TestCreateSongSection_WhenSuccessful_ShouldCreateSection(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.CreateSongSectionRequest{
		SongID: songData.Songs[0].ID,
		Name:   "Chorus 1-New",
		TypeID: songData.Users[0].SongSectionTypes[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/sections", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase()

	var section model.SongSection
	db.Find(&section, &model.SongSection{Name: request.Name})

	assertCreatedSongSection(t, section, request, len(songData.Songs[0].Sections))
}

func assertCreatedSongSection(
	t *testing.T,
	songSection model.SongSection,
	request requests.CreateSongSectionRequest,
	order int,
) {
	assert.NotEmpty(t, songSection.ID)
	assert.Equal(t, request.SongID, songSection.SongID)
	assert.Equal(t, request.Name, songSection.Name)
	assert.Equal(t, request.TypeID, songSection.SongSectionTypeID)
	assert.Zero(t, songSection.Rehearsals)
	assert.Equal(t, uint(order), songSection.Order)
}
