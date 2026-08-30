package instrument

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/reorder"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type MoveInstrument struct {
	repository repository.UserDataRepository
	jwtService service.JwtService
}

func NewMoveInstrument(repository repository.UserDataRepository, jwtService service.JwtService) MoveInstrument {
	return MoveInstrument{
		repository: repository,
		jwtService: jwtService,
	}
}

func (m MoveInstrument) Handle(request requests.MoveInstrumentRequest, token string) *wrapper.ErrorCode {
	userID, errCode := m.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var instruments []model.Instrument
	err := m.repository.GetInstruments(&instruments, userID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	errCode = reorder.MoveEntity(
		instruments,
		request.ID,
		request.OverID,
		&reorder.Config{
			EntityNotFoundMsg:     "instrument not found",
			OverEntityNotFoundMsg: "over instrument not found",
		},
	)
	if errCode != nil {
		return errCode
	}

	err = m.repository.UpdateAllInstruments(&instruments)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
