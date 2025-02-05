package song

import (
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type GetInstruments struct {
	repository repository.SongRepository
	jwtService service.JwtService
}

func NewGetInstruments(repository repository.SongRepository, jwtService service.JwtService) GetInstruments {
	return GetInstruments{
		repository: repository,
		jwtService: jwtService,
	}
}

func (g GetInstruments) Handle(token string) (result []model.Instrument, e *wrapper.ErrorCode) {
	userID, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return result, errCode
	}

	err := g.repository.GetInstruments(&result, userID)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	return result, nil
}
