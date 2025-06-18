package album

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/internal/message/topics"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	albumData "repertoire/server/test/integration/test/data/album"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestDeleteImageFromAlbum_WhenAlbumIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/albums/images/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteImageFromAlbum_WhenAlbumHasNoImage_ShouldReturnConflictError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	album := albumData.Albums[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/albums/images/"+album.ID.String())

	// then
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestDeleteImageFromAlbum_WhenSuccessful_ShouldUpdateAlbumAndDeleteImage(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	album := albumData.Albums[1]

	messages := utils.SubscribeToTopic(topics.AlbumsUpdatedTopic)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/albums/images/"+album.ID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)
	db.Find(&album, album.ID)

	assert.Nil(t, album.ImageURL)

	assertion.AssertMessage(t, messages, func(ids []uuid.UUID) {
		assert.Len(t, ids, 1)
		assert.Equal(t, album.ID, ids[0])
	})
}
