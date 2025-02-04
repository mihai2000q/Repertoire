package member

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

func TestCreateBandMember_WhenSuccessful_ShouldCreateMember(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artist := artistData.Artists[0]
	request := requests.CreateBandMemberRequest{
		ArtistID: artist.ID,
		Name:     "Guitarist 1-New",
		RoleIDs:  []uuid.UUID{artistData.Users[0].BandMemberRoles[2].ID},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/artists/band-members", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var section model.BandMember
	db.Preload("BandMembers").
		Preload("BandMembers.Roles").
		Find(&section, &model.BandMember{Name: request.Name})

	assertCreatedBandMember(t, section, request, len(artist.BandMembers))
}

func assertCreatedBandMember(
	t *testing.T,
	bandMember model.BandMember,
	request requests.CreateBandMemberRequest,
	order int,
) {
	assert.NotEmpty(t, bandMember.ID)
	assert.Equal(t, request.ArtistID, bandMember.ArtistID)
	assert.Equal(t, request.Name, bandMember.Name)
	assert.Nil(t, bandMember.ImageURL)
	assert.Equal(t, uint(order), bandMember.Order)
	assert.Len(t, bandMember.Roles, len(request.RoleIDs))
	for i, role := range bandMember.Roles {
		assert.Equal(t, request.RoleIDs[i], role.ID)
	}
}
