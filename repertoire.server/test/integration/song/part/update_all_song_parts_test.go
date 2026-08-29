package part

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

func TestUpdateAllSongParts_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateAllSongPartsRequest{
		SongID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/parts/all", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateAllSongParts_WhenSuccessful_ShouldUpdateAllSongParts(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]
	request := requests.UpdateAllSongPartsRequest{
		SongID:       song.ID,
		InstrumentID: &songData.Users[0].Instruments[0].ID,
		BandMemberID: &songData.Artists[0].BandMembers[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/parts/all", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var updatedSong model.Song
	db := utils.GetDatabase(t)
	db.Preload("Parts").Find(&updatedSong, request.SongID)

	for _, part := range updatedSong.Parts {
		if request.InstrumentID != nil {
			assert.Equal(t, request.InstrumentID, part.InstrumentID)
		}
		if request.BandMemberID != nil {
			assert.Equal(t, request.BandMemberID, part.BandMemberID)
		}
	}
}
