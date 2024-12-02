package types

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

func TestMoveSongSectionType_WhenTypeIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.MoveSongSectionTypeRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections/types/move", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveSongSectionType_WhenOverTypeIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	user := songData.Users[0]
	request := requests.MoveSongSectionTypeRequest{
		ID:     user.SongSectionTypes[0].ID,
		OverID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithUser(user).
		PUT(w, "/api/songs/sections/types/move", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveSongSectionType_WhenSuccessful_ShouldMoveTypes(t *testing.T) {
	tests := []struct {
		name      string
		user      model.User
		index     int
		overIndex int
	}{
		{
			"From upper position to lower",
			songData.Users[0],
			2,
			0,
		},
		{
			"From lower position to upper",
			songData.Users[0],
			0,
			2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

			request := requests.MoveSongSectionTypeRequest{
				ID:     test.user.SongSectionTypes[test.index].ID,
				OverID: test.user.SongSectionTypes[test.overIndex].ID,
			}

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().
				WithUser(test.user).
				PUT(w, "/api/songs/section/types/move", request)

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			var sectionTypes []model.SongSectionType
			db := utils.GetDatabase()
			db.Order("\"order\"").Find(&sectionTypes, &model.SongSectionType{UserID: test.user.ID})

			assertMovedTunings(t, request, sectionTypes, test.index, test.overIndex)
		})
	}
}

func assertMovedTunings(
	t *testing.T,
	request requests.MoveSongSectionTypeRequest,
	sectionTypes []model.SongSectionType,
	index int,
	overIndex int,
) {
	if index < overIndex {
		assert.Equal(t, sectionTypes[overIndex-1].ID, request.OverID)
	} else if index > overIndex {
		assert.Equal(t, sectionTypes[overIndex+1].ID, request.OverID)
	}

	assert.Equal(t, sectionTypes[overIndex].ID, request.ID)
	for i, sectionType := range sectionTypes {
		assert.Equal(t, uint(i), sectionType.Order)
	}
}
