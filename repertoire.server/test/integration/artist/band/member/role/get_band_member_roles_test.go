package role

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	artistData "repertoire/server/test/integration/test/data/artist"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestGetBandMemberRoles_WhenSuccessful_ShouldGetRoles(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	user := artistData.Users[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithUser(user).
		GET(w, "/api/artists/band-members/roles")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var responseRoles []model.BandMemberRole
	_ = json.Unmarshal(w.Body.Bytes(), &responseRoles)

	db := utils.GetDatabase(t)

	var bandMemberRoles []model.BandMemberRole
	db.Find(&bandMemberRoles, model.BandMemberRole{UserID: user.ID})

	for i := range bandMemberRoles {
		assertion.ResponseBandMemberRole(t, bandMemberRoles[i], responseRoles[i])
	}
}
