package playlist

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestPlaylistDeleted_WhenSuccessful_ShouldPublishMessage(t *testing.T) {
	// given
	messages := utils.SubscribeToTopic(topics.DeleteFromSearchEngineTopic)

	playlist := model.Playlist{ID: uuid.New()}

	// when
	err := utils.PublishToTopic(topics.PlaylistDeletedTopic, playlist)

	// then
	assert.NoError(t, err)

	assertion.AssertMessage(t, messages, topics.DeleteFromSearchEngineTopic, func(ids []string) {
		assert.Len(t, ids, 1)
		assertion.PlaylistSearchID(t, playlist.ID, ids[0])
	})
}
