package validator

import (
	"errors"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type BandMemberValidator interface {
	Validate(ids []uuid.UUID, song model.Song) ([]model.BandMember, *httperror.ErrorCode)
}

type bandMemberValidator struct {
	artistRepository repository.ArtistRepository
}

func NewBandMemberValidator(artistRepository repository.ArtistRepository) BandMemberValidator {
	return bandMemberValidator{
		artistRepository: artistRepository,
	}
}

func (b bandMemberValidator) Validate(ids []uuid.UUID, song model.Song) ([]model.BandMember, *httperror.ErrorCode) {
	if len(ids) == 0 {
		return nil, nil
	}

	var bandMembers []model.BandMember
	if err := b.artistRepository.GetBandMembersByIDs(&bandMembers, ids); err != nil {
		return nil, httperror.DatabaseError(err)
	}
	if len(bandMembers) != len(ids) {
		return nil, httperror.NotFoundError(errors.New("band members not found"))
	}

	for _, bandMember := range bandMembers {
		if song.ArtistID == nil || *song.ArtistID != bandMember.ArtistID {
			return nil, httperror.ConflictError(errors.New("band member is not part of the artist associated with this song"))
		}
	}

	return bandMembers, nil
}
