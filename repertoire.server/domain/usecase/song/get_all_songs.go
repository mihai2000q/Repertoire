package song

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
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

func (g GetAllSongs) Handle(request requests.GetSongsRequest, token string) (result wrapper.WithTotalCount[model.Song], e *wrapper.ErrorCode) {
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
