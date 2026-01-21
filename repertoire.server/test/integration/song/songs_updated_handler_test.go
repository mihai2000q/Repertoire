package song

import (
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

	assertion.AssertMessage(t, messages, func(songSearches []model.SongSearch) {
		assert.Len(t, songSearches, len(songs))
		for i := range songs {
			assertion.SongSearch(t, songSearches[i], songs[i])
		}
	})
}
