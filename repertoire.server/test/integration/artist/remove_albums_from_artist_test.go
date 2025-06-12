package artist

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	artistData "repertoire/server/test/integration/test/data/artist"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRemoveAlbumsFromArtist_WhenAlbumArtistIsDifferent_ShouldReturnConflictError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artist := artistData.Artists[0]

	request := requests.RemoveAlbumsFromArtistRequest{
		ID: artist.ID,
		AlbumIDs: []uuid.UUID{
			artist.Albums[3].ID,
			artistData.Artists[1].Albums[0].ID,
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/artists/remove-albums", request)

	// then
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRemoveAlbumsFromArtist_WhenSuccessful_ShouldDeleteAlbumsFromArtist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artist := artistData.Artists[0]
	oldAlbumsLength := len(artist.Albums)

	request := requests.RemoveAlbumsFromArtistRequest{
		ID: artist.ID,
		AlbumIDs: []uuid.UUID{
			artist.Albums[3].ID,
			artist.Albums[1].ID,
		},
	}

	messages := utils.SubscribeToTopic(topics.AlbumsUpdatedTopic)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/artists/remove-albums", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)
	db.Preload("Albums").Find(&artist, artist.ID)
	assertRemoveAlbumsFromArtist(t, request, artist, oldAlbumsLength)

	// all the songs from the albums also lost their artist
	var albums []model.Album
	db.Preload("Songs").Find(&albums, request.AlbumIDs)
	for _, album := range albums {
		assert.Nil(t, album.ArtistID)
		for _, song := range album.Songs {
			assert.Nil(t, song.ArtistID)
		}
	}

	assertion.AssertMessage(t, messages, func(ids []uuid.UUID) {
		assert.Equal(t, request.AlbumIDs, ids)
	})
}

func assertRemoveAlbumsFromArtist(
	t *testing.T,
	request requests.RemoveAlbumsFromArtistRequest,
	artist model.Artist,
	oldAlbumsLength int,
) {
	assert.Equal(t, artist.ID, request.ID)

	assert.Len(t, artist.Albums, oldAlbumsLength-len(request.AlbumIDs))
	for _, album := range artist.Albums {
		assert.NotContains(t, request.AlbumIDs, album.ID)
	}
}
