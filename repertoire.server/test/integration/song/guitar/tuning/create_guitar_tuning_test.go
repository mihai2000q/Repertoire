package tuning

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

func TestCreateGuitarTuning_WhenSuccessful_ShouldCreateTuning(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	user := songData.Users[0]

	request := requests.CreateGuitarTuningRequest{
		Name: "Eb Standard-New",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithUser(user).
		POST(w, "/api/songs/guitar-tunings", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase()

	var guitarTuning model.GuitarTuning
	db.Find(&guitarTuning, &model.GuitarTuning{Name: request.Name})

	assertCreatedGuitarTuning(t, guitarTuning, request, len(songData.Users[0].GuitarTunings), user.ID)
}

func assertCreatedGuitarTuning(
	t *testing.T,
	guitarTuning model.GuitarTuning,
	request requests.CreateGuitarTuningRequest,
	order int,
	userID uuid.UUID,
) {
	assert.NotEmpty(t, guitarTuning.ID)
	assert.Equal(t, request.Name, guitarTuning.Name)
	assert.Equal(t, userID, guitarTuning.UserID)
	assert.Equal(t, uint(order), guitarTuning.Order)
}
