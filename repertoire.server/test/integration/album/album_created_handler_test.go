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
				UpdatedAt: time.Now().UTC(),
				UserID:    uuid.New(),
			},
		},
		{
			"With New Artist",
			model.Album{
				ID:        uuid.New(),
				Title:     "Album 1",
				UpdatedAt: time.Now().UTC(),
				UserID:    uuid.New(),
				Artist: &model.Artist{
					ID:        uuid.New(),
					Name:      "Artist 1",
					UpdatedAt: time.Now().UTC(),
				},
			},
		},
		{
			"With Existing Artist",
			model.Album{
				ID:        uuid.New(),
				Title:     "Album 1",
				UpdatedAt: time.Now().UTC(),
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

			assertion.AssertMessage(t, messages, func(documents []any) {
				if test.album.Artist != nil {
					assert.Len(t, documents, 2)
					assertion.ArtistSearch(t, documents[1].(model.ArtistSearch), *test.album.Artist)
				} else {
					assert.Len(t, documents, 1)
				}
				assertion.AlbumSearch(t, documents[0].(model.AlbumSearch), test.album)
			})
		})
	}
}
