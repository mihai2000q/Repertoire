package playlist

import (
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPlaylistDeleted_WhenSuccessful_ShouldPublishMessages(t *testing.T) {
	// given
	searchMessages := utils.SubscribeToTopic(topics.DeleteFromSearchEngineTopic)
	storageMessages := utils.SubscribeToTopic(topics.DeleteDirectoriesStorageTopic)

	playlist := model.Playlist{ID: uuid.New()}

	// when
	err := utils.PublishToTopic(topics.PlaylistDeletedTopic, playlist)

	// then
	assert.NoError(t, err)

	assertion.AssertMessage(t, searchMessages, func(ids []string) {
		assert.Len(t, ids, 1)
		assertion.PlaylistSearchID(t, playlist.ID, ids[0])
	})
	assertion.AssertMessage(t, storageMessages, func(paths []string) {
		assert.Len(t, paths, 1)
	})
}
