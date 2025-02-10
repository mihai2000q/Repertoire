package member

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"slices"

	"github.com/google/uuid"
)

type DeleteBandMember struct {
	artistRepository repository.ArtistRepository
}

func NewDeleteBandMember(repository repository.ArtistRepository) DeleteBandMember {
	return DeleteBandMember{
		artistRepository: repository,
	}
}

func (d DeleteBandMember) Handle(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode {
	var artist model.Artist
	err := d.artistRepository.GetWithBandMembers(&artist, songID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return wrapper.NotFoundError(errors.New("artist not found"))
	}

	index := slices.IndexFunc(artist.BandMembers, func(a model.BandMember) bool {
		return a.ID == id
	})
	if index == -1 {
		return wrapper.NotFoundError(errors.New("band member not found"))
	}

	sectionsLength := len(artist.BandMembers)
	for i := index + 1; i < sectionsLength; i++ {
		artist.BandMembers[i].Order = artist.BandMembers[i].Order - 1
	}

	err = d.artistRepository.UpdateWithAssociations(&artist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	err = d.artistRepository.DeleteBandMember(id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
