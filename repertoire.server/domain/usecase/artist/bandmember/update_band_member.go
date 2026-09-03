package bandmember

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type UpdateBandMember struct {
	artistRepository repository.ArtistRepository
}

func NewUpdateBandMember(artistRepository repository.ArtistRepository) UpdateBandMember {
	return UpdateBandMember{
		artistRepository: artistRepository,
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

	var roles []model.BandMemberRole
	if err := u.artistRepository.GetBandMemberRolesByIDs(&roles, request.RoleIDs); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := u.artistRepository.ReplaceRolesFromBandMember(roles, &bandMember); err != nil {
		return httperror.DatabaseError(err)
	}

	bandMember.Name = request.Name
	bandMember.Color = request.Color
	if err := u.artistRepository.UpdateBandMember(&bandMember); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
