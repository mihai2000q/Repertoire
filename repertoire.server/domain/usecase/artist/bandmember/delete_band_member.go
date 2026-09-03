package bandmember

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
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

func (d DeleteBandMember) Handle(id uuid.UUID, songID uuid.UUID) *httperror.ErrorCode {
	var artist model.Artist
	if err := d.artistRepository.GetWithBandMembers(&artist, songID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return httperror.NotFoundError(errors.New("artist not found"))
	}

	index := slices.IndexFunc(artist.BandMembers, func(a model.BandMember) bool {
		return a.ID == id
	})
	if index == -1 {
		return httperror.NotFoundError(errors.New("band member not found"))
	}

	sectionsLength := len(artist.BandMembers)
	for i := index + 1; i < sectionsLength; i++ {
		artist.BandMembers[i].Order = artist.BandMembers[i].Order - 1
	}

	if err := d.artistRepository.UpdateWithAssociations(&artist); err != nil {
		return httperror.DatabaseError(err)
	}
	if err := d.artistRepository.DeleteBandMember(id); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
