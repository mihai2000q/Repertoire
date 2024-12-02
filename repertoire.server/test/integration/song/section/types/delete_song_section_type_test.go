package types

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestDeleteSongSectionType_WhenSuccessful_ShouldDeleteType(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	sectionType := songData.Users[0].SongSectionTypes[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/sections/types/"+sectionType.ID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase()

	var deletedType model.SongSectionType
	db.Find(&deletedType, sectionType.ID)

	assert.Empty(t, deletedType)
}
