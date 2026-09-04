package bandmember

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type CreateBandMember struct {
	artistRepository repository.ArtistRepository
}

func NewCreateBandMember(artistRepository repository.ArtistRepository) CreateBandMember {
	return CreateBandMember{
		artistRepository: artistRepository,
	}
}

func (c CreateBandMember) Handle(request requests.CreateBandMemberRequest) (uuid.UUID, *httperror.ErrorCode) {
	var artist model.Artist
	if err := c.artistRepository.GetWithBandMembers(&artist, request.ArtistID); err != nil {
		return uuid.Nil, httperror.DatabaseError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return uuid.Nil, httperror.NotFoundError(errors.New("artist not found"))
	}
	if !artist.IsBand {
		return uuid.Nil, httperror.ConflictError(errors.New("artist is not band"))
	}

	var roles []model.BandMemberRole
	if err := c.artistRepository.GetBandMemberRolesByIDs(&roles, request.RoleIDs); err != nil {
		return uuid.Nil, httperror.DatabaseError(err)
	}
	if len(roles) != len(request.RoleIDs) {
		return uuid.Nil, httperror.NotFoundError(errors.New("roles not found"))
	}

	member := model.BandMember{
		ID:       uuid.New(),
		Name:     request.Name,
		Color:    request.Color,
		Order:    uint(len(artist.BandMembers)),
		ArtistID: request.ArtistID,
		Roles:    roles,
	}
	if err := c.artistRepository.CreateBandMember(&member); err != nil {
		return uuid.Nil, httperror.DatabaseError(err)
	}

	return member.ID, nil
}
