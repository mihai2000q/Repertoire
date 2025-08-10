package search

import (
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	searchData "repertoire/server/test/integration/test/data/search"
	"repertoire/server/test/integration/test/utils"
	"strconv"
	"testing"
	"time"

	"github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/assert"
)

func TestUpdateFromSearchEngine_WhenSuccessful_ShouldUpdateDataFromMeilisearch(t *testing.T) {
	// given
	utils.SeedAndCleanupSearchData(t, searchData.GetSearchDocuments())

	userID := searchData.SongSearches[0].(model.SongSearch).UserID
	newEntities := []any{
		model.SongSearch{
			SearchBase: model.SearchBase{
				ID:        searchData.SongSearches[0].(model.SongSearch).ID,
				UpdatedAt: time.Now().UTC(),
				CreatedAt: searchData.SongSearches[0].(model.SongSearch).CreatedAt,
				Type:      searchData.SongSearches[0].(model.SongSearch).Type,
				UserID:    userID,
			},
			Title:    "New Song Name",
			ImageUrl: searchData.SongSearches[0].(model.SongSearch).ImageUrl,
			Artist:   searchData.SongSearches[0].(model.SongSearch).Artist,
			Album:    searchData.SongSearches[0].(model.SongSearch).Album,
		},
		model.ArtistSearch{
			SearchBase: model.SearchBase{
				ID:        searchData.ArtistSearches[0].(model.ArtistSearch).ID,
				UpdatedAt: time.Now().UTC(),
				CreatedAt: searchData.ArtistSearches[0].(model.ArtistSearch).CreatedAt,
				Type:      searchData.ArtistSearches[0].(model.ArtistSearch).Type,
				UserID:    userID,
			},
			Name:     "New Artist Name",
			ImageUrl: searchData.ArtistSearches[0].(model.ArtistSearch).ImageUrl,
		},
	}

	searchClient := utils.GetSearchClient(t)
	tasks, _ := searchClient.GetTasks(&meilisearch.TasksQuery{})

	// when
	err := utils.PublishToTopic(topics.UpdateFromSearchEngineTopic, newEntities)

	// then
	assert.NoError(t, err)

	utils.WaitForSearchTasksToStart(searchClient, tasks.Total)
	utils.WaitForAllSearchTasks(searchClient)
	for _, expectedEntity := range newEntities {
		unmarshalledExpectedEntity := utils.UnmarshallDocument[map[string]any](expectedEntity)
		var actualEntity *map[string]any
		getErr := searchClient.Index("search").GetDocument(
			unmarshalledExpectedEntity["id"].(string),
			&meilisearch.DocumentQuery{},
			&actualEntity,
		)
		assert.NoError(t, getErr)
		assert.NotNil(t, actualEntity)
		assert.Equal(t, unmarshalledExpectedEntity, *actualEntity)
	}

	tasks, _ = searchClient.GetTasks(&meilisearch.TasksQuery{})
	latestTaskID := strconv.FormatInt((*tasks).Results[0].UID, 10)
	cachedUserID, _ := core.MeiliCache.Get("task-" + latestTaskID)
	assert.Equal(t, userID.String(), cachedUserID)
}
