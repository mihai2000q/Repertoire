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

	assertion.AssertMessage(t, messages, func(documents []any) {
		assert.Len(t, documents, 1)
		artistSearch := utils.UnmarshallDocument[model.ArtistSearch](documents[0])
		assertion.ArtistSearch(t, artistSearch, artist)
	})
}
