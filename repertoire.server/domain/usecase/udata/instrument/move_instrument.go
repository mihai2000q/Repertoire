package instrument

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
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

	index, overIndex, err := m.getIndexes(instruments, request.ID, request.OverID)
	if err != nil {
		return wrapper.NotFoundError(err)
	}
	instruments = m.move(instruments, index, overIndex)

	err = m.repository.UpdateAllInstruments(&instruments)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (MoveInstrument) getIndexes(instruments []model.Instrument, id uuid.UUID, overID uuid.UUID) (int, int, error) {
	var index *int
	var overIndex *int
	for i := 0; i < len(instruments); i++ {
		if instruments[i].ID == id {
			index = &i
		} else if instruments[i].ID == overID {
			overIndex = &i
		}
	}

	if index == nil {
		return -1, -1, errors.New("instrument not found")
	}
	if overIndex == nil {
		return -1, -1, errors.New("over instrument not found")
	}

	return *index, *overIndex, nil
}

func (MoveInstrument) move(tunings []model.Instrument, index int, overIndex int) []model.Instrument {
	if index < overIndex {
		for i := index + 1; i <= overIndex; i++ {
			tunings[i].Order = uint(i - 1)
		}
	} else {
		for i := overIndex; i <= index; i++ {
			tunings[i].Order = uint(i + 1)
		}
	}

	tunings[index].Order = uint(overIndex)

	return tunings
}
