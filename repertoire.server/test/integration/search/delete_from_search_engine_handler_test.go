package search

import (
	"github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/assert"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	searchData "repertoire/server/test/integration/test/data/search"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestDeleteFromSearchEngine_WhenSuccessful_ShouldDeleteDataFromMeilisearch(t *testing.T) {
	// given
	utils.SeedAndCleanupSearchData(t, searchData.GetSearchDocuments())

	ids := []string{
		searchData.ArtistSearches[0].(model.ArtistSearch).ID,
		searchData.SongSearches[0].(model.SongSearch).ID,
	}

	// when
	err := utils.PublishToTopic(topics.DeleteFromSearchEngineTopic, ids)

	// then
	assert.NoError(t, err)

	searchClient := utils.GetSearchClient(t)
	utils.WaitForAllSearchTasks(searchClient)
	for _, id := range ids {
		var entity *map[string]any
		getErr := searchClient.Index("search").GetDocument(id, &meilisearch.DocumentQuery{}, &entity)
		assert.Nil(t, entity)
		assert.Error(t, getErr)
	}
}
