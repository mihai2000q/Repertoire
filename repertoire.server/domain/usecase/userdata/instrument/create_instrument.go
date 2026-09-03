package instrument

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type CreateInstrument struct {
	repository repository.UserDataRepository
	jwtService service.JwtService
}

func NewCreateInstrument(repository repository.UserDataRepository, jwtService service.JwtService) CreateInstrument {
	return CreateInstrument{
		repository: repository,
		jwtService: jwtService,
	}
}

func (c CreateInstrument) Handle(request requests.CreateInstrumentRequest, token string) *httperror.ErrorCode {
	userID, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var count int64
	if err := c.repository.GetInstrumentsCount(&count, userID); err != nil {
		return httperror.DatabaseError(err)
	}

	guitarTuning := &model.Instrument{
		ID:     uuid.New(),
		Name:   request.Name,
		Order:  uint(count),
		UserID: userID,
	}
	if err := c.repository.CreateInstrument(guitarTuning); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
