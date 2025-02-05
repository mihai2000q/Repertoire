package tuning

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	userDataData "repertoire/server/test/integration/test/data/udata"
	"repertoire/server/test/integration/test/utils"
	"slices"
	"testing"
)

func TestDeleteGuitarTuning_WhenTuningIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userDataData.Users, userDataData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/guitar-tunings/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteGuitarTuning_WhenSuccessful_ShouldDeleteTuning(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userDataData.Users, userDataData.SeedData)

	tuning := userDataData.Users[0].GuitarTunings[1]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/guitar-tunings/"+tuning.ID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var tunings []model.GuitarTuning
	db.Order("\"order\"").Find(&tunings, &model.GuitarTuning{UserID: tuning.UserID})

	assert.True(t,
		slices.IndexFunc(tunings, func(t model.GuitarTuning) bool {
			return t.ID == tuning.ID
		}) == -1,
		"Guitar Tuning has not been deleted",
	)

	for i := range tunings {
		assert.Equal(t, uint(i), tunings[i].Order)
	}
}
