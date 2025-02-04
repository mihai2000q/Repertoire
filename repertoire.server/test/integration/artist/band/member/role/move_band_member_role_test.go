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

func TestMoveBandMemberRole_WhenRoleIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	request := requests.MoveBandMemberRoleRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/artists/band-members/roles/move", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveBandMemberRole_WhenOverRoleIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	user := artistData.Users[0]
	request := requests.MoveBandMemberRoleRequest{
		ID:     user.BandMemberRoles[0].ID,
		OverID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithUser(user).
		PUT(w, "/api/artists/band-members/roles/move", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveBandMemberRole_WhenSuccessful_ShouldMoveRoles(t *testing.T) {
	tests := []struct {
		name      string
		user      model.User
		index     int
		overIndex int
	}{
		{
			"From upper position to lower",
			artistData.Users[0],
			2,
			0,
		},
		{
			"From lower position to upper",
			artistData.Users[0],
			0,
			2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

			request := requests.MoveBandMemberRoleRequest{
				ID:     test.user.BandMemberRoles[test.index].ID,
				OverID: test.user.BandMemberRoles[test.overIndex].ID,
			}

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().
				WithUser(test.user).
				PUT(w, "/api/artists/band-members/roles/move", request)

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			var bandMemberRoles []model.BandMemberRole
			db := utils.GetDatabase(t)
			db.Order("\"order\"").Find(&bandMemberRoles, &model.BandMemberRole{UserID: test.user.ID})

			assertMovedRoles(t, request, bandMemberRoles, test.index, test.overIndex)
		})
	}
}

func assertMovedRoles(
	t *testing.T,
	request requests.MoveBandMemberRoleRequest,
	bandMemberRoles []model.BandMemberRole,
	index int,
	overIndex int,
) {
	if index < overIndex {
		assert.Equal(t, bandMemberRoles[overIndex-1].ID, request.OverID)
	} else if index > overIndex {
		assert.Equal(t, bandMemberRoles[overIndex+1].ID, request.OverID)
	}

	assert.Equal(t, bandMemberRoles[overIndex].ID, request.ID)
	for i, bandMemberRole := range bandMemberRoles {
		assert.Equal(t, uint(i), bandMemberRole.Order)
	}
}
