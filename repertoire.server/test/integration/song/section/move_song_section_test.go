package section

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMoveSongSection_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.MoveSongSectionRequest{
		SongID: uuid.New(),
		ID:     uuid.New(),
		OverID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections/move", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveSongSection_WhenSectionIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]
	request := requests.MoveSongSectionRequest{
		SongID: song.ID,
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
	request := requests.MoveSongSectionRequest{
		SongID: song.ID,
		ID:     songData.SongSections[0].ID,
		OverID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections/move", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveSongSection_WhenSuccessful_ShouldMoveTypes(t *testing.T) {
	tests := []struct {
		name      string
		request   requests.MoveSongSectionRequest
		index     int
		overIndex int
	}{
		{
			"From upper position to lower",
			requests.MoveSongSectionRequest{
				SongID: songData.SongSections[0].SongID,
				ID:     songData.SongSections[2].ID,
				OverID: songData.SongSections[0].ID,
			},
			2,
			0,
		},
		{
			"From lower position to upper",
			requests.MoveSongSectionRequest{
				SongID: songData.SongSections[0].SongID,
				ID:     songData.SongSections[0].ID,
				OverID: songData.SongSections[2].ID,
			},
			0,
			2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().PUT(w, "/api/songs/sections/move", test.request)

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			var sections []model.SongSection
			db := utils.GetDatabase(t)
			db.Where(&model.SongSection{SongID: test.request.SongID}).Order("\"order\"").Find(&sections)

			assertMovedSections(t, test.request, sections, test.index, test.overIndex)
		})
	}
}

func assertMovedSections(
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
	for i, section := range sections {
		assert.Equal(t, uint(i), section.Order)
	}
}
