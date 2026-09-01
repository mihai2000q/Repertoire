package instrument

import (
	"errors"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"slices"

	"github.com/google/uuid"
)

type DeleteInstrument struct {
	repository repository.UserDataRepository
	jwtService service.JwtService
}

func NewDeleteInstrument(repository repository.UserDataRepository, jwtService service.JwtService) DeleteInstrument {
	return DeleteInstrument{
		repository: repository,
		jwtService: jwtService,
	}
}

func (d DeleteInstrument) Handle(id uuid.UUID, token string) *wrapper.ErrorCode {
	userID, errCode := d.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var instruments []model.Instrument
	err := d.repository.GetInstruments(&instruments, userID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	index := slices.IndexFunc(instruments, func(t model.Instrument) bool {
		return t.ID == id
	})
	if index == -1 {
		return wrapper.NotFoundError(errors.New("instrument not found"))
	}

	for i := index + 1; i < len(instruments); i++ {
		instruments[i].Order = instruments[i].Order - 1
	}

	err = d.repository.UpdateAllInstruments(&instruments)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = d.repository.DeleteInstrument(id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
