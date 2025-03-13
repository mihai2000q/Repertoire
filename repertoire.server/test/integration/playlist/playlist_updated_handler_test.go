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
)

func TestPlaylistUpdated_WhenSuccessful_ShouldPublishMessage(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	messages := utils.SubscribeToTopic(topics.UpdateFromSearchEngineTopic)

	ids := []uuid.UUID{
		playlistData.Playlists[0].ID,
		playlistData.Playlists[1].ID,
		playlistData.Playlists[2].ID,
	}

	// when
	err := utils.PublishToTopic(topics.PlaylistUpdatedTopic, ids)

	// then
	assert.NoError(t, err)

	db := utils.GetDatabase(t)
	var playlists []model.Playlist
	db.Model(&model.Playlist{}).Find(&playlists, ids)

	assertion.AssertMessage(t, messages, topics.UpdateFromSearchEngineTopic, func(documents []any) {
		assert.Len(t, documents, len(playlists))
		for i := range playlists {
			playlistSearch := utils.UnmarshallDocument[model.PlaylistSearch](documents[i])
			assertion.PlaylistSearch(t, playlistSearch, playlists[i])
		}
	})
}
