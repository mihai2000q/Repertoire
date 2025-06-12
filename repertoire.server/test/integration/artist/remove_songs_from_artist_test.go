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
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRemoveSongsFromArtist_WhenSongArtistIsDifferent_ShouldReturnConflictError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artist := artistData.Artists[0]

	request := requests.RemoveSongsFromArtistRequest{
		ID: artist.ID,
		SongIDs: []uuid.UUID{
			artistData.Songs[4].ID,
			artistData.Artists[1].Songs[0].ID,
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/artists/remove-songs", request)

	// then
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRemoveSongsFromArtist_WhenSuccessful_ShouldDeleteSongsFromArtist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artist := artistData.Artists[1]
	oldSongsLength := len(slices.DeleteFunc(slices.Clone(artistData.Songs), func(s model.Song) bool {
		return s.ArtistID != nil && *s.ArtistID != artist.ID
	}))

	request := requests.RemoveSongsFromArtistRequest{
		ID: artist.ID,
		SongIDs: []uuid.UUID{
			artist.Songs[2].ID,
			artist.Songs[0].ID,
		},
	}

	messages := utils.SubscribeToTopic(topics.SongsUpdatedTopic)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/artists/remove-songs", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)
	db.Preload("Songs").Find(&artist, artist.ID)
	assertRemoveSongsFromArtist(t, request, artist, oldSongsLength)

	assertion.AssertMessage(t, messages, func(ids []uuid.UUID) {
		assert.Equal(t, request.SongIDs, ids)
	})
}

func assertRemoveSongsFromArtist(
	t *testing.T,
	request requests.RemoveSongsFromArtistRequest,
	artist model.Artist,
	oldSongsLength int,
) {
	assert.Equal(t, artist.ID, request.ID)

	assert.Len(t, artist.Songs, oldSongsLength-len(request.SongIDs))
	for _, song := range artist.Songs {
		assert.NotContains(t, request.SongIDs, song.ID)
	}
}
