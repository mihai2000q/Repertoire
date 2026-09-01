package guitartuning

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/reorder"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type MoveGuitarTuning struct {
	repository repository.UserDataRepository
	jwtService service.JwtService
}

func NewMoveGuitarTuning(repository repository.UserDataRepository, jwtService service.JwtService) MoveGuitarTuning {
	return MoveGuitarTuning{
		repository: repository,
		jwtService: jwtService,
	}
}

func (m MoveGuitarTuning) Handle(request requests.MoveGuitarTuningRequest, token string) *wrapper.ErrorCode {
	userID, errCode := m.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var tunings []model.GuitarTuning
	err := m.repository.GetGuitarTunings(&tunings, userID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	errCode = reorder.MoveEntity(
		tunings,
		request.ID,
		request.OverID,
		&reorder.Config{
			EntityNotFoundMsg:     "tuning not found",
			OverEntityNotFoundMsg: "over tuning not found",
		},
	)
	if errCode != nil {
		return errCode
	}

	err = m.repository.UpdateAllGuitarTunings(&tunings)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
