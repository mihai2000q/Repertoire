package bandmember

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type UpdateBandMember struct {
	artistRepository   repository.ArtistRepository
	transactionManager transaction.Manager
}

func NewUpdateBandMember(
	artistRepository repository.ArtistRepository,
	transactionManager transaction.Manager,
) UpdateBandMember {
	return UpdateBandMember{
		artistRepository:   artistRepository,
		transactionManager: transactionManager,
	}
}

func (u UpdateBandMember) Handle(request requests.UpdateBandMemberRequest) *httperror.ErrorCode {
	var bandMember model.BandMember
	if err := u.artistRepository.GetBandMember(&bandMember, request.ID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(bandMember).IsZero() {
		return httperror.NotFoundError(errors.New("band member not found"))
	}

	err := u.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		txArtistRepo := factory.NewArtistRepository()

		var roles []model.BandMemberRole
		if err := txArtistRepo.GetBandMemberRolesByIDs(&roles, request.RoleIDs); err != nil {
			return err
		}
		if err := txArtistRepo.ReplaceRolesFromBandMember(roles, &bandMember); err != nil {
			return err
		}

		bandMember.Name = request.Name
		bandMember.Color = request.Color
		if err := txArtistRepo.UpdateBandMember(&bandMember); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
