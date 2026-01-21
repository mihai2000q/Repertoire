package artist

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

func TestArtistCreated_WhenSuccessful_ShouldPublishMessage(t *testing.T) {
	// given
	messages := utils.SubscribeToTopic(topics.AddToSearchEngineTopic)

	artist := model.Artist{
		ID:        uuid.New(),
		Name:      "Artist",
		UpdatedAt: time.Now(),
	}

	// when
	err := utils.PublishToTopic(topics.ArtistCreatedTopic, artist)

	// then
	assert.NoError(t, err)

	assertion.AssertMessage(t, messages, func(artistSearches []model.ArtistSearch) {
		assert.Len(t, artistSearches, 1)
		assertion.ArtistSearch(t, artistSearches[0], artist)
	})
}
