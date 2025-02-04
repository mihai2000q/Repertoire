package member

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

func TestDeleteBandMember_WhenArtistIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/artists/band-members/"+uuid.New().String()+"/from/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteBandMember_WhenMemberIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artist := artistData.Artists[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/artists/band-members/"+uuid.New().String()+"/from/"+artist.ID.String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteBandMember_WhenSuccessful_ShouldDeleteMember(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	bandMember := artistData.Artists[0].BandMembers[1]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/artists/band-members/"+bandMember.ID.String()+"/from/"+bandMember.ArtistID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var sections []model.BandMember
	db.Order("\"order\"").Find(&sections, &model.BandMember{ArtistID: bandMember.ArtistID})

	assert.True(t,
		slices.IndexFunc(sections, func(t model.BandMember) bool {
			return t.ID == bandMember.ID
		}) == -1,
		"Artist Section has not been deleted",
	)

	for i := range sections {
		assert.Equal(t, uint(i), sections[i].Order)
	}
}
