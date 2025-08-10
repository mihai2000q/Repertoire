package song

import (
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSongDeleted_WhenSuccessful_ShouldPublishMessages(t *testing.T) {
	// given
	searchMessages := utils.SubscribeToTopic(topics.DeleteFromSearchEngineTopic)
	storageMessages := utils.SubscribeToTopic(topics.DeleteDirectoriesStorageTopic)

	songs := []model.Song{
		{ID: uuid.New()},
		{ID: uuid.New()},
	}

	// when
	err := utils.PublishToTopic(topics.SongsDeletedTopic, songs)

	// then
	assert.NoError(t, err)

	assertion.AssertMessage(t, searchMessages, func(ids []string) {
		assert.Len(t, ids, len(songs))
		for i := range songs {
			assertion.SongSearchID(t, songs[i].ID, ids[i])
		}
	})
	assertion.AssertMessage(t, storageMessages, func(paths []string) {
		assert.Len(t, paths, len(songs))
	})
}
