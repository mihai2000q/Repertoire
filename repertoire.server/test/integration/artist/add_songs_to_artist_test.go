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
	"gorm.io/gorm"
)

func TestAddSongsToArtist_WhenArtistIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	request := requests.AddSongsToArtistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/artists/add-songs", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAddSongsToArtist_WhenSongsAreNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	request := requests.AddSongsToArtistRequest{
		ID:      artistData.Artists[0].ID,
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/artists/add-songs", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAddSongsToArtist_WhenSongAlreadyHasArtist_ShouldReturnConflictError(t *testing.T) {
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
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestAddSongsToArtist_WhenSuccessful_ShouldAddSongsToArtist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artist := artistData.Artists[0]
	db := utils.GetDatabase(t) // too many nested songs on albums on artist and so on...
	var oldSongsLength int64
	db.Model(&model.Song{}).Where(&model.Song{ArtistID: &artist.ID}).Count(&oldSongsLength)

	request := requests.AddSongsToArtistRequest{
		ID: artist.ID,
		SongIDs: []uuid.UUID{
			artistData.Songs[0].ID,
			artistData.Songs[1].ID,
			artistData.Albums[1].Songs[0].ID,
			artistData.Albums[1].Songs[1].ID,
		},
	}

	songMessages := utils.SubscribeToTopic(topics.SongsUpdatedTopic)
	albumMessages := utils.SubscribeToTopic(topics.AlbumsUpdatedTopic)

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

	assertion.AssertMessage(t, songMessages, func(ids []uuid.UUID) {
		var songIDs []uuid.UUID
		for _, id := range request.SongIDs {
			songIDs = append(songIDs, id)
		}
		for _, albumSong := range albumSongs {
			songIDs = slices.DeleteFunc(songIDs, func(id uuid.UUID) bool { return id == albumSong.ID })
		}
		assert.ElementsMatch(t, songIDs, ids)
	})
	assertion.AssertMessage(t, albumMessages, func(ids []uuid.UUID) {
		var albumIDs []uuid.UUID
		albumsSet := make(map[uuid.UUID]bool)
		for _, song := range albumSongs {
			if !albumsSet[song.ID] {
				albumIDs = append(albumIDs, *song.AlbumID)
				albumsSet[*song.AlbumID] = true
			}
		}
		assert.ElementsMatch(t, albumIDs, ids)
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
	albumsVisited := make(map[uuid.UUID]bool)
	for _, s := range albumSongs {
		if albumsVisited[s.Album.ID] {
			continue
		}
		for _, albumSong := range s.Album.Songs {
			nestedSongIDs = append(nestedSongIDs, albumSong.ID)
		}
		albumsVisited[s.Album.ID] = true
	}

	assert.Equal(t, artist.ID, request.ID)

	totalNewSongIDs := slices.Clone(nestedSongIDs)
	for _, songID := range request.SongIDs {
		if !slices.Contains(totalNewSongIDs, songID) {
			totalNewSongIDs = append(totalNewSongIDs, songID)
		}
	}

	totalSongsLen := int(oldSongsLength) + len(totalNewSongIDs)
	assert.Len(t, artist.Songs, totalSongsLen)

	var allSongIDs []uuid.UUID
	for _, song := range artist.Songs {
		allSongIDs = append(allSongIDs, song.ID)
		assert.Equal(t, artist.ID, *song.ArtistID)
		if song.Album != nil {
			assert.Equal(t, artist.ID, *song.Album.ArtistID)
		}
	}
	assert.Subset(t, allSongIDs, request.SongIDs)
	assert.Subset(t, allSongIDs, nestedSongIDs)
	assert.Subset(t, allSongIDs, totalNewSongIDs)
}
