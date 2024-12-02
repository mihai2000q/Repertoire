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

func TestMoveSongSection_WhenSectionIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.MoveSongSectionRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections/move", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveSongSection_WhenOverSectionIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]
	request := requests.MoveSongSectionTypeRequest{
		ID:     song.Sections[0].ID,
		OverID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections/move", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveSongSectionType_WhenSuccessful_ShouldMoveTypes(t *testing.T) {
	tests := []struct {
		name      string
		song      model.Song
		index     int
		overIndex int
	}{
		{
			"From upper position to lower",
			songData.Songs[0],
			2,
			0,
		},
		{
			"From lower position to upper",
			songData.Songs[0],
			0,
			2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

			request := requests.MoveSongSectionRequest{
				ID:     test.song.Sections[test.index].ID,
				OverID: test.song.Sections[test.overIndex].ID,
			}

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().PUT(w, "/api/songs/section/move", request)

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			var sections []model.SongSection
			db := utils.GetDatabase()
			db.Order("\"order\"").Find(&sections, &model.SongSection{SongID: test.song.ID})

			assertMovedTunings(t, request, sections, test.index, test.overIndex)
		})
	}
}

func assertMovedTunings(
	t *testing.T,
	request requests.MoveSongSectionRequest,
	sections []model.SongSection,
	index int,
	overIndex int,
) {
	if index < overIndex {
		assert.Equal(t, sections[overIndex-1].ID, request.OverID)
	} else if index > overIndex {
		assert.Equal(t, sections[overIndex+1].ID, request.OverID)
	}

	assert.Equal(t, sections[overIndex].ID, request.ID)
	for i, sectionType := range sections {
		assert.Equal(t, uint(i), sectionType.Order)
	}
}
