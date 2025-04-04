package artist

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"repertoire/server/internal/message/topics"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	artistData "repertoire/server/test/integration/test/data/artist"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSaveImageToArtist_WhenArtistIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	var requestBody bytes.Buffer
	multiWriter := multipart.NewWriter(&requestBody)
	utils.AttachFileToMultipartBody("test-file.jpeg", "image", multiWriter)
	_ = multiWriter.WriteField("id", uuid.New().String())
	_ = multiWriter.Close()

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUTForm(w, "/api/artists/images", &requestBody, multiWriter.FormDataContentType())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSaveImageToArtist_WhenSuccessful_ShouldUpdateArtistAndSaveImage(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artist := artistData.Artists[0]

	var requestBody bytes.Buffer
	multiWriter := multipart.NewWriter(&requestBody)
	utils.AttachFileToMultipartBody("test-file.jpeg", "image", multiWriter)
	_ = multiWriter.WriteField("id", artist.ID.String())
	_ = multiWriter.Close()

	messages := utils.SubscribeToTopic(topics.ArtistUpdatedTopic)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUTForm(w, "/api/artists/images", &requestBody, multiWriter.FormDataContentType())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)
	db.Find(&artist, artist.ID)

	assert.NotNil(t, artist.ImageURL)

	assertion.AssertMessage(t, messages, func(id uuid.UUID) {
		assert.Equal(t, artist.ID, id)
	})
}
