package artist

import (
	"repertoire/server/internal"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	artistData "repertoire/server/test/integration/test/data/artist"
	"repertoire/server/test/integration/test/utils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestArtistsDeleted_WhenSuccessful_ShouldPublishMessages(t *testing.T) {
	tests := []struct {
		name    string
		artists []model.Artist
	}{
		{
			"Multiple artists",
			[]model.Artist{
				{ID: uuid.New(), Name: "Artist 1", ImageURL: &[]internal.FilePath{"something.png"}[0]},
				{ID: uuid.New(), Name: "Artist 2"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupSearchData(t, artistData.GetSearchDocuments())

			deleteMessages := utils.SubscribeToTopic(topics.DeleteFromSearchEngineTopic)
			deleteStorageMessages := utils.SubscribeToTopic(topics.DeleteDirectoriesStorageTopic)

			// when
			err := utils.PublishToTopic(topics.ArtistsDeletedTopic, test.artists)

			// then
			assert.NoError(t, err)

			// then - assert data
			assertion.AssertMessage(t, deleteMessages, func(ids []string) {
				assert.Len(t, ids, len(test.artists))

				for i := range ids {
					assertion.ArtistSearchID(t, test.artists[i].ID, ids[i])
				}
			})

			assertion.AssertMessage(t, deleteStorageMessages, func(paths []string) {
				assert.Len(t, paths, len(test.artists))
			})
		})
	}
}

func TestArtistsDeleted_WhenWithSongsAndOrAlbums_ShouldPublishMessages(t *testing.T) {
	tests := []struct {
		name    string
		artists []model.Artist
	}{
		{
			"Delete with songs",
			[]model.Artist{
				{
					ID:        uuid.New(),
					Name:      "Artist 1 - with songs",
					UpdatedAt: time.Now().UTC(),
					UserID:    uuid.New(),
					Songs: []model.Song{
						{ID: uuid.New()},
						{ID: uuid.New()},
					},
				},
				{ID: uuid.New(), Name: "Artist 2"},
			},
		},
		{
			"Delete with albums",
			[]model.Artist{
				{
					ID:        uuid.New(),
					Name:      "Artist 1 - with albums",
					UpdatedAt: time.Now().UTC(),
					UserID:    uuid.New(),
					Albums: []model.Album{
						{ID: uuid.New()},
						{ID: uuid.New()},
					},
				},
				{ID: uuid.New(), Name: "Artist 2"},
			},
		},
		{
			"Delete with songs and albums",
			[]model.Artist{
				{
					ID:        uuid.New(),
					Name:      "Artist 1 - with songs and albums",
					UpdatedAt: time.Now().UTC(),
					UserID:    uuid.New(),
					Songs: []model.Song{
						{ID: uuid.New()},
						{ID: uuid.New()},
					},
					Albums: []model.Album{
						{ID: uuid.New()},
						{ID: uuid.New()},
					},
				},
				{ID: uuid.New(), Name: "Artist 2"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupSearchData(t, artistData.GetSearchDocuments())

			deleteMessages := utils.SubscribeToTopic(topics.DeleteFromSearchEngineTopic)
			deleteStorageMessages := utils.SubscribeToTopic(topics.DeleteDirectoriesStorageTopic)

			// when
			err := utils.PublishToTopic(topics.ArtistsDeletedTopic, test.artists)

			// then
			assert.NoError(t, err)

			// then - prepare data for assertion
			var artistIDs []uuid.UUID
			var albumIDs []uuid.UUID
			var songIDs []uuid.UUID
			for _, art := range test.artists {
				artistIDs = append(artistIDs, art.ID)
				for _, song := range art.Songs {
					songIDs = append(songIDs, song.ID)
				}
				for _, album := range art.Albums {
					albumIDs = append(albumIDs, album.ID)
				}
			}

			// then - assert data
			assertion.AssertMessage(t, deleteMessages, func(ids []string) {
				assert.Len(t, ids, len(artistIDs)+len(albumIDs)+len(songIDs))

				for i := range artistIDs {
					assertion.ArtistSearchID(t, artistIDs[i], ids[i])
				}
				for i := range albumIDs {
					assertion.AlbumSearchID(t, albumIDs[i], ids[i+len(artistIDs)])
				}
				for i := range songIDs {
					assertion.SongSearchID(t, songIDs[i], ids[i+len(artistIDs)+len(albumIDs)])
				}
			})

			assertion.AssertMessage(t, deleteStorageMessages, func(paths []string) {
				assert.Len(t, paths, len(artistIDs)+len(albumIDs)+len(songIDs))
			})
		})
	}
}

func TestArtistsDeleted_WhenWithRelatedSongsOrAlbums_ShouldPublishMessages(t *testing.T) {
	tests := []struct {
		name    string
		artists []model.Artist
	}{
		{
			"With related songs and albums",
			[]model.Artist{{
				ID:        artistData.ArtistSearchID,
				Name:      "Artist 1",
				ImageURL:  &[]internal.FilePath{"something.png"}[0],
				UpdatedAt: time.Now().UTC(),
				UserID:    uuid.New(),
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupSearchData(t, artistData.GetSearchDocuments())

			deleteMessages := utils.SubscribeToTopic(topics.DeleteFromSearchEngineTopic)
			updateMessages := utils.SubscribeToTopic(topics.UpdateFromSearchEngineTopic)
			deleteStorageMessages := utils.SubscribeToTopic(topics.DeleteDirectoriesStorageTopic)

			// when
			err := utils.PublishToTopic(topics.ArtistsDeletedTopic, test.artists)

			// then
			assert.NoError(t, err)

			assertion.AssertMessage(t, deleteMessages, func(ids []string) {
				assert.Len(t, ids, 1)
				assertion.ArtistSearchID(t, test.artists[0].ID, ids[0])
			})

			assertion.AssertMessage(t, updateMessages, func(documents []any) {
				assert.Len(t, documents, len(artistData.SongSearches)+len(artistData.AlbumSearches))
				for i, albumSearch := range artistData.AlbumSearches {
					assert.Equal(t, documents[i].(model.AlbumSearch).ID, albumSearch.(model.AlbumSearch).ID)
					assert.Nil(t, albumSearch.(model.AlbumSearch).Artist)
				}
				for i, songSearch := range artistData.SongSearches {
					assert.Equal(t, documents[i].(model.SongSearch).ID, songSearch.(model.SongSearch).ID)
					assert.Nil(t, songSearch.(model.SongSearch).Artist)
				}
			})

			assertion.AssertMessage(t, deleteStorageMessages, func(paths []string) {
				assert.Len(t, paths, 1)
			})
		})
	}
}
