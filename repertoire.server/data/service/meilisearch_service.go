package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/meilisearch/meilisearch-go"
	"repertoire/server/internal/enums"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type MeiliSearchService interface {
	Get(
		query string,
		currentPage *int,
		pageSize *int,
		searchType *enums.SearchType,
		userID uuid.UUID,
	) (wrapper.WithTotalCount[model.SearchBase], *wrapper.ErrorCode)
}

type meiliSearchService struct {
	client meilisearch.ServiceManager
}

func NewMeiliSearchService(client meilisearch.ServiceManager) MeiliSearchService {
	return meiliSearchService{client: client}
}

func (m meiliSearchService) Get(
	query string,
	currentPage *int,
	pageSize *int,
	searchType *enums.SearchType,
	userID uuid.UUID,
) (wrapper.WithTotalCount[model.SearchBase], *wrapper.ErrorCode) {
	request := &meilisearch.SearchRequest{}

	if currentPage != nil && pageSize != nil {
		request.Offset = int64((*currentPage - 1) * *pageSize)
		request.Limit = int64(*pageSize)
	}

	if searchType != nil {
		request.Filter = "type = " + string(*searchType) + " AND user_id = " + userID.String()
	} else {
		request.Filter = "user_id = " + userID.String()
	}

	searchResult, err := m.client.Index("search").Search(query, request)
	if err != nil {
		return wrapper.WithTotalCount[model.SearchBase]{}, wrapper.InternalServerError(err)
	}

	for _, hit := range searchResult.Hits {
		fmt.Println(hit)
	}

	result := wrapper.WithTotalCount[model.SearchBase]{
		Models:     searchResult.Hits,
		TotalCount: searchResult.TotalHits,
	}

	return result, nil
}
