package role

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	artistData "repertoire/server/test/integration/test/data/artist"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestCreateBandMemberRole_WhenSuccessful_ShouldCreateRole(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	user := artistData.Users[0]

	request := requests.CreateBandMemberRoleRequest{
		Name: "Guitarist-New",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithUser(user).
		POST(w, "/api/artists/band-members/roles", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var bandMemberRole model.BandMemberRole
	db.Find(&bandMemberRole, &model.BandMemberRole{Name: request.Name})

	assertCreatedBandMemberRole(t, bandMemberRole, request, len(artistData.Users[0].BandMemberRoles), user.ID)
}

func assertCreatedBandMemberRole(
	t *testing.T,
	bandMemberRole model.BandMemberRole,
	request requests.CreateBandMemberRoleRequest,
	order int,
	userID uuid.UUID,
) {
	assert.NotEmpty(t, bandMemberRole.ID)
	assert.Equal(t, request.Name, bandMemberRole.Name)
	assert.Equal(t, userID, bandMemberRole.UserID)
	assert.Equal(t, uint(order), bandMemberRole.Order)
}
