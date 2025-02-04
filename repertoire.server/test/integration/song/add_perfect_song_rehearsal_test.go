package song

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
	core.NewTestHandler().PUT(w, "/api/songs/perfect-rehearsal", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAddPerfectSongRehearsal_WhenSuccessful_ShouldUpdateSong(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]
	request := requests.AddPerfectSongRehearsalRequest{
		ID: song.ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/perfect-rehearsal", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var newSong model.Song
	db := utils.GetDatabase(t)
	db.Preload("Sections").Find(&newSong, request.ID)

	for i, section := range song.Sections {
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

	assert.NotNil(t, song.LastTimePlayed)
	assert.WithinDuration(t, time.Now(), *song.LastTimePlayed, 1*time.Minute)
}
