package song

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"
	"time"
)

func TestUpdateSong_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateSongRequest{
		ID:    uuid.New(),
		Title: "New Title",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateSong_WhenSuccessful_ShouldUpdateSong(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateSongRequest{
		ID:          songData.Songs[0].ID,
		Title:       "New Title",
		ReleaseDate: &[]time.Time{time.Now()}[0],
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var song model.Song
	db := utils.GetDatabase(t)
	db.Find(&song, request.ID)

	assertUpdatedSong(t, request, song)
}

func assertUpdatedSong(t *testing.T, request requests.UpdateSongRequest, song model.Song) {
	assert.Equal(t, request.ID, song.ID)
	assert.Equal(t, request.Title, song.Title)
	assert.Equal(t, request.Description, song.Description)
	assert.Equal(t, request.IsRecorded, song.IsRecorded)
	assert.Equal(t, request.Bpm, song.Bpm)
	assert.Equal(t, request.SongsterrLink, song.SongsterrLink)
	assert.Equal(t, request.YoutubeLink, song.YoutubeLink)
	assertion.Time(t, request.ReleaseDate, song.ReleaseDate)
	assert.Equal(t, request.Difficulty, song.Difficulty)
	assert.Equal(t, request.GuitarTuningID, song.GuitarTuningID)
}
