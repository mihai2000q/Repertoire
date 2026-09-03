package song

import (
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
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

func (g GetInstruments) Handle(token string) (result []model.Instrument, e *httperror.ErrorCode) {
	userID, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return result, errCode
	}

	if err := g.repository.GetInstruments(&result, userID); err != nil {
		return result, httperror.DatabaseError(err)
	}

	return result, nil
}
