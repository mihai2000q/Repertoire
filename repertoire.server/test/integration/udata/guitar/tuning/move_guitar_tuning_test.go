package tuning

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	userDataData "repertoire/server/test/integration/test/data/udata"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestMoveGuitarTuning_WhenTuningIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userDataData.Users, userDataData.SeedData)

	request := requests.MoveGuitarTuningRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/guitar-tunings/move", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveGuitarTuning_WhenOverTuningIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userDataData.Users, userDataData.SeedData)

	user := userDataData.Users[0]
	request := requests.MoveGuitarTuningRequest{
		ID:     user.GuitarTunings[0].ID,
		OverID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithUser(user).
		PUT(w, "/api/songs/guitar-tunings/move", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveGuitarTuning_WhenSuccessful_ShouldMoveTunings(t *testing.T) {
	tests := []struct {
		name      string
		user      model.User
		index     int
		overIndex int
	}{
		{
			"From upper position to lower",
			userDataData.Users[0],
			2,
			0,
		},
		{
			"From lower position to upper",
			userDataData.Users[0],
			0,
			2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, userDataData.Users, userDataData.SeedData)

			request := requests.MoveGuitarTuningRequest{
				ID:     test.user.GuitarTunings[test.index].ID,
				OverID: test.user.GuitarTunings[test.overIndex].ID,
			}

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().
				WithUser(test.user).
				PUT(w, "/api/songs/guitar-tunings/move", request)

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			var tunings []model.GuitarTuning
			db := utils.GetDatabase(t)
			db.Order("\"order\"").Find(&tunings, &model.GuitarTuning{UserID: test.user.ID})

			assertMovedTunings(t, request, tunings, test.index, test.overIndex)
		})
	}
}

func assertMovedTunings(
	t *testing.T,
	request requests.MoveGuitarTuningRequest,
	tunings []model.GuitarTuning,
	index int,
	overIndex int,
) {
	if index < overIndex {
		assert.Equal(t, tunings[overIndex-1].ID, request.OverID)
	} else if index > overIndex {
		assert.Equal(t, tunings[overIndex+1].ID, request.OverID)
	}

	assert.Equal(t, tunings[overIndex].ID, request.ID)
	for i, tuning := range tunings {
		assert.Equal(t, uint(i), tuning.Order)
	}
}
