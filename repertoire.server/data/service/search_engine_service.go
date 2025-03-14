package service

import (
	"github.com/google/uuid"
	"github.com/meilisearch/meilisearch-go"
	"repertoire/server/internal/enums"
	"repertoire/server/internal/wrapper"
)

type SearchEngineService interface {
	Search(
		query string,
		currentPage *int,
		pageSize *int,
		searchType *enums.SearchType,
		userID uuid.UUID,
	) (wrapper.WithTotalCount[any], *wrapper.ErrorCode)
	GetDocuments(filter string) ([]map[string]any, error)
	Add(items []any) error
	Update(items []any) error
	Delete(ids []string) error
}

type searchEngineService struct {
	client meilisearch.ServiceManager
}

func NewSearchEngineService(client meilisearch.ServiceManager) SearchEngineService {
	return searchEngineService{client: client}
}

func (s searchEngineService) Search(
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

func (s searchEngineService) GetDocuments(filter string) ([]map[string]any, error) {
	var result meilisearch.DocumentsResult
	err := s.client.Index("search").GetDocuments(&meilisearch.DocumentsQuery{
		Filter: filter,
	}, &result)
	if err != nil {
		return []map[string]any{}, err
	}

	return result.Results, nil
}

func (s searchEngineService) Add(items []any) error {
	_, err := s.client.Index("search").AddDocuments(&items)
	return err
}

func (s searchEngineService) Update(items []any) error {
	_, err := s.client.Index("search").UpdateDocuments(&items)
	if err != nil {
		return err
	}
	return nil
}

func (s searchEngineService) Delete(ids []string) error {
	_, err := s.client.Index("search").DeleteDocuments(ids)
	if err != nil {
		return err
	}
	return nil
}
