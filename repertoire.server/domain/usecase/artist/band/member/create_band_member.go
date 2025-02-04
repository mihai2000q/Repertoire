package member

import (
	"errors"
	"github.com/google/uuid"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type CreateBandMember struct {
	artistRepository repository.ArtistRepository
}

func NewCreateBandMember(repository repository.ArtistRepository) CreateBandMember {
	return CreateBandMember{
		artistRepository: repository,
	}
}

func (c CreateBandMember) Handle(request requests.CreateBandMemberRequest) (uuid.UUID, *wrapper.ErrorCode) {
	var artist model.Artist
	err := c.artistRepository.GetWithBandMembers(&artist, request.ArtistID)
	if err != nil {
		return uuid.Nil, wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return uuid.Nil, wrapper.NotFoundError(errors.New("artist not found"))
	}
	if !artist.IsBand {
		return uuid.Nil, wrapper.BadRequestError(errors.New("artist is not band"))
	}

	var roles []model.BandMemberRole
	err = c.artistRepository.GetBandMemberRolesByIDs(&roles, request.RoleIDs)
	if err != nil {
		return uuid.Nil, wrapper.InternalServerError(err)
	}

	member := model.BandMember{
		ID:       uuid.New(),
		Name:     request.Name,
		Order:    uint(len(artist.BandMembers)),
		ArtistID: request.ArtistID,
		Roles:    roles,
	}
	err = c.artistRepository.CreateBandMember(&member)
	if err != nil {
		return uuid.Nil, wrapper.InternalServerError(err)
	}

	return member.ID, nil
}
