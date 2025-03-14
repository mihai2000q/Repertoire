package song

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestSongsUpdated_WhenSuccessful_ShouldPublishMessage(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	messages := utils.SubscribeToTopic(topics.UpdateFromSearchEngineTopic)

	ids := []uuid.UUID{
		songData.Songs[0].ID,
		songData.Songs[1].ID,
		songData.Songs[4].ID,
	}

	// when
	err := utils.PublishToTopic(topics.SongsUpdatedTopic, ids)

	// then
	assert.NoError(t, err)

	db := utils.GetDatabase(t)
	var songs []model.Song
	db.Model(&model.Song{}).
		Joins("Artist").
		Joins("Album").
		Find(&songs, ids)

	assertion.AssertMessage(t, messages, topics.UpdateFromSearchEngineTopic, func(documents []any) {
		assert.Len(t, documents, len(songs))
		for i := range songs {
			songSearch := utils.UnmarshallDocument[model.SongSearch](documents[i])
			assertion.SongSearch(t, songSearch, songs[i])
		}
	})
}
