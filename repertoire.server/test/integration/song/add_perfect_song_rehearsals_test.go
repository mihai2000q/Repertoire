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

	songs := []model.Song{
		songData.Songs[0],
		songData.Songs[4],
	}
	request := requests.AddPerfectSongRehearsalsRequest{
		IDs: []uuid.UUID{
			songs[0].ID,
			songs[1].ID,
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/perfect-rehearsals", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var newSongs []model.Song
	db := utils.GetDatabase(t)
	db.Preload("Sections").Preload("Sections.History").Find(&newSongs, request.IDs)

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
