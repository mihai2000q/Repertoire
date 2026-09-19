package guitartuning

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/reorder"
	"repertoire/server/model"
)

type MoveGuitarTuning struct {
	userDataRepository repository.UserDataRepository
	jwtService         service.JwtService
}

func NewMoveGuitarTuning(
	userDataRepository repository.UserDataRepository,
	jwtService service.JwtService,
) MoveGuitarTuning {
	return MoveGuitarTuning{
		userDataRepository: userDataRepository,
		jwtService:         jwtService,
	}
}

func (m MoveGuitarTuning) Handle(request requests.MoveGuitarTuningRequest, token string) *httperror.ErrorCode {
	userID, errCode := m.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var tunings []model.GuitarTuning
	if err := m.userDataRepository.GetGuitarTunings(&tunings, userID); err != nil {
		return httperror.DatabaseError(err)
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

	if err := m.userDataRepository.UpdateAllGuitarTunings(&tunings); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
