package service

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"repertoire/server/internal/enums"
	"repertoire/server/internal/wrapper"
)

type SearchEngineServiceMock struct {
	mock.Mock
}

func (s *SearchEngineServiceMock) Search(
	query string,
	currentPage *int,
	pageSize *int,
	searchType *enums.SearchType,
	userID uuid.UUID,
	filter []string,
	sort []string,
) (wrapper.WithTotalCount[any], *wrapper.ErrorCode) {
	args := s.Called(query, currentPage, pageSize, searchType, userID, filter, sort)

	var errCode *wrapper.ErrorCode
	if a := args.Get(1); a != nil {
		errCode = a.(*wrapper.ErrorCode)
	}

	return args.Get(0).(wrapper.WithTotalCount[any]), errCode
}

func (s *SearchEngineServiceMock) GetDocument(id string) (map[string]any, error) {
	args := s.Called(id)
	return args.Get(0).(map[string]any), args.Error(1)
}

func (s *SearchEngineServiceMock) GetDocuments(filter string) ([]map[string]any, error) {
	args := s.Called(filter)
	return args.Get(0).([]map[string]any), args.Error(1)
}

func (s *SearchEngineServiceMock) Add(items []map[string]any) (int64, error) {
	args := s.Called(items)
	return args.Get(0).(int64), args.Error(1)
}

func (s *SearchEngineServiceMock) Update(items []map[string]any) (int64, error) {
	args := s.Called(items)
	return args.Get(0).(int64), args.Error(1)
}

func (s *SearchEngineServiceMock) Delete(ids []string) (int64, error) {
	args := s.Called(ids)
	return args.Get(0).(int64), args.Error(1)
}

func (s *SearchEngineServiceMock) HasTaskSucceeded(status string) bool {
	args := s.Called(status)
	return args.Bool(0)
}
