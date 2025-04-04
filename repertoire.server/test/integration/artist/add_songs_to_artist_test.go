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
	"gorm.io/gorm"
)

func TestAddSongsToArtist_WhenSongAlreadyHasArtist_ShouldReturnBadRequestError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artist := artistData.Artists[0]

	request := requests.AddSongsToArtistRequest{
		ID: artist.ID,
		SongIDs: []uuid.UUID{
			artistData.Songs[0].ID,
			artistData.Songs[3].ID, // already has an artist
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/artists/add-songs", request)

	// then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddSongsToArtist_WhenSuccessful_ShouldAddSongsToArtist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artist := artistData.Artists[0]
	db := utils.GetDatabase(t) // too many nested songs on albums on artist and so on...
	var oldSongsLength int64
	db.Model(&model.Song{}).Where("artist_id = ?", artist.ID).Count(&oldSongsLength)

	request := requests.AddSongsToArtistRequest{
		ID: artist.ID,
		SongIDs: []uuid.UUID{
			artistData.Songs[0].ID,
			artistData.Songs[1].ID,
			artistData.Albums[1].Songs[0].ID,
		},
	}

	messages := utils.SubscribeToTopic(topics.SongsUpdatedTopic)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/artists/add-songs", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db.Session(&gorm.Session{NewDB: true}).
		Preload("Songs").
		Preload("Songs.Album").
		Find(&artist, artist.ID)

	var albumSongs []model.Song
	db.Preload("Album").
		Preload("Album.Songs").
		Where("album_id IS NOT NULL").
		Find(&albumSongs, request.SongIDs)

	assertAddedSongsToArtist(t, request, artist, oldSongsLength, albumSongs)

	assertion.AssertMessage(t, messages, func(ids []uuid.UUID) {
		assert.Equal(t, request.SongIDs, ids)
	})
}

func assertAddedSongsToArtist(
	t *testing.T,
	request requests.AddSongsToArtistRequest,
	artist model.Artist,
	oldSongsLength int64,
	albumSongs []model.Song,
) {
	// calculate the number of nested songs, as the whole album is added to artist
	// (including other non-mentioned songs in the request)
	var nestedSongIDs []uuid.UUID
	for _, s := range albumSongs {
		for _, albumSong := range s.Album.Songs {
			if albumSong.ID == s.ID {
				continue
			}
			nestedSongIDs = append(nestedSongIDs, albumSong.ID)
		}
	}

	assert.Equal(t, artist.ID, request.ID)

	totalSongs := int(oldSongsLength) + len(request.SongIDs) + len(nestedSongIDs)
	assert.Len(t, artist.Songs, totalSongs)

	var songIDs []uuid.UUID
	for _, song := range artist.Songs {
		songIDs = append(songIDs, song.ID)
		assert.Equal(t, artist.ID, *song.ArtistID)
		if song.Album != nil {
			assert.Equal(t, artist.ID, *song.Album.ArtistID)
		}
	}
	assert.Subset(t, songIDs, request.SongIDs)
	assert.Subset(t, songIDs, nestedSongIDs)
}
