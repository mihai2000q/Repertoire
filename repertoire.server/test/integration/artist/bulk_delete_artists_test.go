package artist

import (
	"github.com/google/uuid"
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

	"github.com/stretchr/testify/assert"
)

func TestBulkDeleteArtists_WhenArtistsAreNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	request := requests.BulkDeleteArtistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/artists/bulk-delete", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBulkDeleteArtists_WhenSuccessful_ShouldDeleteArtists(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artists := []model.Artist{artistData.Artists[0], artistData.Artists[1]}
	request := requests.BulkDeleteArtistsRequest{
		IDs: []uuid.UUID{artistData.Artists[0].ID, artistData.Artists[1].ID},
	}

	messages := utils.SubscribeToTopic(topics.ArtistsDeletedTopic)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/artists/bulk-delete", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var deletedArtists []model.Artist
	db.Find(&deletedArtists, request.IDs)
	assert.Empty(t, deletedArtists)

	for _, art := range artists {
		if len(art.Albums) > 0 {
			var ids []uuid.UUID
			for _, album := range art.Albums {
				ids = append(ids, album.ID)
			}

			var albums []model.Album
			db.Find(&albums, ids)
			assert.NotEmpty(t, albums)
		}

		if len(art.Songs) > 0 {
			var ids []uuid.UUID
			for _, song := range art.Songs {
				ids = append(ids, song.ID)
			}

			var songs []model.Song
			db.Find(&songs, ids)
			assert.NotEmpty(t, songs)
		}
	}

	assertion.AssertMessage(t, messages, func(payloadArtists []model.Artist) {
		assert.Len(t, payloadArtists, len(artists))
		for _, art := range payloadArtists {
			assert.Contains(t, request.IDs, art.ID)
		}
	})
}

func TestBulkDeleteArtists_WhenWithAlbums_ShouldDeleteArtistsAndAlbums(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artists := []model.Artist{artistData.Artists[0], artistData.Artists[1]}
	request := requests.BulkDeleteArtistsRequest{
		IDs:        []uuid.UUID{artistData.Artists[0].ID, artistData.Artists[1].ID},
		WithAlbums: true,
	}

	messages := utils.SubscribeToTopic(topics.ArtistsDeletedTopic)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/artists/bulk-delete", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var deletedArtists []model.Artist
	db.Find(&deletedArtists, request.IDs)
	assert.Empty(t, deletedArtists)

	for _, art := range artists {
		var ids []uuid.UUID
		for _, album := range art.Albums {
			ids = append(ids, album.ID)
		}

		if len(ids) == 0 {
			continue
		}

		var albums []model.Album
		db.Find(&albums, ids)
		assert.Empty(t, albums)
	}

	assertion.AssertMessage(t, messages, func(payloadArtists []model.Artist) {
		assert.Len(t, payloadArtists, len(artists))
		for _, art := range payloadArtists {
			assert.Contains(t, request.IDs, art.ID)
		}
	})
}

func TestBulkDeleteArtists_WhenWithSongs_ShouldDeleteArtistsAndSongs(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artists := []model.Artist{artistData.Artists[0], artistData.Artists[1]}
	request := requests.BulkDeleteArtistsRequest{
		IDs:       []uuid.UUID{artistData.Artists[0].ID, artistData.Artists[1].ID},
		WithSongs: true,
	}

	messages := utils.SubscribeToTopic(topics.ArtistsDeletedTopic)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/artists/bulk-delete", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var deletedArtists []model.Artist
	db.Find(&deletedArtists, request.IDs)
	assert.Empty(t, deletedArtists)

	for _, art := range artists {
		var ids []uuid.UUID
		for _, song := range art.Songs {
			ids = append(ids, song.ID)
		}

		if len(ids) == 0 {
			continue
		}

		var songs []model.Song
		db.Find(&songs, ids)
		assert.Empty(t, songs)
	}

	assertion.AssertMessage(t, messages, func(payloadArtists []model.Artist) {
		assert.Len(t, payloadArtists, len(artists))
		for _, art := range payloadArtists {
			assert.Contains(t, request.IDs, art.ID)
		}
	})
}

func TestBulkDeleteArtists_WhenWithAlbumsAndSongs_ShouldDeleteArtistsAndAlbumsAndSongs(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artists := []model.Artist{artistData.Artists[0], artistData.Artists[1]}
	request := requests.BulkDeleteArtistsRequest{
		IDs:        []uuid.UUID{artistData.Artists[0].ID, artistData.Artists[1].ID},
		WithSongs:  true,
		WithAlbums: true,
	}

	messages := utils.SubscribeToTopic(topics.ArtistsDeletedTopic)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/artists/bulk-delete", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var deletedArtists []model.Artist
	db.Find(&deletedArtists, request.IDs)
	assert.Empty(t, deletedArtists)

	for _, art := range artists {
		var albumIds []uuid.UUID
		for _, song := range art.Albums {
			albumIds = append(albumIds, song.ID)
		}

		if len(albumIds) > 0 {
			var albums []model.Album
			db.Find(&albums, albumIds)
			assert.Empty(t, albums)
		}

		var songIds []uuid.UUID
		for _, song := range art.Songs {
			songIds = append(songIds, song.ID)
		}

		if len(songIds) > 0 {
			var songs []model.Song
			db.Find(&songs, songIds)
			assert.Empty(t, songs)
		}
	}

	assertion.AssertMessage(t, messages, func(payloadArtists []model.Artist) {
		assert.Len(t, payloadArtists, len(artists))
		for _, art := range payloadArtists {
			assert.Contains(t, request.IDs, art.ID)
		}
	})
}
