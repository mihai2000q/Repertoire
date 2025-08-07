package album

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	albumData "repertoire/server/test/integration/test/data/album"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBulkDeleteAlbums_WhenAlbumsAreNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	request := requests.BulkDeleteAlbumsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/albums/bulk-delete", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBulkDeleteAlbums_WhenSuccessful_ShouldDeleteAlbums(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	messages := utils.SubscribeToTopic(topics.AlbumsDeletedTopic)

	albums := []model.Album{albumData.Albums[0], albumData.Albums[1]}
	request := requests.BulkDeleteAlbumsRequest{
		IDs:       []uuid.UUID{albumData.Albums[0].ID, albumData.Albums[1].ID},
		WithSongs: false,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/albums/bulk-delete", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var deletedAlbum []model.Album
	db.Find(&deletedAlbum, request.IDs)
	assert.Empty(t, deletedAlbum)

	for _, album := range albums {
		if len(album.Songs) > 0 {
			var ids []uuid.UUID
			for _, song := range album.Songs {
				ids = append(ids, song.ID)
			}

			var songs []model.Song
			db.Find(&songs, ids)
			assert.NotEmpty(t, songs)
		}
	}

	assertion.AssertMessage(t, messages, func(payloadAlbums []model.Album) {
		assert.Len(t, payloadAlbums, len(albums))
		for i := range payloadAlbums {
			assert.Contains(t, request.IDs, payloadAlbums[i].ID)
		}
	})
}

func TestBulkDeleteAlbums_WhenWithSongs_ShouldBulkDeleteAlbumsAndSongs(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	messages := utils.SubscribeToTopic(topics.AlbumsDeletedTopic)

	albums := []model.Album{albumData.Albums[0], albumData.Albums[1]}
	request := requests.BulkDeleteAlbumsRequest{
		IDs:       []uuid.UUID{albumData.Albums[0].ID, albumData.Albums[1].ID},
		WithSongs: true,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/albums/bulk-delete", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var deletedAlbum []model.Album
	db.Find(&deletedAlbum, request.IDs)
	assert.Empty(t, deletedAlbum)

	for _, album := range albums {
		var ids []uuid.UUID
		for _, song := range album.Songs {
			ids = append(ids, song.ID)
		}

		if len(ids) == 0 {
			continue
		}

		var songs []model.Song
		db.Find(&songs, ids)
		assert.Empty(t, songs)
	}

	assertion.AssertMessage(t, messages, func(payloadAlbums []model.Album) {
		assert.Len(t, payloadAlbums, len(albums))
		for i := range payloadAlbums {
			assert.Contains(t, request.IDs, payloadAlbums[i].ID)
		}
	})
}
