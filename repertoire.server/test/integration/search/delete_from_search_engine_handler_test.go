package search

import (
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	searchData "repertoire/server/test/integration/test/data/search"
	"repertoire/server/test/integration/test/utils"
	"strconv"
	"testing"

	"github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/assert"
)

func TestDeleteFromSearchEngine_WhenSuccessful_ShouldDeleteDataFromMeilisearch(t *testing.T) {
	// given
	utils.SeedAndCleanupSearchData(t, searchData.GetSearchDocuments())

	userID := searchData.SongSearches[0].(model.SongSearch).UserID

	ids := []string{
		searchData.ArtistSearches[0].(model.ArtistSearch).ID,
		searchData.SongSearches[0].(model.SongSearch).ID,
	}

	searchClient := utils.GetSearchClient(t)
	tasks, _ := searchClient.GetTasks(&meilisearch.TasksQuery{})

	// when
	err := utils.PublishToTopic(topics.DeleteFromSearchEngineTopic, ids)

	// then
	assert.NoError(t, err)

	utils.WaitForSearchTasksToStart(searchClient, tasks.Total)
	utils.WaitForAllSearchTasks(searchClient)
	for _, id := range ids {
		var entity *map[string]any
		getErr := searchClient.Index("search").GetDocument(id, &meilisearch.DocumentQuery{}, &entity)
		assert.Nil(t, entity)
		assert.Error(t, getErr)
	}

	tasks, _ = searchClient.GetTasks(&meilisearch.TasksQuery{})
	latestTaskID := strconv.FormatInt((*tasks).Results[0].UID, 10)
	cachedUserID, _ := core.MeiliCache.Get("task-" + latestTaskID)
	assert.Equal(t, userID.String(), cachedUserID)
}
