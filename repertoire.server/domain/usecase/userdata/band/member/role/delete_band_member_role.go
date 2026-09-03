package role

import (
	"errors"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
	"slices"

	"github.com/google/uuid"
)

type DeleteBandMemberRole struct {
	repository repository.UserDataRepository
	jwtService service.JwtService
}

func NewDeleteBandMemberRole(
	repository repository.UserDataRepository,
	jwtService service.JwtService,
) DeleteBandMemberRole {
	return DeleteBandMemberRole{
		repository: repository,
		jwtService: jwtService,
	}
}

func (d DeleteBandMemberRole) Handle(id uuid.UUID, token string) *httperror.ErrorCode {
	userID, errCode := d.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var roles []model.BandMemberRole
	if err := d.repository.GetBandMemberRoles(&roles, userID); err != nil {
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

	if err := d.repository.UpdateAllBandMemberRoles(&roles); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := d.repository.DeleteBandMemberRole(id); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
