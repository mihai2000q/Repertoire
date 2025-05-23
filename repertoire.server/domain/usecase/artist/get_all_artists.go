package artist

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type GetAllArtists struct {
	repository repository.ArtistRepository
	jwtService service.JwtService
}

func NewGetAllArtists(repository repository.ArtistRepository, jwtService service.JwtService) GetAllArtists {
	return GetAllArtists{
		repository: repository,
		jwtService: jwtService,
	}
}

func (g GetAllArtists) Handle(request requests.GetArtistsRequest, token string) (result wrapper.WithTotalCount[model.EnhancedArtist], e *wrapper.ErrorCode) {
	userID, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return result, errCode
	}

	err := g.repository.GetAllByUser(
		&result.Models,
		userID,
		request.CurrentPage,
		request.PageSize,
		request.OrderBy,
		request.SearchBy,
	)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	err = g.repository.GetAllByUserCount(&result.TotalCount, userID, request.SearchBy)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	return result, nil
}
