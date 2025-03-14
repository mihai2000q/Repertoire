package user

import (
	"github.com/stretchr/testify/assert"
	"repertoire/server/internal/message/topics"
	"repertoire/server/test/integration/test/assertion"
	userData "repertoire/server/test/integration/test/data/user"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestUserDeleted_WhenSuccessful_ShouldPublishMessage(t *testing.T) {
	// given
	utils.SeedAndCleanupSearchData(t, userData.GetSearchDocuments())

	userID := userData.UserSearchID
	searches := userData.Searches

	deleteSearchMessage := utils.SubscribeToTopic(topics.DeleteFromSearchEngineTopic)
	deleteStorageMessage := utils.SubscribeToTopic(topics.DeleteDirectoriesStorageTopic)

	// when
	err := utils.PublishToTopic(topics.UserDeletedTopic, userID)

	// then
	assert.NoError(t, err)

	assertion.AssertMessage(t, deleteSearchMessage, func(ids []string) {
		assert.Len(t, ids, len(searches))
		for i := range searches {
			assert.Equal(t, ids[i], searches[i]["id"].(string))
		}
	})
	assertion.AssertMessage(t, deleteStorageMessage, func(paths []string) {
		assert.Len(t, paths, 1)
	})
}
