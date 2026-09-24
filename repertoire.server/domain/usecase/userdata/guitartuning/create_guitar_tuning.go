package guitartuning

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type CreateGuitarTuning struct {
	userDataRepository repository.UserDataRepository
	jwtService         service.JwtService
}

func NewCreateGuitarTuning(
	userDataRepository repository.UserDataRepository,
	jwtService service.JwtService,
) CreateGuitarTuning {
	return CreateGuitarTuning{
		userDataRepository: userDataRepository,
		jwtService:         jwtService,
	}
}

func (c CreateGuitarTuning) Handle(request requests.CreateGuitarTuningRequest, token string) *httperror.ErrorCode {
	userID, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var count int64
	if err := c.userDataRepository.GetGuitarTuningsCount(&count, userID); err != nil {
		return httperror.DatabaseError(err)
	}

	guitarTuning := &model.GuitarTuning{
		ID:     uuid.New(),
		Name:   request.Name,
		Order:  uint(count),
		UserID: userID,
	}
	if err := c.userDataRepository.CreateGuitarTuning(guitarTuning); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
