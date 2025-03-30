package service

import (
	"io"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/search"
	"repertoire/server/internal/wrapper"
)

type SearchService interface {
	Get(request requests.SearchGetRequest, token string) (wrapper.WithTotalCount[any], *wrapper.ErrorCode)
	MeiliWebhook(requestBody io.ReadCloser) *wrapper.ErrorCode
}

type searchService struct {
	get          search.Get
	meiliWebhook search.MeiliWebhook
}

func NewSearchService(get search.Get, meiliWebhook search.MeiliWebhook) SearchService {
	return &searchService{
		get:          get,
		meiliWebhook: meiliWebhook,
	}
}

func (s searchService) Get(request requests.SearchGetRequest, token string) (wrapper.WithTotalCount[any], *wrapper.ErrorCode) {
	return s.get.Handle(request, token)
}

func (s searchService) MeiliWebhook(requestBody io.ReadCloser) *wrapper.ErrorCode {
	return s.meiliWebhook.Handle(requestBody)
}
