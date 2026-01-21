package instrument

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	userDataData "repertoire/server/test/integration/test/data/udata"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateInstrument_WhenSuccessful_ShouldCreateInstrument(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userDataData.Users, userDataData.SeedData)

	user := userDataData.Users[0]

	request := requests.CreateInstrumentRequest{
		Name: "New Voice New Me",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithUser(user).
		POST(w, "/api/user-data/instruments", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var Instrument model.Instrument
	db.Find(&Instrument, &model.Instrument{Name: request.Name})

	assertCreatedInstrument(t, Instrument, request, len(userDataData.Users[0].Instruments), user.ID)
}

func assertCreatedInstrument(
	t *testing.T,
	instrument model.Instrument,
	request requests.CreateInstrumentRequest,
	order int,
	userID uuid.UUID,
) {
	assert.NotEmpty(t, instrument.ID)
	assert.Equal(t, request.Name, instrument.Name)
	assert.Equal(t, userID, instrument.UserID)
	assert.Equal(t, uint(order), instrument.Order)
}
