package song

import (
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type GetGuitarTunings struct {
	songRepository repository.SongRepository
	jwtService     service.JwtService
}

func NewGetGuitarTunings(
	songRepository repository.SongRepository,
	jwtService service.JwtService,
) GetGuitarTunings {
	return GetGuitarTunings{
		songRepository: songRepository,
		jwtService:     jwtService,
	}
}

func (g GetGuitarTunings) Handle(token string) (result []model.GuitarTuning, e *httperror.ErrorCode) {
	userID, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return result, errCode
	}

	if err := g.songRepository.GetGuitarTunings(&result, userID); err != nil {
		return result, httperror.DatabaseError(err)
	}

	return result, nil
}
