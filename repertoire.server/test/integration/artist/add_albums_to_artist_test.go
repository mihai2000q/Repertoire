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

func TestAddAlbumsToArtist_WhenAlbumAlreadyHasArtist_ShouldReturnBadRequestError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artist := artistData.Artists[0]

	request := requests.AddAlbumsToArtistRequest{
		ID: artist.ID,
		AlbumIDs: []uuid.UUID{
			artist.Albums[0].ID,
			artistData.Albums[1].ID,
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/artists/add-albums", request)

	// then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddAlbumsToArtist_WhenSuccessful_ShouldAddAlbumsToArtist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artist := artistData.Artists[0]
	oldAlbumsLength := len(artist.Albums)

	request := requests.AddAlbumsToArtistRequest{
		ID: artist.ID,
		AlbumIDs: []uuid.UUID{
			artistData.Albums[0].ID,
			artistData.Albums[1].ID,
		},
	}

	messages := utils.SubscribeToTopic(topics.AlbumsUpdatedTopic)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/artists/add-albums", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)
	db.Preload("Albums").Preload("Albums.Songs").Find(&artist, artist.ID)
	assertAddedAlbumsToArtist(t, request, artist, oldAlbumsLength)

	assertion.AssertMessage(t, messages, func(ids []uuid.UUID) {
		assert.Equal(t, request.AlbumIDs, ids)
	})
}

func assertAddedAlbumsToArtist(
	t *testing.T,
	request requests.AddAlbumsToArtistRequest,
	artist model.Artist,
	oldAlbumsLength int,
) {
	assert.Equal(t, artist.ID, request.ID)

	assert.Len(t, artist.Albums, oldAlbumsLength+len(request.AlbumIDs))
	var albumIDs []uuid.UUID
	for _, album := range artist.Albums {
		albumIDs = append(albumIDs, album.ID)
		assert.Equal(t, artist.ID, *album.ArtistID)

		for _, song := range album.Songs {
			assert.Equal(t, artist.ID, *song.ArtistID)
		}
	}
	assert.Subset(t, albumIDs, request.AlbumIDs)
}
