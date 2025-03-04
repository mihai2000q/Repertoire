package service

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/meilisearch/meilisearch-go"
	"log"
	"repertoire/server/internal/enums"
	"repertoire/server/internal/wrapper"
)

type SearchEngineService interface {
	Get(
		query string,
		currentPage *int,
		pageSize *int,
		searchType *enums.SearchType,
		userID uuid.UUID,
	) (wrapper.WithTotalCount[any], *wrapper.ErrorCode)
	Add(items []any)
}

type searchEngineService struct {
	client meilisearch.ServiceManager
}

func NewSearchEngineService(client meilisearch.ServiceManager) SearchEngineService {
	return searchEngineService{client: client}
}

func (s searchEngineService) Get(
	query string,
	currentPage *int,
	pageSize *int,
	searchType *enums.SearchType,
	userID uuid.UUID,
) (wrapper.WithTotalCount[any], *wrapper.ErrorCode) {
	request := &meilisearch.SearchRequest{}

	if currentPage != nil && pageSize != nil {
		request.Offset = int64((*currentPage - 1) * *pageSize)
		request.Limit = int64(*pageSize)
	}

	if searchType != nil {
		request.Filter = "type = " + string(*searchType) + " AND userId = " + userID.String()
	} else {
		request.Filter = "userId = " + userID.String()
	}

	searchResult, err := s.client.Index("search").Search(query, request)
	if err != nil {
		return wrapper.WithTotalCount[any]{}, wrapper.InternalServerError(err)
	}

	result := wrapper.WithTotalCount[any]{
		Models:     searchResult.Hits,
		TotalCount: searchResult.EstimatedTotalHits,
	}

	return result, nil
}

func (s searchEngineService) Add(items []any) {
	go func() {
		var results []map[string]interface{}
		for _, item := range items {
			jsonData, err := json.Marshal(item)
			if err != nil {
				log.Println(err)
			}

			var result map[string]interface{}
			err = json.Unmarshal(jsonData, &result)
			if err != nil {
				log.Println(err)
			}
			results = append(results, result)
		}

		_, err := s.client.Index("search").AddDocuments(results)
		if err != nil {
			log.Println(err)
		}
	}()
}
