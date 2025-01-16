package tuning

import (
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type GetGuitarTunings struct {
	repository repository.SongRepository
	jwtService service.JwtService
}

func NewGetGuitarTunings(repository repository.SongRepository, jwtService service.JwtService) GetGuitarTunings {
	return GetGuitarTunings{
		repository: repository,
		jwtService: jwtService,
	}
}

func (g GetGuitarTunings) Handle(token string) (result []model.GuitarTuning, e *wrapper.ErrorCode) {
	userID, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return result, errCode
	}

	err := g.repository.GetGuitarTunings(&result, userID)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	return result, nil
}
