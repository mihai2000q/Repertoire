package section

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestGetSongSectionTypes_WhenSuccessful_ShouldGetTypes(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	user := songData.Users[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithUser(user).
		GET(w, "/api/songs/sections/types")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var responseTypes []model.SongSectionType
	_ = json.Unmarshal(w.Body.Bytes(), &responseTypes)

	db := utils.GetDatabase(t)

	var sectionTypes []model.SongSectionType
	db.Find(&sectionTypes, model.SongSectionType{UserID: user.ID})

	for i := range sectionTypes {
		assertion.ResponseSongSectionType(t, sectionTypes[i], responseTypes[i])
	}
}
