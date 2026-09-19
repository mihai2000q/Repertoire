package role

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

type DeleteBandMemberRole struct {
	userDataRepository repository.UserDataRepository
	jwtService         service.JwtService
	transactionManager transaction.Manager
}

func NewDeleteBandMemberRole(
	userDataRepository repository.UserDataRepository,
	jwtService service.JwtService,
	transactionManager transaction.Manager,
) DeleteBandMemberRole {
	return DeleteBandMemberRole{
		userDataRepository: userDataRepository,
		jwtService:         jwtService,
		transactionManager: transactionManager,
	}
}

func (d DeleteBandMemberRole) Handle(id uuid.UUID, token string) *httperror.ErrorCode {
	userID, errCode := d.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var roles []model.BandMemberRole
	if err := d.userDataRepository.GetBandMemberRoles(&roles, userID); err != nil {
		return httperror.DatabaseError(err)
	}

	index := slices.IndexFunc(roles, func(s model.BandMemberRole) bool {
		return s.ID == id
	})
	if index == -1 {
		return httperror.NotFoundError(errors.New("band member role not found"))
	}

	for i := index + 1; i < len(roles); i++ {
		roles[i].Order = roles[i].Order - 1
	}

	err := d.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		txUserDataRepo := factory.NewUserDataRepository()
		if err := txUserDataRepo.UpdateAllBandMemberRoles(&roles); err != nil {
			return err
		}
		if err := txUserDataRepo.DeleteBandMemberRole(id); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
