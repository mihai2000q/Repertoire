package types

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	userDataData "repertoire/server/test/integration/test/data/udata"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestCreateSongSectionType_WhenSuccessful_ShouldCreateType(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userDataData.Users, userDataData.SeedData)

	user := userDataData.Users[0]

	request := requests.CreateSongSectionTypeRequest{
		Name: "Chorus-New",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithUser(user).
		POST(w, "/api/user-data/song-section-types", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var sectionType model.SongSectionType
	db.Find(&sectionType, &model.SongSectionType{Name: request.Name})

	assertCreatedSongSectionType(t, sectionType, request, len(userDataData.Users[0].SongSectionTypes), user.ID)
}

func assertCreatedSongSectionType(
	t *testing.T,
	songSectionType model.SongSectionType,
	request requests.CreateSongSectionTypeRequest,
	order int,
	userID uuid.UUID,
) {
	assert.NotEmpty(t, songSectionType.ID)
	assert.Equal(t, request.Name, songSectionType.Name)
	assert.Equal(t, userID, songSectionType.UserID)
	assert.Equal(t, uint(order), songSectionType.Order)
}
