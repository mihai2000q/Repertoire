package album

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"repertoire/server/internal"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	albumData "repertoire/server/test/integration/test/data/album"
	"repertoire/server/test/integration/test/utils"
	"testing"
	"time"
)

func TestAlbumDeleted_WhenSuccessful_ShouldPublishMessages(t *testing.T) {
	tests := []struct {
		name  string
		album model.Album
	}{
		{
			"Normal",
			model.Album{
				ID:        uuid.New(),
				Title:     "Album 1",
				ImageURL:  &[]internal.FilePath{"something.png"}[0],
				UpdatedAt: time.Now().UTC(),
				UserID:    uuid.New(),
			},
		},
		{
			"With related songs",
			model.Album{
				ID:        albumData.AlbumSearchID,
				Title:     "Album 2",
				ImageURL:  &[]internal.FilePath{"something.png"}[0],
				UpdatedAt: time.Now().UTC(),
				UserID:    uuid.New(),
			},
		},
		{
			"Delete with songs",
			model.Album{
				ID:        uuid.New(),
				Title:     "Album 3",
				UpdatedAt: time.Now().UTC(),
				UserID:    uuid.New(),
				Songs: []model.Song{
					{ID: uuid.New()},
					{ID: uuid.New()},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupSearchData(t, albumData.GetSearchDocuments())

			deleteMessages := utils.SubscribeToTopic(topics.DeleteFromSearchEngineTopic)
			var updateMessages utils.SubscribedToTopic
			if len(test.album.Songs) == 0 {
				updateMessages = utils.SubscribeToTopic(topics.UpdateFromSearchEngineTopic)
			}
			deleteStorageMessages := utils.SubscribeToTopic(topics.DeleteDirectoriesStorageTopic)

			// when
			err := utils.PublishToTopic(topics.AlbumDeletedTopic, test.album)

			// then
			assert.NoError(t, err)

			assertion.AssertMessage(t, deleteMessages, func(ids []string) {
				assert.Len(t, ids, len(test.album.Songs)+1)
				assertion.AlbumSearchID(t, test.album.ID, ids[0])
				for i, song := range test.album.Songs {
					assertion.SongSearchID(t, song.ID, ids[i+1])
				}
			})

			if len(test.album.Songs) == 0 {
				assertion.AssertMessage(t, updateMessages, func(documents []any) {
					assert.Len(t, documents, len(albumData.SongSearches))
					for i, songSearch := range albumData.SongSearches {
						assert.Equal(t, documents[i].(model.SongSearch).ID, songSearch.(model.SongSearch).ID)
						assert.Nil(t, songSearch.(model.SongSearch).Album)
					}
				})
			}

			assertion.AssertMessage(t, deleteStorageMessages, func(directoryPaths []string) {
				assert.Len(t, directoryPaths, len(test.album.Songs)+1)
			})
		})
	}
}
