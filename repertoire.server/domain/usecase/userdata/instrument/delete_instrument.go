package instrument

import (
	"errors"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
	"slices"

	"github.com/google/uuid"
)

type DeleteInstrument struct {
	userDataRepository repository.UserDataRepository
	jwtService         service.JwtService
	transactionManager transaction.Manager
}

func NewDeleteInstrument(
	userDataRepository repository.UserDataRepository,
	jwtService service.JwtService,
	transactionManager transaction.Manager,
) DeleteInstrument {
	return DeleteInstrument{
		userDataRepository: userDataRepository,
		jwtService:         jwtService,
		transactionManager: transactionManager,
	}
}

func (d DeleteInstrument) Handle(id uuid.UUID, token string) *httperror.ErrorCode {
	userID, errCode := d.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var instruments []model.Instrument
	if err := d.userDataRepository.GetInstruments(&instruments, userID); err != nil {
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

	err := d.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		txUserDataRepo := factory.NewUserDataRepository()
		if err := txUserDataRepo.UpdateAllInstruments(&instruments); err != nil {
			return err
		}
		if err := txUserDataRepo.DeleteInstrument(id); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
