package album

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"repertoire/server/test/integration/test/core"
	albumData "repertoire/server/test/integration/test/data/album"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestSaveImageFromAlbum_WhenAlbumIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	var requestBody bytes.Buffer
	multiWriter := multipart.NewWriter(&requestBody)
	utils.AttachFileToMultipartBody("test-file.jpeg", "image", multiWriter)
	_ = multiWriter.WriteField("id", uuid.New().String())
	_ = multiWriter.Close()

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUTForm(w, "/api/albums/images", &requestBody, multiWriter.FormDataContentType())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSaveImageFromAlbum_WhenSuccessful_ShouldUpdateAlbumAndSaveImage(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	album := albumData.Albums[0]

	var requestBody bytes.Buffer
	multiWriter := multipart.NewWriter(&requestBody)
	utils.AttachFileToMultipartBody("test-file.jpeg", "image", multiWriter)
	_ = multiWriter.WriteField("id", album.ID.String())
	_ = multiWriter.Close()

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUTForm(w, "/api/albums/images", &requestBody, multiWriter.FormDataContentType())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)
	db.Find(&album, album.ID)

	assert.NotNil(t, album.ImageURL)
}
