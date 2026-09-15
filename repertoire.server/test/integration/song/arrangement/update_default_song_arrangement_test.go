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
		ID:     nil,
		SongID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/arrangements/default", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateDefaultSongArrangement_WhenArrangementIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateDefaultSongArrangementRequest{
		ID:     &[]uuid.UUID{uuid.New()}[0],
		SongID: songData.SongArrangements[2].SongID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/arrangements/default", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateDefaultSongArrangement_WhenSuccessfulWithArrangement_ShouldUpdateDefaultArrangementOnSongWithTheArrangement(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateDefaultSongArrangementRequest{
		ID:     &[]uuid.UUID{songData.SongArrangements[2].ID}[0],
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
	assert.Equal(t, request.ID, newSong.DefaultArrangementID)
}

func TestUpdateDefaultSongArrangement_WhenSuccessfulWithNullID_ShouldUpdateDefaultArrangementOnSongWithNull(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateDefaultSongArrangementRequest{
		ID:     nil,
		SongID: songData.SongArrangements[0].SongID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/arrangements/default", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var newSong model.Song
	db.Find(&newSong, &model.Song{ID: request.SongID})
	assert.Nil(t, newSong.DefaultArrangementID)
}
