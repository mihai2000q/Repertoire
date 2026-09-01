package part

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

func TestDeleteSongPart_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/parts/"+uuid.New().String()+"/from/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteSongPart_WhenPartIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/parts/"+uuid.New().String()+"/from/"+song.ID.String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteSongPart_WhenSuccessful_ShouldDeletePart(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// song with parts and previous stats
	song := songData.Songs[0]
	part := songData.SongParts[1]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/parts/"+part.ID.String()+"/from/"+part.SongID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var newSong model.Song
	db.
		Preload("Parts", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"song_order\"")
		}).
		Preload("Sections").
		Preload("Sections.SectionParts", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\"")
		}).
		Find(&newSong, song.ID)

	assert.True(t,
		slices.IndexFunc(newSong.Parts, func(t model.SongPart) bool {
			return t.ID == part.ID
		}) == -1,
		"Song Part has not been deleted",
	)

	for i, s := range newSong.Parts {
		assert.Equal(t, uint(i), s.SongOrder)
	}

	assert.LessOrEqual(t, newSong.Confidence, song.Confidence)
	assert.LessOrEqual(t, newSong.Rehearsals, song.Rehearsals)
	assert.LessOrEqual(t, newSong.Progress, song.Progress)

	for _, sec := range newSong.Sections {
		for i, sp := range sec.SectionParts {
			assert.Equal(t, uint(i), sp.Order)
		}
	}
}
