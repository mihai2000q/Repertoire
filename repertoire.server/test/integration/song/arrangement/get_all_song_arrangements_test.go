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

	var songs []model.SongArrangement
	db.Where(&model.SongArrangement{}).
		Preload("SectionOccurrences", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("LEFT JOIN song_sections ON song_sections.id = song_section_occurrences.section_id").
				Order("song_sections.order")
		}).
		Preload("SectionOccurrences.Section").
		Preload("SectionOccurrences.Section.SongSectionType").
		Where(model.SongArrangement{SongID: songID}).
		Order("\"order\"").
		Find(&songs)

	for i := range responseSongArrangements {
		assertion.ResponseSongArrangement(t, songs[i], responseSongArrangements[i], true)
	}
}
