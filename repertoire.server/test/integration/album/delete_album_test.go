package album

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
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

			messages := utils.SubscribeToTopic(topics.AlbumsDeletedTopic)

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().DELETE(w, "/api/albums/"+test.album.ID.String())

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			db := utils.GetDatabase(t)

			var deletedAlbum model.Album
			db.Find(&deletedAlbum, test.album.ID)
			assert.Empty(t, deletedAlbum)

			if len(test.album.Songs) > 0 {
				var ids []uuid.UUID
				for _, song := range test.album.Songs {
					ids = append(ids, song.ID)
				}

				var songs []model.Song
				db.Find(&songs, ids)
				assert.NotEmpty(t, songs)
			}

			assertion.AssertMessage(t, messages, func(payloadAlbums []model.Album) {
				assert.Len(t, payloadAlbums, 1)
				assert.Equal(t, test.album.ID, payloadAlbums[0].ID)
			})
		})
	}
}

func TestDeleteAlbum_WhenWithSongs_ShouldDeleteAlbumAndSongs(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	album := albumData.Albums[0]

	messages := utils.SubscribeToTopic(topics.AlbumsDeletedTopic)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/albums/"+album.ID.String()+"?withSongs=true")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var deletedAlbum model.Album
	db.Find(&deletedAlbum, album.ID)
	assert.Empty(t, deletedAlbum)

	var ids []uuid.UUID
	for _, song := range album.Songs {
		ids = append(ids, song.ID)
	}

	var songs []model.Song
	db.Find(&songs, ids)
	assert.Empty(t, songs)

	assertion.AssertMessage(t, messages, func(payloadAlbum model.Album) {
		assert.Equal(t, album.ID, payloadAlbum.ID)
	})
}
