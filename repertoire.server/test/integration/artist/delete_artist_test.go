package artist

import (
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	artistData "repertoire/server/test/integration/test/data/artist"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteArtist_WhenArtistIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/artists/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteArtist_WhenSuccessful_ShouldDeleteArtist(t *testing.T) {
	tests := []struct {
		name   string
		artist model.Artist
	}{
		{
			"Without Files",
			artistData.Artists[1],
		},
		{
			"With Image",
			artistData.Artists[0],
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

			messages := utils.SubscribeToTopic(topics.ArtistsDeletedTopic)

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().DELETE(w, "/api/artists/"+test.artist.ID.String())

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			db := utils.GetDatabase(t)

			var deletedArtist model.Artist
			db.Find(&deletedArtist, test.artist.ID)
			assert.Empty(t, deletedArtist)

			if len(test.artist.Albums) > 0 {
				var ids []uuid.UUID
				for _, album := range test.artist.Albums {
					ids = append(ids, album.ID)
				}

				var albums []model.Album
				db.Find(&albums, ids)
				assert.NotEmpty(t, albums)
			}

			if len(test.artist.Songs) > 0 {
				var ids []uuid.UUID
				for _, song := range test.artist.Songs {
					ids = append(ids, song.ID)
				}

				var songs []model.Song
				db.Find(&songs, ids)
				assert.NotEmpty(t, songs)
			}

			assertion.AssertMessage(t, messages, func(payloadArtist model.Artist) {
				assert.Equal(t, test.artist.ID, payloadArtist.ID)
			})
		})
	}
}

func TestDeleteArtist_WhenWithAlbums_ShouldDeleteArtistAndAlbums(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artist := artistData.Artists[1]

	messages := utils.SubscribeToTopic(topics.ArtistsDeletedTopic)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/artists/"+artist.ID.String()+"?withAlbums=true")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var deletedArtist model.Artist
	db.Find(&deletedArtist, artist.ID)
	assert.Empty(t, deletedArtist)

	var ids []uuid.UUID
	for _, album := range artist.Albums {
		ids = append(ids, album.ID)
	}

	var albums []model.Album
	db.Find(&albums, ids)
	assert.Empty(t, albums)

	assertion.AssertMessage(t, messages, func(payloadArtist model.Artist) {
		assert.Equal(t, artist.ID, payloadArtist.ID)
	})
}

func TestDeleteArtist_WhenWithSongs_ShouldDeleteArtistAndSongs(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artist := artistData.Artists[1]

	messages := utils.SubscribeToTopic(topics.ArtistsDeletedTopic)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/artists/"+artist.ID.String()+"?withSongs=true")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var deletedArtist model.Artist
	db.Find(&deletedArtist, artist.ID)
	assert.Empty(t, deletedArtist)

	var ids []uuid.UUID
	for _, song := range artist.Songs {
		ids = append(ids, song.ID)
	}

	var songs []model.Song
	db.Find(&songs, ids)
	assert.Empty(t, songs)

	assertion.AssertMessage(t, messages, func(payloadArtist model.Artist) {
		assert.Equal(t, artist.ID, payloadArtist.ID)
	})
}

func TestDeleteArtist_WhenWithAlbumsAndSongs_ShouldDeleteArtistAndAlbumsAndSongs(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artist := artistData.Artists[1]

	messages := utils.SubscribeToTopic(topics.ArtistsDeletedTopic)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/artists/"+artist.ID.String()+"?withSongs=true&withAlbums=true")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var deletedArtist model.Artist
	db.Find(&deletedArtist, artist.ID)
	assert.Empty(t, deletedArtist)

	var albumIds []uuid.UUID
	for _, song := range artist.Albums {
		albumIds = append(albumIds, song.ID)
	}

	var albums []model.Album
	db.Find(&albums, albumIds)
	assert.Empty(t, albums)

	var songIds []uuid.UUID
	for _, song := range artist.Songs {
		songIds = append(songIds, song.ID)
	}

	var songs []model.Song
	db.Find(&songs, songIds)
	assert.Empty(t, songs)

	assertion.AssertMessage(t, messages, func(payloadArtist model.Artist) {
		assert.Equal(t, artist.ID, payloadArtist.ID)
	})
}
