package album

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	albumData "repertoire/server/test/integration/test/data/album"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
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

func TestGetAlbum_WhenSuccessful_ShouldReturnAlbumWithSongsOrderedByTrackNo(t *testing.T) {
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
		return db.Order("album_track_no")
	}).Joins("Artist").
		Find(&album, album.ID)

	assertion.ResponseAlbum(t, album, responseAlbum, true, true)
}

func TestGetAlbum_WhenRequestHasSongsOrderBy_ShouldReturnAlbumAndSongsOrderedByQueryParam(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	album := albumData.Albums[0]
	songsOrder := "title desc"

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/albums/"+album.ID.String()+"?songsOrderBy="+songsOrder)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var responseAlbum model.Album
	_ = json.Unmarshal(w.Body.Bytes(), &responseAlbum)

	db := utils.GetDatabase(t)
	db.Preload("Songs", func(db *gorm.DB) *gorm.DB {
		return db.Order(songsOrder)
	}).
		Joins("Artist").
		Find(&album, album.ID)

	assertion.ResponseAlbum(t, album, responseAlbum, true, true)
}
