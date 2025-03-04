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
) (wrapper.WithTotalCount[any], *wrapper.ErrorCode) {
	args := s.Called(query, currentPage, pageSize, searchType, userID)

	var errCode *wrapper.ErrorCode
	if a := args.Get(1); a != nil {
		errCode = a.(*wrapper.ErrorCode)
	}

	return args.Get(0).(wrapper.WithTotalCount[any]), errCode
}

func (s *SearchEngineServiceMock) Add(items []any) {
	s.Called(items)
}
