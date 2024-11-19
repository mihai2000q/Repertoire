package tuning

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type MoveGuitarTuning struct {
	repository repository.SongRepository
	jwtService service.JwtService
}

func NewMoveGuitarTuning(repository repository.SongRepository, jwtService service.JwtService) MoveGuitarTuning {
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

	index, overIndex, err := m.getIndexes(tunings, request.ID, request.OverID)
	if err != nil {
		return wrapper.NotFoundError(err)
	}
	tunings = m.move(tunings, index, overIndex)

	err = m.repository.UpdateAllGuitarTunings(&tunings)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (c MoveGuitarTuning) getIndexes(tunings []model.GuitarTuning, id uuid.UUID, overID uuid.UUID) (int, int, error) {
	var index *int
	var overIndex *int
	for i := 0; i < len(tunings); i++ {
		if tunings[i].ID == id {
			index = &i
		} else if tunings[i].ID == overID {
			overIndex = &i
		}
	}

	if index == nil {
		return -1, -1, errors.New("tuning not found")
	}
	if overIndex == nil {
		return -1, -1, errors.New("over tuning not found")
	}

	return *index, *overIndex, nil
}

func (c MoveGuitarTuning) move(tunings []model.GuitarTuning, index int, overIndex int) []model.GuitarTuning {
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
