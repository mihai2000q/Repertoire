package album

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
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
)

func TestRemoveSongsFromAlbum_WhenAlbumIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	request := requests.RemoveSongsFromAlbumRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/albums/remove-songs", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRemoveSongsFromAlbum_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	album := albumData.Albums[0]

	request := requests.RemoveSongsFromAlbumRequest{
		ID: album.ID,
		SongIDs: []uuid.UUID{
			album.Songs[3].ID,
			album.Songs[1].ID,
			uuid.New(),
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/albums/remove-songs", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRemoveSongsFromAlbum_WhenSuccessful_ShouldDeleteSongsFromAlbum(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	album := albumData.Albums[0]
	oldSongsLength := len(album.Songs)

	request := requests.RemoveSongsFromAlbumRequest{
		ID: album.ID,
		SongIDs: []uuid.UUID{
			album.Songs[3].ID,
			album.Songs[1].ID,
		},
	}

	messages := utils.SubscribeToTopic(topics.SongsUpdatedTopic)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/albums/remove-songs", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)
	db.Preload("Songs", func(db *gorm.DB) *gorm.DB {
		return db.Order("songs.album_track_no")
	}).Find(&album, album.ID)
	assertRemoveSongsFromAlbum(t, request, album, oldSongsLength)

	assertion.AssertMessage(t, messages, topics.SongsUpdatedTopic, func(ids []uuid.UUID) {
		assert.Equal(t, ids, request.SongIDs)
	})
}

func assertRemoveSongsFromAlbum(
	t *testing.T,
	request requests.RemoveSongsFromAlbumRequest,
	album model.Album,
	oldSongsLength int,
) {
	assert.Equal(t, album.ID, request.ID)

	assert.Len(t, album.Songs, oldSongsLength-len(request.SongIDs))
	for i, song := range album.Songs {
		assert.NotContains(t, request.SongIDs, song.ID)
		assert.Equal(t, uint(i)+1, *song.AlbumTrackNo)
	}
}
