package types

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

type DeleteSongSectionType struct {
	userDataRepository repository.UserDataRepository
	jwtService         service.JwtService
	transactionManager transaction.Manager
}

func NewDeleteSongSectionType(
	userDataRepository repository.UserDataRepository,
	jwtService service.JwtService,
	transactionManager transaction.Manager,
) DeleteSongSectionType {
	return DeleteSongSectionType{
		userDataRepository: userDataRepository,
		jwtService:         jwtService,
		transactionManager: transactionManager,
	}
}

func (d DeleteSongSectionType) Handle(id uuid.UUID, token string) *httperror.ErrorCode {
	userID, errCode := d.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var types []model.SongSectionType
	if err := d.userDataRepository.GetSectionTypes(&types, userID); err != nil {
		return httperror.DatabaseError(err)
	}

	index := slices.IndexFunc(types, func(s model.SongSectionType) bool {
		return s.ID == id
	})
	if index == -1 {
		return httperror.NotFoundError(errors.New("song section type not found"))
	}

	for i := index + 1; i < len(types); i++ {
		types[i].Order = types[i].Order - 1
	}

	err := d.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		txUserDataRepo := factory.NewUserDataRepository()
		if err := txUserDataRepo.UpdateAllSectionTypes(&types); err != nil {
			return err
		}
		if err := txUserDataRepo.DeleteSectionType(id); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
