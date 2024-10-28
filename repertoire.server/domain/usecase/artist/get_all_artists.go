package artist

import (
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/model"
	"repertoire/utils/wrapper"
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

func (g GetAllArtists) Handle(request requests.GetArtistsRequest, token string) (result wrapper.WithTotalCount[model.Artist], e *wrapper.ErrorCode) {
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
	)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	err = g.repository.GetAllByUserCount(&result.TotalCount, userID)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	return result, nil
}
