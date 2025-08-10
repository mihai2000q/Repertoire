package user

import (
	"repertoire/server/internal/message/topics"
	"repertoire/server/test/integration/test/assertion"
	userData "repertoire/server/test/integration/test/data/user"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserDeleted_WhenSuccessful_ShouldPublishMessages(t *testing.T) {
	// given
	utils.SeedAndCleanupSearchData(t, userData.GetSearchDocuments())

	userID := userData.UserSearchID
	searches := userData.Searches

	searchMessage := utils.SubscribeToTopic(topics.DeleteFromSearchEngineTopic)
	storageMessage := utils.SubscribeToTopic(topics.DeleteDirectoriesStorageTopic)

	// when
	err := utils.PublishToTopic(topics.UserDeletedTopic, userID)

	// then
	assert.NoError(t, err)

	assertion.AssertMessage(t, searchMessage, func(ids []string) {
		assert.Len(t, ids, len(searches))
		for i := range searches {
			assert.Equal(t, ids[i], searches[i]["id"].(string))
		}
	})
	assertion.AssertMessage(t, storageMessage, func(paths []string) {
		assert.Len(t, paths, 1)
	})
}
