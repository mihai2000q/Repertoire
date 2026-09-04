package instrument

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	userDataData "repertoire/server/test/integration/test/data/userdata"
	"repertoire/server/test/integration/test/utils"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteInstrument_WhenInstrumentIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userDataData.Users, userDataData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/user-data/instruments/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteInstrument_WhenSuccessful_ShouldDeleteInstrument(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userDataData.Users, userDataData.SeedData)

	instrument := userDataData.Users[0].Instruments[1]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/user-data/instruments/"+instrument.ID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var instruments []model.Instrument
	db.Order("\"order\"").Find(&instruments, &model.Instrument{UserID: instrument.UserID})

	assert.True(t,
		slices.IndexFunc(instruments, func(t model.Instrument) bool {
			return t.ID == instrument.ID
		}) == -1,
		"Instrument has not been deleted",
	)

	for i := range instruments {
		assert.Equal(t, uint(i), instruments[i].Order)
	}
}
