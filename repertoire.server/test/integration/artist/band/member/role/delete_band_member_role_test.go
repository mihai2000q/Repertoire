package role

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	artistData "repertoire/server/test/integration/test/data/artist"
	"repertoire/server/test/integration/test/utils"
	"slices"
	"testing"
)

func TestDeleteBandMemberRole_WhenRoleIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/artists/band-members/roles/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteBandMemberRole_WhenSuccessful_ShouldDeleteRole(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	bandMemberRole := artistData.Users[0].BandMemberRoles[1]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/artists/band-members/roles/"+bandMemberRole.ID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var bandMemberRoles []model.BandMemberRole
	db.Order("\"order\"").Find(&bandMemberRoles, &model.BandMemberRole{UserID: bandMemberRole.UserID})

	assert.True(t,
		slices.IndexFunc(bandMemberRoles, func(t model.BandMemberRole) bool {
			return t.ID == bandMemberRole.ID
		}) == -1,
		"Band Member Role has not been deleted",
	)

	for i := range bandMemberRoles {
		assert.Equal(t, uint(i), bandMemberRoles[i].Order)
	}
}
