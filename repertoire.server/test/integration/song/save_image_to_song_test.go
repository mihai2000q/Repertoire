package song

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"repertoire/server/internal/message/topics"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestSaveImageFromSong_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	var requestBody bytes.Buffer
	multiWriter := multipart.NewWriter(&requestBody)
	utils.AttachFileToMultipartBody("test-file.jpeg", "image", multiWriter)
	_ = multiWriter.WriteField("id", uuid.New().String())
	_ = multiWriter.Close()

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUTForm(w, "/api/songs/images", &requestBody, multiWriter.FormDataContentType())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSaveImageFromSong_WhenSuccessful_ShouldUpdateSongAndSaveImage(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]

	messages := utils.SubscribeToTopic(topics.SongsUpdatedTopic)

	var requestBody bytes.Buffer
	multiWriter := multipart.NewWriter(&requestBody)
	utils.AttachFileToMultipartBody("test-file.jpeg", "image", multiWriter)
	_ = multiWriter.WriteField("id", song.ID.String())
	_ = multiWriter.Close()

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUTForm(w, "/api/songs/images", &requestBody, multiWriter.FormDataContentType())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)
	db.Find(&song, song.ID)

	assert.NotNil(t, song.ImageURL)

	assertion.AssertMessage(t, messages, topics.SongsUpdatedTopic, func(ids []uuid.UUID) {
		assert.Len(t, ids, 1)
		assert.Equal(t, song.ID, ids[0])
	})
}
