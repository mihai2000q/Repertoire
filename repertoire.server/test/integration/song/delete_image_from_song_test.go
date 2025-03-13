package song

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/internal/message/topics"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestDeleteImageFromSong_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/images/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteImageFromSong_WhenSongHasNoImage_ShouldReturnBadRequestError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[1]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/images/"+song.ID.String())

	// then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteImageFromSong_WhenSuccessful_ShouldUpdateSongAndDeleteImage(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]

	messages := utils.SubscribeToTopic(topics.SongsUpdatedTopic)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/images/"+song.ID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)
	db.Find(&song, song.ID)

	assert.Nil(t, song.ImageURL)

	assertion.AssertMessage(t, messages, topics.SongsUpdatedTopic, func(ids []uuid.UUID) {
		assert.Len(t, ids, 1)
		assert.Equal(t, song.ID, ids[0])
	})
}
