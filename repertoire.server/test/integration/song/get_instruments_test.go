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

func TestGetInstruments_WhenSuccessful_ShouldGetInstruments(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	user := songData.Users[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithUser(user).
		GET(w, "/api/songs/instruments")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var responseInstruments []model.Instrument
	_ = json.Unmarshal(w.Body.Bytes(), &responseInstruments)

	db := utils.GetDatabase(t)

	var instruments []model.Instrument
	db.Find(&instruments, model.Instrument{UserID: user.ID})

	for i := range instruments {
		assertion.ResponseInstrument(t, instruments[i], responseInstruments[i])
	}
}
