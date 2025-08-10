package album

import (
	"repertoire/server/internal"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	albumData "repertoire/server/test/integration/test/data/album"
	"repertoire/server/test/integration/test/utils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAlbumsDeleted_WhenSuccessful_ShouldPublishMessages(t *testing.T) {
	tests := []struct {
		name   string
		albums []model.Album
	}{
		{
			"Delete with songs",
			[]model.Album{
				{ID: uuid.New(), Title: "Album 1", ImageURL: &[]internal.FilePath{"something.png"}[0]},
				{ID: uuid.New(), Title: "Album 2"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupSearchData(t, albumData.GetSearchDocuments())

			deleteMessages := utils.SubscribeToTopic(topics.DeleteFromSearchEngineTopic)
			deleteStorageMessages := utils.SubscribeToTopic(topics.DeleteDirectoriesStorageTopic)

			// when
			err := utils.PublishToTopic(topics.AlbumsDeletedTopic, test.albums)

			// then
			assert.NoError(t, err)

			assertion.AssertMessage(t, deleteMessages, func(ids []string) {
				assert.Len(t, ids, len(test.albums))

				for i := range test.albums {
					assertion.AlbumSearchID(t, test.albums[i].ID, ids[i])
				}
			})

			assertion.AssertMessage(t, deleteStorageMessages, func(directoryPaths []string) {
				assert.Len(t, directoryPaths, len(test.albums))
			})
		})
	}
}

func TestAlbumsDeleted_WhenWithSongs_ShouldPublishMessages(t *testing.T) {
	tests := []struct {
		name   string
		albums []model.Album
	}{
		{
			"Delete with songs",
			[]model.Album{
				{
					ID:        uuid.New(),
					Title:     "Album 1 - with songs",
					UpdatedAt: time.Now().UTC(),
					UserID:    uuid.New(),
					Songs: []model.Song{
						{ID: uuid.New()},
						{ID: uuid.New()},
					},
				},
				{ID: uuid.New(), Title: "Album 2"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupSearchData(t, albumData.GetSearchDocuments())

			deleteMessages := utils.SubscribeToTopic(topics.DeleteFromSearchEngineTopic)
			deleteStorageMessages := utils.SubscribeToTopic(topics.DeleteDirectoriesStorageTopic)

			// when
			err := utils.PublishToTopic(topics.AlbumsDeletedTopic, test.albums)

			// then
			assert.NoError(t, err)

			// then - prepare data for assertion
			var albumIDs []uuid.UUID
			var songIDs []uuid.UUID
			for _, alb := range test.albums {
				albumIDs = append(albumIDs, alb.ID)
				for _, song := range alb.Songs {
					songIDs = append(songIDs, song.ID)
				}
			}

			// then - assert data
			assertion.AssertMessage(t, deleteMessages, func(ids []string) {
				assert.Len(t, ids, len(albumIDs)+len(songIDs))

				for i := range albumIDs {
					assertion.AlbumSearchID(t, albumIDs[i], ids[i])
				}
				for i := range songIDs {
					assertion.SongSearchID(t, songIDs[i], ids[i+len(albumIDs)])
				}
			})

			assertion.AssertMessage(t, deleteStorageMessages, func(directoryPaths []string) {
				assert.Len(t, directoryPaths, len(albumIDs)+len(songIDs))
			})
		})
	}
}

func TestAlbumsDeleted_WhenWithRelatedSongs_ShouldPublishMessages(t *testing.T) {
	tests := []struct {
		name   string
		albums []model.Album
	}{
		{
			"With related songs",
			[]model.Album{{
				ID:        albumData.AlbumSearchID,
				Title:     "Album 1",
				ImageURL:  &[]internal.FilePath{"something.png"}[0],
				UpdatedAt: time.Now().UTC(),
				UserID:    uuid.New(),
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupSearchData(t, albumData.GetSearchDocuments())

			deleteMessages := utils.SubscribeToTopic(topics.DeleteFromSearchEngineTopic)
			updateMessages := utils.SubscribeToTopic(topics.UpdateFromSearchEngineTopic)
			deleteStorageMessages := utils.SubscribeToTopic(topics.DeleteDirectoriesStorageTopic)

			// when
			err := utils.PublishToTopic(topics.AlbumsDeletedTopic, test.albums)

			// then
			assert.NoError(t, err)

			assertion.AssertMessage(t, deleteMessages, func(ids []string) {
				assert.Len(t, ids, len(test.albums[0].Songs)+1)
				assertion.AlbumSearchID(t, test.albums[0].ID, ids[0])
			})

			assertion.AssertMessage(t, updateMessages, func(documents []any) {
				assert.Len(t, documents, len(albumData.SongSearches))
				for i, songSearch := range albumData.SongSearches {
					assert.Equal(t, documents[i].(model.SongSearch).ID, songSearch.(model.SongSearch).ID)
					assert.Nil(t, songSearch.(model.SongSearch).Album)
				}
			})

			assertion.AssertMessage(t, deleteStorageMessages, func(directoryPaths []string) {
				assert.Len(t, directoryPaths, 1)
			})
		})
	}
}
