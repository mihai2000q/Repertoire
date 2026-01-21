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

func TestPlaylistsDeleted_WhenSuccessful_ShouldPublishMessages(t *testing.T) {
	// given
	searchMessages := utils.SubscribeToTopic(topics.DeleteFromSearchEngineTopic)
	storageMessages := utils.SubscribeToTopic(topics.DeleteDirectoriesStorageTopic)

	playlists := []model.Playlist{
		{ID: uuid.New()},
		{ID: uuid.New()},
	}

	// when
	err := utils.PublishToTopic(topics.PlaylistsDeletedTopic, playlists)

	// then
	assert.NoError(t, err)

	assertion.AssertMessage(t, searchMessages, func(ids []string) {
		assert.Len(t, ids, len(playlists))
		for i := range ids {
			assertion.PlaylistSearchID(t, playlists[i].ID, ids[i])
		}
	})
	assertion.AssertMessage(t, storageMessages, func(paths []string) {
		assert.Len(t, paths, len(playlists))
	})
}
