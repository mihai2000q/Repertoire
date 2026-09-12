package section

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

func TestGetAllSongSections_WhenSuccessful_ShouldReturnSongSections(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	songID := songData.SongSections[0].SongID

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/songs/sections?songId="+songID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var responseSongSections []model.SongSection
	_ = json.Unmarshal(w.Body.Bytes(), &responseSongSections)

	db := utils.GetDatabase(t)

	var sections []model.SongSection
	db.Where(&model.SongSection{}).
		Joins("SongSectionType").
		Preload("SectionParts", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("Part").
				Joins("Part.Instrument").
				Joins("Part.BandMember").
				Preload("Part.BandMember.Roles").
				Preload("Part.Sections").
				Preload("Part.Sections.SongSectionType").
				Order("song_section_parts.order")
		}).
		Where(model.SongSection{SongID: songID}).
		Order("\"order\"").
		Find(&sections)

	for i := range responseSongSections {
		assertion.ResponseSongSection(t, sections[i], responseSongSections[i], true)
	}
}
