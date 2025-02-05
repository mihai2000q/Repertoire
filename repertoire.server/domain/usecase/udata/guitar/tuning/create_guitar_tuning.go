package tuning

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type CreateGuitarTuning struct {
	repository repository.UserDataRepository
	jwtService service.JwtService
}

func NewCreateGuitarTuning(repository repository.UserDataRepository, jwtService service.JwtService) CreateGuitarTuning {
	return CreateGuitarTuning{
		repository: repository,
		jwtService: jwtService,
	}
}

func (c CreateGuitarTuning) Handle(request requests.CreateGuitarTuningRequest, token string) *wrapper.ErrorCode {
	userID, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var count int64
	err := c.repository.GetGuitarTuningsCount(&count, userID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	guitarTuning := &model.GuitarTuning{
		ID:     uuid.New(),
		Name:   request.Name,
		Order:  uint(count),
		UserID: userID,
	}
	err = c.repository.CreateGuitarTuning(guitarTuning)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
