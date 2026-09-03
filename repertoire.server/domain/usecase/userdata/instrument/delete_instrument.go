package instrument

import (
	"errors"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
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

func (d DeleteInstrument) Handle(id uuid.UUID, token string) *httperror.ErrorCode {
	userID, errCode := d.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var instruments []model.Instrument
	if err := d.repository.GetInstruments(&instruments, userID); err != nil {
		return httperror.DatabaseError(err)
	}

	index := slices.IndexFunc(instruments, func(t model.Instrument) bool {
		return t.ID == id
	})
	if index == -1 {
		return httperror.NotFoundError(errors.New("instrument not found"))
	}

	for i := index + 1; i < len(instruments); i++ {
		instruments[i].Order = instruments[i].Order - 1
	}

	if err := d.repository.UpdateAllInstruments(&instruments); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := d.repository.DeleteInstrument(id); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
