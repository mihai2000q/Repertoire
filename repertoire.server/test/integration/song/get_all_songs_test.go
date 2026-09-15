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

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetAllSongs_WhenSuccessful_ShouldReturnSongs(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/songs")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var responseSongs []model.EnhancedSong
	_ = json.Unmarshal(w.Body.Bytes(), &responseSongs)

	db := utils.GetDatabase(t)

	var songs []model.Song
	db.Joins("Album").
		Joins("Artist").
		Joins("GuitarTuning").
		Preload("Parts").
		Preload("Sections").
		Preload("Sections.SongSectionType").
		Find(&songs)

	for i := range responseSongs {
		assertion.ResponseEnhancedSong(t, songs[i], responseSongs[i], nil)
	}
}

func TestGetAllSongs_WhenSuccessfulWithArrangements_ShouldReturnSongsWithArrangements(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/songs?with=Arrangements")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var responseSongs []model.EnhancedSong
	_ = json.Unmarshal(w.Body.Bytes(), &responseSongs)

	db := utils.GetDatabase(t)

	var songs []model.Song
	db.Joins("Album").
		Joins("Artist").
		Joins("GuitarTuning").
		Preload("Parts").
		Preload("Sections").
		Preload("Sections.SongSectionType").
		Preload("Arrangements", func(db *gorm.DB) *gorm.DB {
			return db.Order("song_arrangements.order")
		}).
		Find(&songs)

	for i := range responseSongs {
		assertion.ResponseEnhancedSong(t, songs[i], responseSongs[i], []string{"Arrangements"})
	}
}
