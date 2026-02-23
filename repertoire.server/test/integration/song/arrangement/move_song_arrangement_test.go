package arrangement

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

func TestMoveSongArrangement_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.MoveSongArrangementRequest{
		SongID: uuid.New(),
		ID:     uuid.New(),
		OverID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/arrangements/move", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveSongArrangement_WhenArrangementIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]
	request := requests.MoveSongArrangementRequest{
		SongID: song.ID,
		ID:     uuid.New(),
		OverID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/arrangements/move", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveSongArrangement_WhenOverArrangementIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	songArrangement := songData.SongArrangements[0]
	request := requests.MoveSongArrangementRequest{
		SongID: songArrangement.SongID,
		ID:     songArrangement.ID,
		OverID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/arrangements/move", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveSongArrangementType_WhenSuccessful_ShouldMoveTypes(t *testing.T) {
	tests := []struct {
		name      string
		request   requests.MoveSongArrangementRequest
		index     int
		overIndex int
	}{
		{
			"From upper position to lower",
			requests.MoveSongArrangementRequest{
				SongID: songData.SongArrangements[0].SongID,
				ID:     songData.SongArrangements[2].ID,
				OverID: songData.SongArrangements[0].ID,
			},
			2,
			0,
		},
		{
			"From lower position to upper",
			requests.MoveSongArrangementRequest{
				SongID: songData.SongArrangements[0].SongID,
				ID:     songData.SongArrangements[0].ID,
				OverID: songData.SongArrangements[2].ID,
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
			core.NewTestHandler().PUT(w, "/api/songs/arrangements/move", test.request)

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			var arrangements []model.SongArrangement
			db := utils.GetDatabase(t)
			db.Order("\"order\"").Find(&arrangements, &model.SongArrangement{SongID: test.request.SongID})

			assertMovedArrangements(t, test.request, arrangements, test.index, test.overIndex)
		})
	}
}

func assertMovedArrangements(
	t *testing.T,
	request requests.MoveSongArrangementRequest,
	arrangements []model.SongArrangement,
	index int,
	overIndex int,
) {
	if index < overIndex {
		assert.Equal(t, arrangements[overIndex-1].ID, request.OverID)
	} else if index > overIndex {
		assert.Equal(t, arrangements[overIndex+1].ID, request.OverID)
	}

	assert.Equal(t, arrangements[overIndex].ID, request.ID)
	for i, arrangement := range arrangements {
		assert.Equal(t, uint(i), arrangement.Order)
	}
}
