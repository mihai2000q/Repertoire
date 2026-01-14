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
			songData.Songs[4].ID,
		},
	}

	var songs []model.Song
	db := utils.GetDatabase(t)
	db.Preload("Sections").Preload("Sections.History").Find(&songs, request.IDs)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/perfect-rehearsals", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var newSongs []model.Song
	db = db.Session(&gorm.Session{NewDB: true})
	db.Preload("Sections").
		Preload("Sections.History", func(db *gorm.DB) *gorm.DB { return db.Order("created_at desc") }).
		Find(&newSongs, request.IDs)

	for i := range newSongs {
		totalOccurrences := uint(0)
		for _, section := range newSongs[i].Sections {
			totalOccurrences += section.Occurrences
		}

		if totalOccurrences > 0 {
			assertion.PerfectSongRehearsal(t, songs[i], newSongs[i])
		} else {
			assert.Equal(t, newSongs[i].Rehearsals, songs[i].Rehearsals)
			assert.Equal(t, newSongs[i].Progress, songs[i].Progress)
			assert.Equal(t, newSongs[i].LastTimePlayed, songs[i].LastTimePlayed)
		}
	}
}
