package service

import (
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
