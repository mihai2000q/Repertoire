package section

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

func TestUpdateAllSongSections_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateAllSongSectionsRequest{
		SongID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections/all", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateAllSongSections_WhenSuccessful_ShouldUpdateAllSongSections(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]
	request := requests.UpdateAllSongSectionsRequest{
		SongID:       song.ID,
		InstrumentID: &songData.Users[0].Instruments[0].ID,
		BandMemberID: &songData.Artists[0].BandMembers[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections/all", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var updatedSong model.Song
	db := utils.GetDatabase(t)
	db.Preload("Sections").Find(&updatedSong, request.SongID)

	for _, section := range updatedSong.Sections {
		if request.InstrumentID != nil {
			assert.Equal(t, request.InstrumentID, section.InstrumentID)
		}
		if request.BandMemberID != nil {
			assert.Equal(t, request.BandMemberID, section.BandMemberID)
		}
	}
}
