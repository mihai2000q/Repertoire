package member

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type UpdateBandMember struct {
	artistRepository repository.ArtistRepository
}

func NewUpdateBandMember(repository repository.ArtistRepository) UpdateBandMember {
	return UpdateBandMember{
		artistRepository: repository,
	}
}

func (u UpdateBandMember) Handle(request requests.UpdateBandMemberRequest) *wrapper.ErrorCode {
	var bandMember model.BandMember
	err := u.artistRepository.GetBandMember(&bandMember, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(bandMember).IsZero() {
		return wrapper.NotFoundError(errors.New("band member not found"))
	}

	var roles []model.BandMemberRole
	err = u.artistRepository.GetBandMemberRolesByIDs(&roles, request.RoleIDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = u.artistRepository.ReplaceRolesFromBandMember(&roles, &bandMember)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	bandMember.Name = request.Name
	err = u.artistRepository.UpdateBandMember(&bandMember)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
