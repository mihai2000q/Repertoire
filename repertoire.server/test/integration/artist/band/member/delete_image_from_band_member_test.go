package member

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/test/integration/test/core"
	artistData "repertoire/server/test/integration/test/data/artist"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteImageFromBandMember_WhenMemberIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/artists/band-members/images/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteImageFromBandMember_WhenArtistHasNoImage_ShouldReturnConflictError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	bandMember := artistData.Artists[0].BandMembers[1]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/artists/band-members/images/"+bandMember.ID.String())

	// then
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestDeleteImageFromBandMember_WhenSuccessful_ShouldUpdateArtistAndDeleteImage(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	bandMember := artistData.Artists[0].BandMembers[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/artists/band-members/images/"+bandMember.ID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)
	db.Find(&bandMember, bandMember.ID)

	assert.Nil(t, bandMember.ImageURL)
}
