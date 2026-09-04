package song

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/pagination"
	"repertoire/server/model"
)

type GetAllSongs struct {
	songRepository repository.SongRepository
	jwtService     service.JwtService
}

func NewGetAllSongs(
	songRepository repository.SongRepository,
	jwtService service.JwtService,
) GetAllSongs {
	return GetAllSongs{
		songRepository: songRepository,
		jwtService:     jwtService,
	}
}

func (g GetAllSongs) Handle(request requests.GetSongsRequest, token string) (result pagination.WithTotalCount[model.EnhancedSong], e *httperror.ErrorCode) {
	userID, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return result, errCode
	}

	err := g.songRepository.GetAllByUser(
		&result.Models,
		userID,
		request.CurrentPage,
		request.PageSize,
		request.SearchBy,
		request.OrderBy,
		request.With,
	)
	if err != nil {
		return result, httperror.DatabaseError(err)
	}

	err = g.songRepository.GetAllByUserCount(&result.TotalCount, userID, request.SearchBy)
	if err != nil {
		return result, httperror.DatabaseError(err)
	}

	return result, nil
}
