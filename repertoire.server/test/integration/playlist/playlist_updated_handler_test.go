package playlist

import (
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	playlistData "repertoire/server/test/integration/test/data/playlist"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlaylistUpdated_WhenSuccessful_ShouldPublishMessage(t *testing.T) {
	// given
	messages := utils.SubscribeToTopic(topics.UpdateFromSearchEngineTopic)

	playlist := playlistData.Playlists[0]

	// when
	err := utils.PublishToTopic(topics.PlaylistUpdatedTopic, playlist)

	// then
	assert.NoError(t, err)

	assertion.AssertMessage(t, messages, func(documents []any) {
		assert.Len(t, documents, 1)
		playlistSearch := utils.UnmarshallDocument[model.PlaylistSearch](documents[0])
		assertion.PlaylistSearch(t, playlistSearch, playlist)
	})
}
