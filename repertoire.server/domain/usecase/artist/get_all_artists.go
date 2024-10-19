package artist

import (
	"repertoire/api/request"
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

func (g GetAllArtists) Handle(request request.GetArtistsRequest, token string) (result wrapper.WithTotalCount[model.Artist], e *wrapper.ErrorCode) {
	userId, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return result, errCode
	}

	err := g.repository.GetAllByUser(
		&result.Models,
		userId,
		request.CurrentPage,
		request.PageSize,
		request.OrderBy,
	)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	err = g.repository.GetAllByUserCount(&result.TotalCount, userId)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	return result, nil
}
