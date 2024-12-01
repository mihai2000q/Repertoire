package tuning

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestDeleteGuitarTuning_WhenSuccessful_ShouldDeleteTuning(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	tuning := songData.Users[0].GuitarTunings[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/guitar-tunings/"+tuning.ID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase()

	var deletedTuning model.GuitarTuning
	db.Find(&deletedTuning, tuning.ID)

	assert.Empty(t, deletedTuning)
}
