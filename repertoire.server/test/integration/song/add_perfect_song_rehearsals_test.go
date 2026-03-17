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

func TestAddPerfectSongRehearsals_WhenSongsAreNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.AddPerfectSongRehearsalsRequest{
		IDs: []uuid.UUID{
			uuid.New(),
			uuid.New(),
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/perfect-rehearsals", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAddPerfectSongRehearsals_WhenSuccessful_ShouldUpdateSongAndSectionsIfTheyHaveOccurrences(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.AddPerfectSongRehearsalsRequest{
		IDs: []uuid.UUID{
			songData.Songs[0].ID,
			songData.Songs[1].ID,
			songData.Songs[4].ID,
			songData.Songs[5].ID,
		},
	}

	getSongsQuery := func(db *gorm.DB, songs *[]model.Song) {
		db.Preload("Sections", func(db *gorm.DB) *gorm.DB { return db.Order("song_sections.order") }).
			Preload("Sections.History", func(db *gorm.DB) *gorm.DB {
				return db.
					Where(&model.SongSectionHistory{Property: model.RehearsalsProperty}).
					Order("created_at desc")
			}).
			Preload("Sections.ArrangementOccurrences", func(db *gorm.DB) *gorm.DB {
				return db.Joins("LEFT JOIN song_arrangements ON id = arrangement_id").Order("\"order\"")
			}).
			Find(&songs, request.IDs)
	}

	var songs []model.Song
	db := utils.GetDatabase(t)
	getSongsQuery(db, &songs)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/perfect-rehearsals", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var newSongs []model.Song
	db = db.Session(&gorm.Session{NewDB: true})
	getSongsQuery(db, &newSongs)

	for i := range newSongs {
		assertion.PerfectSongRehearsal(t, songs[i], newSongs[i])
	}
}
