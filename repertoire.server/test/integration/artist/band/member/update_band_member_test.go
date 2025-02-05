package member

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	artistData "repertoire/server/test/integration/test/data/artist"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateBandMember_WhenArtistIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	request := requests.UpdateBandMemberRequest{
		ID:      uuid.New(),
		Name:    "New Name",
		RoleIDs: []uuid.UUID{uuid.New()},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/artists/band-members", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateBandMember_WhenSuccessful_ShouldUpdateBandMember(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	bandMember := artistData.Artists[0].BandMembers[0]

	request := requests.UpdateBandMemberRequest{
		ID:      bandMember.ID,
		Name:    "New Name",
		RoleIDs: []uuid.UUID{artistData.Users[0].BandMemberRoles[2].ID},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/artists/band-members", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)
	db.Preload("Roles").Find(&bandMember, bandMember.ID)

	assertUpdatedBandMember(t, request, bandMember)
}

func assertUpdatedBandMember(t *testing.T, request requests.UpdateBandMemberRequest, bandMember model.BandMember) {
	assert.Equal(t, request.ID, bandMember.ID)
	assert.Equal(t, request.Name, bandMember.Name)
	assert.Equal(t, request.Color, bandMember.Color)
	assert.Len(t, bandMember.Roles, len(request.RoleIDs))
	for i, role := range bandMember.Roles {
		assert.Equal(t, request.RoleIDs[i], role.ID)
	}
}
