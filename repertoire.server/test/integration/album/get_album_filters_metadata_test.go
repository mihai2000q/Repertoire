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

	"github.com/stretchr/testify/assert"
)

func TestGetAlbumFiltersMetadata_WhenSuccessful_ShouldReturnAlbumFiltersMetadata(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/albums/filters-metadata")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var metadata model.AlbumFiltersMetadata
	_ = json.Unmarshal(w.Body.Bytes(), &metadata)

	var albums []model.Album
	db := utils.GetDatabase(t)
	db.Preload("Songs").Find(&albums)

	assertion.AlbumFiltersMetadata(t, metadata, albums)
}
