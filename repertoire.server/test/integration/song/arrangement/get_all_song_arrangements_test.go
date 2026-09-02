package arrangement

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
	"gorm.io/gorm"
)

func TestGetAllSongArrangements_WhenSuccessful_ShouldReturnSongArrangements(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	songID := songData.SongArrangements[1].SongID

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/songs/arrangements?songId="+songID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var responseSongArrangements []model.SongArrangement
	_ = json.Unmarshal(w.Body.Bytes(), &responseSongArrangements)

	db := utils.GetDatabase(t)

	var arrangements []model.SongArrangement
	db.Where(&model.SongArrangement{}).
		Preload("PartOccurrences", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("LEFT JOIN song_parts ON song_parts.id = song_part_occurrences.part_id").
				Order("song_parts.song_order")
		}).
		Preload("PartOccurrences.Part").
		Where(model.SongArrangement{SongID: songID}).
		Order("\"order\"").
		Find(&arrangements)

	for i := range responseSongArrangements {
		assertion.ResponseSongArrangement(t, arrangements[i], responseSongArrangements[i], true)
	}
}
