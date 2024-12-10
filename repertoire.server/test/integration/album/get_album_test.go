package album

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	albumData "repertoire/server/test/integration/test/data/album"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestGetAlbum_WhenUserIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/albums/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetAlbum_WhenSuccessful_ShouldReturnAlbum(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	album := albumData.Albums[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/albums/"+album.ID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var responseAlbum model.Album
	_ = json.Unmarshal(w.Body.Bytes(), &responseAlbum)

	db := utils.GetDatabase(t)
	db.Preload("Songs", func(db *gorm.DB) *gorm.DB {
		return db.Order("songs.album_track_no")
	}).
		Preload("Artist").
		Find(&album, album.ID)

	assertion.ResponseAlbum(t, album, responseAlbum, true, true)
}
