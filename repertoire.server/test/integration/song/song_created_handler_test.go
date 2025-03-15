package song

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"repertoire/server/internal"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"
	"time"
)

func TestSongCreated_WhenSuccessful_ShouldPublishMessage(t *testing.T) {
	tests := []struct {
		name string
		song model.Song
	}{
		{
			"Normal",
			model.Song{
				ID:        uuid.New(),
				Title:     "Song 1",
				ImageURL:  &[]internal.FilePath{"something.png"}[0],
				UpdatedAt: time.Now(),
				UserID:    uuid.New(),
			},
		},
		{
			"With New Artist",
			model.Song{
				ID:        uuid.New(),
				Title:     "Song 2",
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
			"With New Album",
			model.Song{
				ID:        uuid.New(),
				Title:     "Song 3",
				UpdatedAt: time.Now(),
				UserID:    uuid.New(),
				Album: &model.Album{
					ID:        uuid.New(),
					Title:     "Album 1",
					UpdatedAt: time.Now(),
				},
			},
		},
		{
			"With New Artist and New Album",
			model.Song{
				ID:        uuid.New(),
				Title:     "Song 4",
				UpdatedAt: time.Now(),
				UserID:    uuid.New(),
				Artist: &model.Artist{
					ID:        uuid.New(),
					Name:      "Artist 2",
					UpdatedAt: time.Now(),
				},
				Album: &model.Album{
					ID:        uuid.New(),
					Title:     "Album 2",
					UpdatedAt: time.Now(),
				},
			},
		},
		{
			"With Existing Artist",
			model.Song{
				ID:        uuid.New(),
				Title:     "Song 5",
				UpdatedAt: time.Now(),
				UserID:    uuid.New(),
				ArtistID:  &[]uuid.UUID{songData.Artists[0].ID}[0],
			},
		},
		{
			"With Existing Artist and Album",
			model.Song{
				ID:        uuid.New(),
				Title:     "Song 6",
				UpdatedAt: time.Now(),
				UserID:    uuid.New(),
				ArtistID:  &[]uuid.UUID{songData.Artists[0].ID}[0],
				AlbumID:   &[]uuid.UUID{songData.Albums[0].ID}[0],
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

			messages := utils.SubscribeToTopic(topics.AddToSearchEngineTopic)

			// when
			err := utils.PublishToTopic(topics.SongCreatedTopic, test.song)

			// then
			assert.NoError(t, err)

			assertion.AssertMessage(t, messages, func(documents []any) {
				var artistIndex *int
				var albumIndex *int
				if test.song.Artist != nil && test.song.Album != nil {
					assert.Len(t, documents, 3)
					artistIndex = &[]int{1}[0]
					albumIndex = &[]int{2}[0]
				} else if test.song.Artist != nil {
					assert.Len(t, documents, 2)
					artistIndex = &[]int{1}[0]
				} else if test.song.Album != nil {
					assert.Len(t, documents, 2)
					albumIndex = &[]int{1}[0]
				}

				if test.song.ArtistID != nil {
					test.song.Artist = &songData.Artists[0]
				}
				if test.song.AlbumID != nil {
					test.song.Album = &songData.Albums[0]
				}
				songSearch := utils.UnmarshallDocument[model.SongSearch](documents[0])
				assertion.SongSearch(t, songSearch, test.song)
				if test.song.ArtistID != nil {
					assert.Equal(t, songSearch.Artist.ID, *test.song.ArtistID)
				}
				if test.song.AlbumID != nil {
					assert.Equal(t, songSearch.Album.ID, *test.song.AlbumID)
				}

				if artistIndex != nil {
					artistSearch := utils.UnmarshallDocument[model.ArtistSearch](documents[*artistIndex])
					assertion.ArtistSearch(t, artistSearch, *test.song.Artist)
				}
				if albumIndex != nil {
					albumSearch := utils.UnmarshallDocument[model.AlbumSearch](documents[*albumIndex])
					assertion.AlbumSearch(t, albumSearch, *test.song.Album)
				}
			})
		})
	}
}
