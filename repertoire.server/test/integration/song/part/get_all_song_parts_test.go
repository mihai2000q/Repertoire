package part

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllSongParts_WhenSuccessful_ShouldReturnSongParts(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	songID := songData.SongParts[0].SongID

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/songs/parts?songId="+songID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var responseSongParts []model.SongPart
	_ = json.Unmarshal(w.Body.Bytes(), &responseSongParts)

	db := utils.GetDatabase(t)

	var parts []model.SongPart
	db.Where(&model.SongPart{}).
		Joins("Instrument").
		Joins("BandMember").
		Preload("BandMember.Roles").
		Preload("Sections").
		Preload("Sections.SongSectionType").
		Where(model.SongPart{SongID: songID}).
		Order("song_order").
		Find(&parts)

	for i := range responseSongParts {
		assertion.ResponseSongPart(t, parts[i], responseSongParts[i], true, true)
	}
}
