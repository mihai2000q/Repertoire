package types

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	userDataData "repertoire/server/test/integration/test/data/udata"
	"repertoire/server/test/integration/test/utils"
	"slices"
	"testing"
)

func TestDeleteSongSectionType_WhenTypeIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userDataData.Users, userDataData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/user-data/song-section-types/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteSongSectionType_WhenSuccessful_ShouldDeleteType(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userDataData.Users, userDataData.SeedData)

	sectionType := userDataData.Users[0].SongSectionTypes[1]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/user-data/song-section-types/"+sectionType.ID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var sectionTypes []model.SongSectionType
	db.Order("\"order\"").Find(&sectionTypes, &model.SongSectionType{UserID: sectionType.UserID})

	assert.True(t,
		slices.IndexFunc(sectionTypes, func(t model.SongSectionType) bool {
			return t.ID == sectionType.ID
		}) == -1,
		"Song Section Type has not been deleted",
	)

	for i := range sectionTypes {
		assert.Equal(t, uint(i), sectionTypes[i].Order)
	}
}
