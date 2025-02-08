package member

import (
	"encoding/json"
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

func TestCreateBandMember_WhenArtistIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	request := requests.CreateBandMemberRequest{
		ArtistID: uuid.New(),
		Name:     "Guitarist 1-New",
		RoleIDs:  []uuid.UUID{artistData.Users[0].BandMemberRoles[2].ID},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/artists/band-members", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateBandMember_WhenArtistIsNotBand_ShouldReturnBadRequestError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artist := artistData.Artists[1]
	request := requests.CreateBandMemberRequest{
		ArtistID: artist.ID,
		Name:     "Guitarist 1-New",
		RoleIDs:  []uuid.UUID{artistData.Users[0].BandMemberRoles[2].ID},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/artists/band-members", request)

	// then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

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
	var response struct{ ID uuid.UUID }
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, response)

	db := utils.GetDatabase(t)

	var member model.BandMember
	db.Preload("Roles").Find(&member, &model.BandMember{Name: request.Name})

	assertCreatedBandMember(t, member, request, len(artist.BandMembers))
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
	assert.Equal(t, request.Color, bandMember.Color)
	assert.Nil(t, bandMember.ImageURL)
	assert.Equal(t, uint(order), bandMember.Order)
	assert.Len(t, bandMember.Roles, len(request.RoleIDs))
	for i, role := range bandMember.Roles {
		assert.Equal(t, request.RoleIDs[i], role.ID)
	}
}
