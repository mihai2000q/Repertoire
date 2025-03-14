package search

import (
	"github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/assert"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	searchData "repertoire/server/test/integration/test/data/search"
	"repertoire/server/test/integration/test/utils"
	"testing"
	"time"
)

func TestUpdateFromSearchEngine_WhenSuccessful_ShouldUpdateDataFromMeilisearch(t *testing.T) {
	// given
	utils.SeedAndCleanupSearchData(t, searchData.GetSearchDocuments())

	newEntities := []any{
		model.SongSearch{
			SearchBase: model.SearchBase{
				ID:     searchData.SongSearches[0].(model.SongSearch).ID,
				Type:   searchData.SongSearches[0].(model.SongSearch).Type,
				UserID: searchData.SongSearches[0].(model.SongSearch).UserID,
			},
			Title:     "New Song Name",
			UpdatedAt: time.Now(),
			ImageUrl:  searchData.SongSearches[0].(model.SongSearch).ImageUrl,
			Artist:    searchData.SongSearches[0].(model.SongSearch).Artist,
			Album:     searchData.SongSearches[0].(model.SongSearch).Album,
		},
		model.ArtistSearch{
			SearchBase: model.SearchBase{
				ID:     searchData.ArtistSearches[0].(model.ArtistSearch).ID,
				Type:   searchData.ArtistSearches[0].(model.ArtistSearch).Type,
				UserID: searchData.ArtistSearches[0].(model.ArtistSearch).UserID,
			},
			Name:      "New Artist Name",
			UpdatedAt: time.Now(),
			ImageUrl:  searchData.ArtistSearches[0].(model.ArtistSearch).ImageUrl,
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
}
