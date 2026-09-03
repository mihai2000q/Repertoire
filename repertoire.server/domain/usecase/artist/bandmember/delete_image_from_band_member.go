package bandmember

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type DeleteImageFromBandMember struct {
	artistRepository repository.ArtistRepository
	storageService   service.StorageService
}

func NewDeleteImageFromBandMember(
	artistRepository repository.ArtistRepository,
	storageService service.StorageService,
) DeleteImageFromBandMember {
	return DeleteImageFromBandMember{
		artistRepository: artistRepository,
		storageService:   storageService,
	}
}

func (d DeleteImageFromBandMember) Handle(id uuid.UUID) *httperror.ErrorCode {
	var member model.BandMember
	if err := d.artistRepository.GetBandMember(&member, id); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(member).IsZero() {
		return httperror.NotFoundError(errors.New("band member not found"))
	}
	if member.ImageURL == nil {
		return nil
	}

	if errCode := d.storageService.DeleteFile(*member.ImageURL); errCode != nil {
		return errCode
	}

	member.ImageURL = nil
	if err := d.artistRepository.UpdateBandMember(&member); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
