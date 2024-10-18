package song

import (
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/models"
	"repertoire/utils/wrapper"
)

type GetAllSongs struct {
	repository repository.SongRepository
	jwtService service.JwtService
}

func NewGetAllSongs(repository repository.SongRepository, jwtService service.JwtService) GetAllSongs {
	return GetAllSongs{
		repository: repository,
		jwtService: jwtService,
	}
}

func (g GetAllSongs) Handle(request requests.GetSongsRequest, token string) (result wrapper.WithTotalCount[models.Song], e *wrapper.ErrorCode) {
	userId, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return result, errCode
	}

	err := g.repository.GetAllByUser(&result.Data, userId, request.CurrentPage, request.PageSize)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	err = g.repository.GetAllByUserCount(&result.TotalCount, userId)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	return result, nil
}
