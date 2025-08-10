package playlist

import (
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/utils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPlaylistCreated_WhenSuccessful_ShouldPublishMessage(t *testing.T) {
	// given
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

	assertion.AssertMessage(t, messages, func(documents []any) {
		assert.Len(t, documents, 1)
		playlistSearch := utils.UnmarshallDocument[model.PlaylistSearch](documents[0])
		assertion.PlaylistSearch(t, playlistSearch, playlist)
	})
}
