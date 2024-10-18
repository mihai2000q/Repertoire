package artist

import (
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/models"
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

func (g GetAllArtists) Handle(request requests.GetArtistsRequest, token string) (result wrapper.WithTotalCount[models.Artist], e *wrapper.ErrorCode) {
	userId, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return result, errCode
	}

	err := g.repository.GetAllByUser(&result.Models, userId, request.CurrentPage, request.PageSize)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	err = g.repository.GetAllByUserCount(&result.TotalCount, userId)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	return result, nil
}
