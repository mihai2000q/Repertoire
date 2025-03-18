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
)

func TestUpdateSongSettings_WhenSettingsIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateSongSettingsRequest{
		SettingsID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/settings", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateSongSettings_WhenSuccessful_ShouldUpdateSongSettings(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[1]
	request := requests.UpdateSongSettingsRequest{
		SettingsID:          song.Settings.ID,
		DefaultInstrumentID: &songData.Users[0].Instruments[0].ID,
		DefaultBandMemberID: &songData.Artists[0].BandMembers[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/settings", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var updatedSettings model.SongSettings
	db := utils.GetDatabase(t)
	db.Find(&updatedSettings, request.SettingsID)

	assertUpdatedSettings(t, request, updatedSettings)
}

func assertUpdatedSettings(t *testing.T, request requests.UpdateSongSettingsRequest, settings model.SongSettings) {
	assert.Equal(t, request.SettingsID, settings.ID)
	assert.Equal(t, request.DefaultInstrumentID, request.DefaultInstrumentID)
	assert.Equal(t, request.DefaultBandMemberID, request.DefaultBandMemberID)
}
