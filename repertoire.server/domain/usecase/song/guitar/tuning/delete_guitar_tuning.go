package tuning

import (
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"slices"

	"github.com/google/uuid"
)

type DeleteGuitarTuning struct {
	repository repository.SongRepository
	jwtService service.JwtService
}

func NewDeleteGuitarTuning(repository repository.SongRepository, jwtService service.JwtService) DeleteGuitarTuning {
	return DeleteGuitarTuning{
		repository: repository,
		jwtService: jwtService,
	}
}

func (d DeleteGuitarTuning) Handle(id uuid.UUID, token string) *wrapper.ErrorCode {
	userID, errCode := d.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var tunings []model.GuitarTuning
	err := d.repository.GetGuitarTunings(&tunings, userID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	index := slices.IndexFunc(tunings, func(t model.GuitarTuning) bool {
		return t.ID == id
	})

	for i := index + 1; i < len(tunings); i++ {
		tunings[i].Order = uint(tunings[i].Order - 1)
	}

	err = d.repository.UpdateAllGuitarTunings(&tunings, index)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = d.repository.DeleteGuitarTuning(id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
