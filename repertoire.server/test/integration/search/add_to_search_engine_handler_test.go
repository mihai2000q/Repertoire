package search

import (
	"repertoire/server/internal/message/topics"
	"repertoire/server/test/integration/test/core"
	"repertoire/server/test/integration/test/utils"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/assert"
)

func TestAddToSearchEngine_WhenSuccessful_ShouldAddDataToMeilisearch(t *testing.T) {
	// given
	userID := "some-user-id"
	data := []map[string]any{
		{
			"id":     uuid.New().String(),
			"name":   "some name",
			"userId": userID,
		},
		{
			"id":     uuid.New().String(),
			"name":   "some-name-2",
			"userId": userID,
		},
	}

	searchClient := utils.GetSearchClient(t)
	tasks, _ := searchClient.GetTasks(&meilisearch.TasksQuery{})

	// when
	err := utils.PublishToTopic(topics.AddToSearchEngineTopic, data)

	// then
	assert.NoError(t, err)

	utils.WaitForSearchTasksToStart(searchClient, tasks.Total)
	utils.WaitForAllSearchTasks(searchClient)
	var result meilisearch.DocumentsResult
	_ = searchClient.Index("search").GetDocuments(&meilisearch.DocumentsQuery{}, &result)
	assert.Len(t, result.Results, len(data))
	for i := range result.Results {
		expected := data[i]
		actual := result.Results[i]
		assert.Equal(t, expected["id"], actual["id"])
		assert.Equal(t, expected["name"], actual["name"])
		assert.Equal(t, expected["userId"], actual["userId"])
	}

	tasks, _ = searchClient.GetTasks(&meilisearch.TasksQuery{})
	latestTaskID := strconv.FormatInt((*tasks).Results[0].UID, 10)
	cachedUserID, _ := core.MeiliCache.Get("task-" + latestTaskID)
	assert.Equal(t, userID, cachedUserID)
}
