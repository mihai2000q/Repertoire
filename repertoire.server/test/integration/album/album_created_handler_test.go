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

func TestAlbumCreated_WhenSuccessful_ShouldPublishMessage(t *testing.T) {
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
				UpdatedAt: time.Now(),
				UserID:    uuid.New(),
			},
		},
		{
			"With New Artist",
			model.Album{
				ID:        uuid.New(),
				Title:     "Album 2",
				UpdatedAt: time.Now(),
				UserID:    uuid.New(),
				Artist: &model.Artist{
					ID:        uuid.New(),
					Name:      "Artist 1",
					UpdatedAt: time.Now(),
				},
			},
		},
		{
			"With Existing Artist",
			model.Album{
				ID:        uuid.New(),
				Title:     "Album 3",
				UpdatedAt: time.Now(),
				UserID:    uuid.New(),
				ArtistID:  &[]uuid.UUID{albumData.Artists[0].ID}[0],
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

			messages := utils.SubscribeToTopic(topics.AddToSearchEngineTopic)

			// when
			err := utils.PublishToTopic(topics.AlbumCreatedTopic, test.album)

			// then
			assert.NoError(t, err)

			assertion.AssertMessage(t, messages, func(documents []map[string]any) {
				if test.album.Artist != nil {
					artistSearch := utils.UnmarshalDocument[model.ArtistSearch](documents[1])
					assertion.ArtistSearch(t, artistSearch, *test.album.Artist)
					assert.Len(t, documents, 2)
				} else {
					assert.Len(t, documents, 1)
				}
				if test.album.ArtistID != nil {
					test.album.Artist = &albumData.Artists[0]
				}
				albumSearch := utils.UnmarshalDocument[model.AlbumSearch](documents[0])
				assertion.AlbumSearch(t, albumSearch, test.album)
			})
		})
	}
}
