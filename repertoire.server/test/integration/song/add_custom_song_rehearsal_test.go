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

func TestAddCustomSongRehearsal_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.AddCustomSongRehearsalRequest{
		ID:            uuid.New(),
		ArrangementID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/custom-rehearsal", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAddCustomSongRehearsal_WhenSongHasNoArrangement_ShouldReturnBadRequestError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.AddCustomSongRehearsalRequest{
		ID:            songData.Songs[0].ID,
		ArrangementID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/custom-rehearsal", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAddCustomSongRehearsal_WhenPartsHaveNoOccurrences_ShouldNotMakeAnyUpdate(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.AddCustomSongRehearsalRequest{
		ID:            songData.SongArrangements[0].SongID,
		ArrangementID: songData.SongArrangements[0].ID,
	}

	var song model.Song
	db := utils.GetDatabase(t)
	db.Preload("Parts").
		Preload("Parts.History").
		Preload("Parts.ArrangementOccurrences").
		Find(&song, song.ID)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/custom-rehearsal", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var newSong model.Song
	db = db.Session(&gorm.Session{NewDB: true})
	db.Preload("Parts").
		Preload("Parts.History").
		Preload("Parts.ArrangementOccurrences").
		Find(&newSong, song.ID)

	assert.Equal(t, song, newSong)
}

func TestAddCustomSongRehearsal_WhenSuccessful_ShouldUpdateSongAndParts(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.AddCustomSongRehearsalRequest{
		ID:            songData.SongArrangements[1].SongID,
		ArrangementID: songData.SongArrangements[1].ID,
	}

	getSongQuery := func(db *gorm.DB, song *model.Song) {
		db.Preload("Parts", func(db *gorm.DB) *gorm.DB { return db.Order("song_parts.song_order") }).
			Preload("Parts.History", func(db *gorm.DB) *gorm.DB {
				return db.
					Where(&model.SongPartHistory{Property: model.RehearsalsProperty}).
					Order("created_at desc")
			}).
			Preload("Parts.ArrangementOccurrences", func(db *gorm.DB) *gorm.DB {
				return db.Joins("LEFT JOIN song_arrangements ON id = arrangement_id").Order("\"order\"")
			}).
			Find(&song, request.ID)
	}

	var song model.Song
	db := utils.GetDatabase(t)
	getSongQuery(db, &song)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/custom-rehearsal", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var newSong model.Song
	db = db.Session(&gorm.Session{NewDB: true})
	getSongQuery(db, &newSong)

	assertion.CustomSongRehearsal(t, song, newSong, request.ArrangementID)
}
