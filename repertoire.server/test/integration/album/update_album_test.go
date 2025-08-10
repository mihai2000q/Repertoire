package album

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/internal"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	albumData "repertoire/server/test/integration/test/data/album"
	"repertoire/server/test/integration/test/utils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateAlbum_WhenAlbumIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	request := requests.UpdateAlbumRequest{
		ID:    uuid.New(),
		Title: "New Title",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/albums", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateAlbum_WhenSuccessful_ShouldUpdateAlbum(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	album := albumData.Albums[0]

	request := requests.UpdateAlbumRequest{
		ID:          album.ID,
		Title:       "New Title",
		ReleaseDate: &[]internal.Date{internal.Date(time.Now())}[0],
		ArtistID:    album.ArtistID,
	}

	messages := utils.SubscribeToTopic(topics.AlbumsUpdatedTopic)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/albums", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)
	db.Find(&album, album.ID)

	assertUpdatedAlbum(t, request, album)

	assertion.AssertMessage(t, messages, func(ids []uuid.UUID) {
		assert.Len(t, ids, 1)
		assert.Equal(t, album.ID, ids[0])
	})
}

func TestUpdateAlbum_WhenUpdatingArtist_ShouldUpdateAlbumAndSongs(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	album := albumData.Albums[0]

	request := requests.UpdateAlbumRequest{
		ID:          album.ID,
		Title:       "New Title",
		ReleaseDate: &[]internal.Date{internal.Date(time.Now())}[0],
		ArtistID:    &albumData.Artists[1].ID,
	}

	messages := utils.SubscribeToTopic(topics.AlbumsUpdatedTopic)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/albums", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)
	db.Preload("Songs").Find(&album, album.ID)

	assertUpdatedAlbum(t, request, album)

	for _, song := range album.Songs {
		assert.Equal(t, song.ArtistID, request.ArtistID)
	}

	assertion.AssertMessage(t, messages, func(ids []uuid.UUID) {
		assert.Len(t, ids, 1)
		assert.Equal(t, album.ID, ids[0])
	})
}

func assertUpdatedAlbum(t *testing.T, request requests.UpdateAlbumRequest, album model.Album) {
	assert.Equal(t, request.ID, album.ID)
	assert.Equal(t, request.Title, album.Title)
	assertion.Date(t, request.ReleaseDate, album.ReleaseDate)
	assert.Equal(t, request.ArtistID, album.ArtistID)
}
