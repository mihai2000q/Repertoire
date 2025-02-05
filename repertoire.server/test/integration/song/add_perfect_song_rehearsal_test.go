package song

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"
	"time"
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

	song := songData.Songs[4]
	request := requests.AddPerfectSongRehearsalRequest{
		ID: song.ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/perfect-rehearsal", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var newSong model.Song
	db := utils.GetDatabase(t)
	db.Preload("Sections").Preload("Sections.History").Find(&newSong, request.ID)

	for i, section := range newSong.Sections {
		if section.Occurrences == 0 { // nothing changed
			newSong.Sections[i].History = nil
			assert.Equal(t, song.Sections[i], newSong.Sections[i])
			continue
		}

		assert.Equal(t, section.Rehearsals, song.Sections[i].Rehearsals+song.Sections[i].Occurrences)
		assert.Greater(t, section.RehearsalsScore, song.Sections[i].RehearsalsScore)
		assert.Greater(t, section.Progress, song.Sections[i].Progress)

		assert.NotEmpty(t, section.History[len(section.History)-1].ID)
		assert.Equal(t, song.Sections[i].Rehearsals, section.History[len(section.History)-1].From)
		assert.Equal(t, section.Rehearsals, section.History[len(section.History)-1].To)
		assert.Equal(t, model.RehearsalsProperty, section.History[len(section.History)-1].Property)
	}

	assert.Greater(t, newSong.Rehearsals, song.Rehearsals)
	assert.Greater(t, newSong.Progress, song.Progress)

	assert.NotNil(t, newSong.LastTimePlayed)
	assert.WithinDuration(t, time.Now(), *newSong.LastTimePlayed, 1*time.Minute)
}
