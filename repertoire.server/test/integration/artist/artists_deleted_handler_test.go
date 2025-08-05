package artist

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"repertoire/server/internal"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	artistData "repertoire/server/test/integration/test/data/artist"
	"repertoire/server/test/integration/test/utils"
	"testing"
	"time"
)

func TestArtistsDeleted_WhenSuccessful_ShouldPublishMessages(t *testing.T) {
	tests := []struct {
		name    string
		artists []model.Artist
	}{
		{
			"Normal",
			[]model.Artist{{
				ID:        uuid.New(),
				Name:      "Artist 1",
				ImageURL:  &[]internal.FilePath{"something.png"}[0],
				UpdatedAt: time.Now().UTC(),
				UserID:    uuid.New(),
			}},
		},
		{
			"With related songs and albums",
			[]model.Artist{{
				ID:        artistData.ArtistSearchID,
				Name:      "Artist 2",
				ImageURL:  &[]internal.FilePath{"something.png"}[0],
				UpdatedAt: time.Now().UTC(),
				UserID:    uuid.New(),
			}},
		},
		{
			"Delete with songs",
			[]model.Artist{{
				ID:        uuid.New(),
				Name:      "Artist 3",
				UpdatedAt: time.Now().UTC(),
				UserID:    uuid.New(),
				Songs: []model.Song{
					{ID: uuid.New()},
					{ID: uuid.New()},
				},
			}},
		},
		{
			"Delete with albums",
			[]model.Artist{{
				ID:        uuid.New(),
				Name:      "Artist 4",
				UpdatedAt: time.Now().UTC(),
				UserID:    uuid.New(),
				Albums: []model.Album{
					{ID: uuid.New()},
					{ID: uuid.New()},
				},
			}},
		},
		{
			"Delete with songs and albums",
			[]model.Artist{{
				ID:        uuid.New(),
				Name:      "Artist 5",
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
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupSearchData(t, artistData.GetSearchDocuments())

			deleteMessages := utils.SubscribeToTopic(topics.DeleteFromSearchEngineTopic)
			var updateMessages utils.SubscribedToTopic
			if len(test.artists[0].Songs) == 0 || len(test.artists[0].Albums) == 0 {
				updateMessages = utils.SubscribeToTopic(topics.UpdateFromSearchEngineTopic)
			}
			deleteStorageMessages := utils.SubscribeToTopic(topics.DeleteDirectoriesStorageTopic)

			// when
			err := utils.PublishToTopic(topics.ArtistsDeletedTopic, test.artists)

			// then
			assert.NoError(t, err)

			assertion.AssertMessage(t, deleteMessages, func(ids []string) {
				assert.Len(t, ids, len(test.artists[0].Songs)+len(test.artists[0].Albums)+1)
				assertion.ArtistSearchID(t, test.artists[0].ID, ids[0])
				for i, song := range test.artists[0].Songs {
					assertion.SongSearchID(t, song.ID, ids[i+1])
				}
				for i, album := range test.artists[0].Albums {
					assertion.AlbumSearchID(t, album.ID, ids[i+1+len(test.artists[0].Songs)])
				}
			})

			if len(test.artists[0].Songs) == 0 || len(test.artists[0].Albums) == 0 {
				assertion.AssertMessage(t, updateMessages, func(documents []any) {
					assert.Len(t, documents, len(artistData.SongSearches)+len(artistData.AlbumSearches))
					for i, songSearch := range artistData.SongSearches {
						assert.Equal(t, documents[i].(model.SongSearch).ID, songSearch.(model.SongSearch).ID)
						assert.Nil(t, songSearch.(model.SongSearch).Artist)
					}
					for i, albumSearch := range artistData.AlbumSearches {
						assert.Equal(t, documents[i].(model.AlbumSearch).ID, albumSearch.(model.AlbumSearch).ID)
						assert.Nil(t, albumSearch.(model.AlbumSearch).Artist)
					}
				})
			}

			assertion.AssertMessage(t, deleteStorageMessages, func(paths []string) {
				assert.Len(t, paths, len(test.artists[0].Songs)+len(test.artists[0].Albums)+1)
			})
		})
	}
}
