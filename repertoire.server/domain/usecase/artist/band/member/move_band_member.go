package member

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type MoveBandMember struct {
	artistRepository repository.ArtistRepository
}

func NewMoveBandMember(repository repository.ArtistRepository) MoveBandMember {
	return MoveBandMember{
		artistRepository: repository,
	}
}

func (m MoveBandMember) Handle(request requests.MoveBandMemberRequest) *wrapper.ErrorCode {
	var artist model.Artist
	err := m.artistRepository.GetWithBandMembers(&artist, request.ArtistID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return wrapper.NotFoundError(errors.New("artist not found"))
	}

	index, overIndex, err := m.getIndexes(artist.BandMembers, request.ID, request.OverID)
	if err != nil {
		return wrapper.NotFoundError(err)
	}
	artist.BandMembers = m.move(artist.BandMembers, index, overIndex)

	err = m.artistRepository.UpdateWithAssociations(&artist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (MoveBandMember) getIndexes(bandMembers []model.BandMember, id uuid.UUID, overID uuid.UUID) (int, int, error) {
	var index *int
	var overIndex *int
	for i := 0; i < len(bandMembers); i++ {
		if bandMembers[i].ID == id {
			index = &i
		} else if bandMembers[i].ID == overID {
			overIndex = &i
		}
	}

	if index == nil {
		return -1, -1, errors.New("band member not found")
	}
	if overIndex == nil {
		return -1, -1, errors.New("over band member not found")
	}

	return *index, *overIndex, nil
}

func (MoveBandMember) move(bandMembers []model.BandMember, index int, overIndex int) []model.BandMember {
	if index < overIndex {
		for i := index + 1; i <= overIndex; i++ {
			bandMembers[i].Order = uint(i - 1)
		}
	} else {
		for i := overIndex; i <= index; i++ {
			bandMembers[i].Order = uint(i + 1)
		}
	}

	bandMembers[index].Order = uint(overIndex)

	return bandMembers
}
