package role

import (
	"errors"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"slices"

	"github.com/google/uuid"
)

type DeleteBandMemberRole struct {
	repository repository.ArtistRepository
	jwtService service.JwtService
}

func NewDeleteBandMemberRole(
	repository repository.ArtistRepository,
	jwtService service.JwtService,
) DeleteBandMemberRole {
	return DeleteBandMemberRole{
		repository: repository,
		jwtService: jwtService,
	}
}

func (d DeleteBandMemberRole) Handle(id uuid.UUID, token string) *wrapper.ErrorCode {
	userID, errCode := d.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var roles []model.BandMemberRole
	err := d.repository.GetBandMemberRoles(&roles, userID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	index := slices.IndexFunc(roles, func(s model.BandMemberRole) bool {
		return s.ID == id
	})
	if index == -1 {
		return wrapper.NotFoundError(errors.New("band member role not found"))
	}

	for i := index + 1; i < len(roles); i++ {
		roles[i].Order = roles[i].Order - 1
	}

	err = d.repository.UpdateAllBandMemberRoles(&roles)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = d.repository.DeleteBandMemberRole(id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
