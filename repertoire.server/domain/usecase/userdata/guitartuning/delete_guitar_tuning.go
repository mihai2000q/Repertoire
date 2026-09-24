package guitartuning

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

type DeleteGuitarTuning struct {
	userDataRepository repository.UserDataRepository
	jwtService         service.JwtService
	transactionManager transaction.Manager
}

func NewDeleteGuitarTuning(
	userDataRepository repository.UserDataRepository,
	jwtService service.JwtService,
	transactionManager transaction.Manager,
) DeleteGuitarTuning {
	return DeleteGuitarTuning{
		userDataRepository: userDataRepository,
		jwtService:         jwtService,
		transactionManager: transactionManager,
	}
}

func (d DeleteGuitarTuning) Handle(id uuid.UUID, token string) *httperror.ErrorCode {
	userID, errCode := d.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var tunings []model.GuitarTuning
	if err := d.userDataRepository.GetGuitarTunings(&tunings, userID); err != nil {
		return httperror.DatabaseError(err)
	}

	index := slices.IndexFunc(tunings, func(t model.GuitarTuning) bool {
		return t.ID == id
	})
	if index == -1 {
		return httperror.NotFoundError(errors.New("guitar tuning not found"))
	}

	for i := index + 1; i < len(tunings); i++ {
		tunings[i].Order = tunings[i].Order - 1
	}

	err := d.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		txUserDataRepo := factory.NewUserDataRepository()
		if err := txUserDataRepo.UpdateAllGuitarTunings(&tunings); err != nil {
			return err
		}
		if err := txUserDataRepo.DeleteGuitarTuning(id); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
