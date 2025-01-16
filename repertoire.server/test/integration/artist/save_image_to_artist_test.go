package artist

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"repertoire/server/test/integration/test/core"
	artistData "repertoire/server/test/integration/test/data/artist"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSaveImageFromArtist_WhenArtistIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
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

func TestSaveImageFromArtist_WhenSuccessful_ShouldUpdateArtistAndSaveImage(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artist := artistData.Artists[0]

	var requestBody bytes.Buffer
	multiWriter := multipart.NewWriter(&requestBody)
	utils.AttachFileToMultipartBody("test-file.jpeg", "image", multiWriter)
	_ = multiWriter.WriteField("id", artist.ID.String())
	_ = multiWriter.Close()

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUTForm(w, "/api/artists/images", &requestBody, multiWriter.FormDataContentType())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)
	db.Find(&artist, artist.ID)

	assert.NotNil(t, artist.ImageURL)
}
