package song

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestGetGuitarTunings_WhenSuccessful_ShouldGetTunings(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	user := songData.Users[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithUser(user).
		GET(w, "/api/songs/guitar-tunings")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var responseTunings []model.GuitarTuning
	_ = json.Unmarshal(w.Body.Bytes(), &responseTunings)

	db := utils.GetDatabase(t)

	var tunings []model.GuitarTuning
	db.Find(&tunings, model.GuitarTuning{UserID: user.ID})

	for i := range tunings {
		assertion.ResponseGuitarTuning(t, tunings[i], responseTunings[i])
	}
}
