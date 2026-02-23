package arrangement

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateDefaultSongArrangement_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateDefaultSongArrangementRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/arrangements/default", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateDefaultSongArrangement_WhenArrangementIsNotFound_ShouldReturnInternalServerError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateDefaultSongArrangementRequest{
		ID:     uuid.New(),
		SongID: songData.SongArrangements[2].SongID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/arrangements/default", request)

	// then
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateDefaultSongArrangement_WhenSuccessful_ShouldUpdateDefaultArrangementOnSong(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateDefaultSongArrangementRequest{
		ID:     songData.SongArrangements[2].ID,
		SongID: songData.SongArrangements[2].SongID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/arrangements/default", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var newSong model.Song
	db.Find(&newSong, &model.Song{ID: request.SongID})
	assert.Equal(t, &request.ID, newSong.DefaultArrangementID)
}
