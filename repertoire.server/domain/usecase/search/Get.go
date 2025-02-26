package search

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/service"
	"repertoire/server/internal/enums"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type Get struct {
	jwtService         service.JwtService
	meiliSearchService service.MeiliSearchService
}

func NewGet(
	jwtService service.JwtService,
	meiliSearchService service.MeiliSearchService,
) Get {
	return Get{
		jwtService:         jwtService,
		meiliSearchService: meiliSearchService,
	}
}

func (g Get) Handle(
	request requests.SearchGetRequest,
	token string,
) (wrapper.WithTotalCount[model.SearchBase], *wrapper.ErrorCode) {
	userID, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return wrapper.WithTotalCount[model.SearchBase]{}, errCode
	}

	searchResult, errCode := g.meiliSearchService.Get(
		request.Query,
		request.CurrentPage,
		request.PageSize,
		request.Type,
		userID,
	)

	for i := range searchResult.Models {
		if searchResult.Models[i].Type == enums.Artist {

		}
	}

	return searchResult, errCode
}
