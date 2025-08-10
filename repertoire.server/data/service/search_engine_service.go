package service

import (
	"repertoire/server/data/search"
	"repertoire/server/internal/enums"
	"repertoire/server/internal/wrapper"
	"strings"

	"github.com/google/uuid"
	"github.com/meilisearch/meilisearch-go"
)

type SearchEngineService interface {
	Search(
		query string,
		currentPage *int,
		pageSize *int,
		searchType *enums.SearchType,
		userID uuid.UUID,
		filter []string,
		sort []string,
	) (wrapper.WithTotalCount[any], *wrapper.ErrorCode)
	GetDocument(id string) (map[string]any, error)
	GetDocuments(filter string) ([]map[string]any, error)
	Add(items []map[string]any) (int64, error)
	Update(items []map[string]any) (int64, error)
	Delete(ids []string) (int64, error)
	HasTaskSucceeded(status string) bool
}

type searchEngineService struct {
	client search.MeiliClient
}

func NewSearchEngineService(client search.MeiliClient) SearchEngineService {
	return searchEngineService{client: client}
}

func (s searchEngineService) Search(
	query string,
	currentPage *int,
	pageSize *int,
	searchType *enums.SearchType,
	userID uuid.UUID,
	filter []string,
	sort []string,
) (wrapper.WithTotalCount[any], *wrapper.ErrorCode) {
	request := &meilisearch.SearchRequest{Sort: sort}

	// pagination
	if currentPage != nil && pageSize != nil {
		request.Offset = int64((*currentPage - 1) * *pageSize)
		request.Limit = int64(*pageSize)
	}

	// filtering
	filters := []string{"userId = " + userID.String()}
	if searchType != nil {
		filters = append(filters, "type = "+string(*searchType))
	}
	if len(filter) > 0 {
		userFilter := "(" + strings.Join(filter, " AND ") + ")" // brackets, in case there is an OR
		filters = append(filters, userFilter)
	}
	request.Filter = strings.Join(filters, " AND ")

	// send search request
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

func (s searchEngineService) GetDocument(id string) (map[string]any, error) {
	var result map[string]any
	err := s.client.Index("search").GetDocument(id, &meilisearch.DocumentQuery{}, &result)
	if err != nil {
		return map[string]any{}, err
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

func (s searchEngineService) Add(items []map[string]any) (int64, error) {
	task, err := s.client.Index("search").AddDocuments(&items)
	if err != nil {
		return 0, err
	}
	return task.TaskUID, nil
}

func (s searchEngineService) Update(items []map[string]any) (int64, error) {
	task, err := s.client.Index("search").UpdateDocuments(&items)
	if err != nil {
		return 0, err
	}
	return task.TaskUID, nil
}

func (s searchEngineService) Delete(ids []string) (int64, error) {
	task, err := s.client.Index("search").DeleteDocuments(ids)
	if err != nil {
		return 0, err
	}
	return task.TaskUID, nil
}

func (s searchEngineService) HasTaskSucceeded(status string) bool {
	return status == string(meilisearch.TaskStatusSucceeded)
}
