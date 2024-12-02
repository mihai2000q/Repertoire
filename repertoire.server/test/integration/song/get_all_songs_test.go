package song

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

func TestGetAllSongs_WhenSuccessful_ShouldReturnSongs(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/songs")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var responseSongs []model.Song
	_ = json.Unmarshal(w.Body.Bytes(), &responseSongs)

	db := utils.GetDatabase()

	var songs []model.Song
	db.Preload("Album").
		Preload("Artist").
		Preload("GuitarTuning").
		Preload("Sections").
		Preload("Sections.SongSectionType").
		Find(&songs)

	for i := range responseSongs {
		assertion.ResponseSong(t, songs[i], responseSongs[i], true, true, true)
	}
}
