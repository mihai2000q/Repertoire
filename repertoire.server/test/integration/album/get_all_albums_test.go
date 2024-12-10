package album

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	albumData "repertoire/server/test/integration/test/data/album"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestGetAllAlbums_WhenSuccessful_ShouldReturnAlbums(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/albums")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var responseAlbums []model.Album
	_ = json.Unmarshal(w.Body.Bytes(), &responseAlbums)

	db := utils.GetDatabase(t)

	var albums []model.Album
	db.Preload("Artist").Find(&albums)

	for i := range responseAlbums {
		assertion.ResponseAlbum(t, albums[i], responseAlbums[i], true, false)
	}
}
