package section

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDeleteSongSection_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/sections/"+uuid.New().String()+"/from/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteSongSection_WhenSectionIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/sections/"+uuid.New().String()+"/from/"+song.ID.String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteSongSection_WhenSuccessful_ShouldDeleteSection(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// song with sections and previous stats
	song := songData.Songs[0]
	section := song.Sections[1]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/sections/"+section.ID.String()+"/from/"+section.SongID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var newSong model.Song
	db.Preload("Sections", func(db gorm.DB) *gorm.DB {
		return db.Order("\"order\"")
	}).Find(&song, song.ID)

	assert.True(t,
		slices.IndexFunc(newSong.Sections, func(t model.SongSection) bool {
			return t.ID == section.ID
		}) == -1,
		"Song Section has not been deleted",
	)

	for i, s := range newSong.Sections {
		assert.Equal(t, uint(i), s.Order)
	}

	assert.LessOrEqual(t, newSong.Confidence, song.Confidence)
	assert.LessOrEqual(t, newSong.Rehearsals, song.Rehearsals)
	assert.LessOrEqual(t, newSong.Progress, song.Progress)
}
