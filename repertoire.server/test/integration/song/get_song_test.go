package song

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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetSong_WhenUserIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/songs/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetSong_WhenSuccessful_ShouldReturnSong(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/songs/"+song.ID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var responseSong model.Song
	_ = json.Unmarshal(w.Body.Bytes(), &responseSong)

	db := utils.GetDatabase(t)
	db.
		Joins("Settings").
		Joins("Settings.DefaultBandMember").
		Joins("Settings.DefaultInstrument").
		Joins("GuitarTuning").
		Joins("Artist").
		Joins("Album").
		Preload("Artist.BandMembers").
		Preload("Artist.BandMembers.Roles").
		Preload("Parts", func(db *gorm.DB) *gorm.DB {
			return db.Order("song_parts.song_order")
		}).
		Preload("Parts.Instrument").
		Preload("Parts.BandMember").
		Preload("Parts.BandMember.Roles").
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("song_sections.order")
		}).
		Preload("Sections.SongSectionType").
		Preload("Sections.SectionParts", func(db *gorm.DB) *gorm.DB {
			return db.Order("song_section_parts.order")
		}).
		Preload("Sections.SectionParts.Part").
		Preload("Sections.SectionParts.Part.Instrument").
		Preload("Sections.SectionParts.Part.BandMember.Roles").
		Find(&song, song.ID)

	assertion.ResponseSong(
		t,
		song,
		responseSong,
		true,
		true,
		true,
	)
}
