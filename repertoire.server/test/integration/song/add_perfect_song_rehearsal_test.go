package song

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
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

func TestAddPerfectSongRehearsal_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.AddPerfectSongRehearsalRequest{
		ID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/perfect-rehearsal", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAddPerfectSongRehearsal_WhenSectionsHaveNoOccurrences_ShouldNotMakeAnyUpdate(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]
	request := requests.AddPerfectSongRehearsalRequest{
		ID: song.ID,
	}

	db := utils.GetDatabase(t)
	db.Preload("Sections").Preload("Sections.History").Find(&song, song.ID)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/perfect-rehearsal", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var newSong model.Song
	db = db.Session(&gorm.Session{NewDB: true})
	db.Preload("Sections").Preload("Sections.History").Find(&newSong, song.ID)

	assert.Equal(t, song, newSong)
}

func TestAddPerfectSongRehearsal_WhenSuccessful_ShouldUpdateSongAndSections(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.AddPerfectSongRehearsalRequest{
		ID: songData.Songs[4].ID,
	}

	var song model.Song
	db := utils.GetDatabase(t)
	db.Preload("Sections").Preload("Sections.History").Find(&song, request.ID)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/perfect-rehearsal", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var newSong model.Song
	db = db.Session(&gorm.Session{NewDB: true})
	db.Preload("Sections").
		Preload("Sections.History", func(db *gorm.DB) *gorm.DB { return db.Order("created_at desc") }).
		Find(&newSong, request.ID)

	assertion.PerfectSongRehearsal(t, song, newSong)
}
