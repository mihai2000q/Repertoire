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

	song := model.Song{ID: uuid.New()}

	// when
	err := utils.PublishToTopic(topics.SongDeletedTopic, song)

	// then
	assert.NoError(t, err)

	assertion.AssertMessage(t, searchMessages, func(ids []string) {
		assert.Len(t, ids, 1)
		assertion.SongSearchID(t, song.ID, ids[0])
	})
	assertion.AssertMessage(t, storageMessages, func(paths []string) {
		assert.Len(t, paths, 1)
	})
}
