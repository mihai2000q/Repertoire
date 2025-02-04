package member

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

func TestSaveImageToBandMember_WhenMemberIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	var requestBody bytes.Buffer
	multiWriter := multipart.NewWriter(&requestBody)
	utils.AttachFileToMultipartBody("test-file.jpeg", "image", multiWriter)
	_ = multiWriter.WriteField("id", uuid.New().String())
	_ = multiWriter.Close()

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUTForm(w, "/api/artists/band-members/images", &requestBody, multiWriter.FormDataContentType())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSaveImageToBandMember_WhenSuccessful_ShouldUpdateBandMemberAndSaveImage(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	bandMember := artistData.Artists[0].BandMembers[0]

	var requestBody bytes.Buffer
	multiWriter := multipart.NewWriter(&requestBody)
	utils.AttachFileToMultipartBody("test-file.jpeg", "image", multiWriter)
	_ = multiWriter.WriteField("id", bandMember.ID.String())
	_ = multiWriter.Close()

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUTForm(w, "/api/artists/band-members/images", &requestBody, multiWriter.FormDataContentType())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)
	db.Find(&bandMember, bandMember.ID)

	assert.NotNil(t, bandMember.ImageURL)
}
