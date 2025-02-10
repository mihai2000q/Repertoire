package song

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"
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
	db.Joins("Album").
		Joins("Artist").
		Joins("GuitarTuning").
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("song_sections.order")
		}).
		Preload("Sections.SongSectionType").
		Preload("Sections.Instrument").
		Preload("Sections.BandMember").
		Preload("Sections.BandMember.Roles").
		Preload("Artist.BandMembers").
		Preload("Artist.BandMembers.Roles").
		Find(&song, song.ID)

	assertion.ResponseSong(
		t,
		song,
		responseSong,
		true,
		true,
		true,
		true,
	)
}
