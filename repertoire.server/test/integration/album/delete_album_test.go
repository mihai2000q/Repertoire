package album

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	albumData "repertoire/server/test/integration/test/data/album"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestDeleteAlbum_WhenAlbumIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/albums/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteAlbum_WhenSuccessful_ShouldDeleteAlbum(t *testing.T) {
	tests := []struct {
		name  string
		album model.Album
	}{
		{
			"Without Files",
			albumData.Albums[0],
		},
		{
			"With Image",
			albumData.Albums[1],
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().DELETE(w, "/api/albums/"+test.album.ID.String())

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			db := utils.GetDatabase()

			var deletedAlbum model.Album
			db.Find(&deletedAlbum, test.album.ID)

			assert.Empty(t, deletedAlbum)
		})
	}
}
