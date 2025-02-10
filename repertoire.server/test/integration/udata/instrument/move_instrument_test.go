package instrument

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

func TestMoveInstrument_WhenInstrumentIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userDataData.Users, userDataData.SeedData)

	request := requests.MoveInstrumentRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/user-data/instruments/move", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveInstrument_WhenOverInstrumentIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userDataData.Users, userDataData.SeedData)

	user := userDataData.Users[0]
	request := requests.MoveInstrumentRequest{
		ID:     user.Instruments[0].ID,
		OverID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithUser(user).
		PUT(w, "/api/user-data/instruments/move", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveInstrument_WhenSuccessful_ShouldMoveInstruments(t *testing.T) {
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

			request := requests.MoveInstrumentRequest{
				ID:     test.user.Instruments[test.index].ID,
				OverID: test.user.Instruments[test.overIndex].ID,
			}

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().
				WithUser(test.user).
				PUT(w, "/api/user-data/instruments/move", request)

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			var tunings []model.Instrument
			db := utils.GetDatabase(t)
			db.Order("\"order\"").Find(&tunings, &model.Instrument{UserID: test.user.ID})

			assertMovedInstruments(t, request, tunings, test.index, test.overIndex)
		})
	}
}

func assertMovedInstruments(
	t *testing.T,
	request requests.MoveInstrumentRequest,
	instruments []model.Instrument,
	index int,
	overIndex int,
) {
	if index < overIndex {
		assert.Equal(t, instruments[overIndex-1].ID, request.OverID)
	} else if index > overIndex {
		assert.Equal(t, instruments[overIndex+1].ID, request.OverID)
	}

	assert.Equal(t, instruments[overIndex].ID, request.ID)
	for i, tuning := range instruments {
		assert.Equal(t, uint(i), tuning.Order)
	}
}
