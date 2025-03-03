package service

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/meilisearch/meilisearch-go"
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
	Add(items []any) *wrapper.ErrorCode
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

	for _, hit := range searchResult.Hits {
		fmt.Println(hit)
	}

	result := wrapper.WithTotalCount[any]{
		Models:     searchResult.Hits,
		TotalCount: searchResult.EstimatedTotalHits,
	}

	return result, nil
}

func (s searchEngineService) Add(items []any) *wrapper.ErrorCode {
	var results []map[string]interface{}
	for _, item := range items {
		jsonData, err := json.Marshal(item)
		if err != nil {
			return wrapper.InternalServerError(err)
		}

		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
		results = append(results, result)
	}

	_, err := s.client.Index("search").AddDocuments(results)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
