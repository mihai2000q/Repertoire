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
	type TestStruct struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		UserID string `json:"userId"`
	}

	userID := uuid.New().String()
	data := []TestStruct{
		{
			ID:     uuid.New().String(),
			Name:   "some name",
			UserID: userID,
		},
		{
			ID:     uuid.New().String(),
			Name:   "some-name-2",
			UserID: userID,
		},
	}

	searchClient := utils.GetSearchClient(t)
	tasks, _ := searchClient.GetTasks(nil)

	// when
	err := utils.PublishToTopic(topics.AddToSearchEngineTopic, data)

	// then
	assert.NoError(t, err)

	utils.WaitForSearchTasksToStart(searchClient, tasks.Total)
	utils.WaitForAllSearchTasks(searchClient)
	var result meilisearch.DocumentsResult
	_ = searchClient.Index("search").GetDocuments(nil, &result)
	assert.Len(t, result.Results, len(data))
	for i := range result.Results {
		expected := data[i]
		var actual TestStruct
		_ = result.Results[i].DecodeInto(&actual)
		assert.Equal(t, expected.ID, actual.ID)
		assert.Equal(t, expected.Name, actual.Name)
		assert.Equal(t, expected.UserID, actual.UserID)
	}

	tasks, _ = searchClient.GetTasks(nil)
	latestTaskID := strconv.FormatInt((*tasks).Results[0].UID, 10)
	cachedUserID, _ := core.MeiliCache.Get("task-" + latestTaskID)
	assert.Equal(t, userID, cachedUserID)
}
