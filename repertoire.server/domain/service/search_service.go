package service

import (
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/search"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type SearchService interface {
	Get(request requests.SearchGetRequest, token string) (wrapper.WithTotalCount[model.SearchBase], *wrapper.ErrorCode)
}

type searchService struct {
	get search.Get
}

func NewSearchService(get search.Get) SearchService {
	return &searchService{get: get}
}

func (s searchService) Get(request requests.SearchGetRequest, token string) (wrapper.WithTotalCount[model.SearchBase], *wrapper.ErrorCode) {
	return s.get.Handle(request, token)
}
