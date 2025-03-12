package album

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	albumData "repertoire/server/test/integration/test/data/album"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestAlbumsUpdated_WhenSuccessful_ShouldPublishMessage(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	messages := utils.SubscribeToTopic(topics.UpdateFromSearchEngineTopic)

	ids := []uuid.UUID{
		albumData.Albums[0].ID,
		albumData.Albums[1].ID,
		albumData.Albums[2].ID,
	}

	// when
	err := utils.PublishToTopic(topics.AlbumsUpdatedTopic, ids)

	// then
	assert.NoError(t, err)

	db := utils.GetDatabase(t)
	var albums []model.Album
	db.Model(&model.Album{}).
		Joins("Artist").
		Preload("Songs").
		Preload("Songs.Artist").
		Preload("Songs.Album").
		Find(&albums, ids)

	assertion.AssertMessage(t, messages, func(documents []any) {
		documentsLength := 0
		for _, album := range albums {
			documentsLength++
			assertion.AlbumSearch(t, documents[documentsLength].(model.AlbumSearch), album)
			for _, song := range album.Songs {
				assertion.SongSearch(t, documents[documentsLength].(model.SongSearch), song)
				documentsLength++
			}
		}
		assert.Len(t, documents, documentsLength)
	})
}
