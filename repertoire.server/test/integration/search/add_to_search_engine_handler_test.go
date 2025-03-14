package search

import (
	"github.com/google/uuid"
	"github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/assert"
	"repertoire/server/internal/message/topics"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestAddToSearchEngine_WhenSuccessful_ShouldAddDataToMeilisearch(t *testing.T) {
	// given
	data := []map[string]any{
		{
			"id":     uuid.New().String(),
			"name":   "some name",
			"userId": "some-user-id",
		},
		{
			"id":     uuid.New().String(),
			"name":   "some-name-2",
			"userId": "some-user-id",
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
}
