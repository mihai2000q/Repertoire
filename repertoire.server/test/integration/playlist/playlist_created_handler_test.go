package playlist

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	playlistData "repertoire/server/test/integration/test/data/playlist"
	"repertoire/server/test/integration/test/utils"
	"testing"
	"time"
)

func TestPlaylistCreated_WhenSuccessful_ShouldPublishMessage(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	messages := utils.SubscribeToTopic(topics.AddToSearchEngineTopic)

	playlist := model.Playlist{
		ID:        uuid.New(),
		Title:     "Playlist",
		UpdatedAt: time.Now(),
	}

	// when
	err := utils.PublishToTopic(topics.PlaylistCreatedTopic, playlist)

	// then
	assert.NoError(t, err)

	assertion.AssertMessage(t, messages, topics.AddToSearchEngineTopic, func(documents []any) {
		assert.Len(t, documents, 1)
		playlistSearch := utils.UnmarshallDocument[model.PlaylistSearch](documents[0])
		assertion.PlaylistSearch(t, playlistSearch, playlist)
	})
}
